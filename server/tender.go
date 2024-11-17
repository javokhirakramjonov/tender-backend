package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"tender-backend/custom_errors"
	"tender-backend/gen_proto"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"tender-backend/rabbit_mq"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TenderService struct {
	db                  *gorm.DB
	redis               *redis.Client
	NotificationService *NotificationService
	notificationServer  *Server
}

// NewTenderService initializes a new TenderService with the database connection.
func NewTenderService(db *gorm.DB, redisClient *redis.Client) *TenderService {
	return &TenderService{
		db:                  db,
		redis:               redisClient,
		notificationServer:  NewNotificationServer(),
		NotificationService: NewNotificationService(db),
	}
}

// CreateTender creates a new tender in the database.
func (t *TenderService) CreateTender(req *request_model.CreateTenderReq, clientID int64) (*model.Tender, *custom_errors.AppError) {
	if err := validateCreateTender(req); err != nil {
		return nil, err
	}

	tender := &model.Tender{
		ClientID:    clientID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
		Status:      "open",
	}

	// Save the tender to the database.
	if err := t.db.Create(tender).Error; err != nil {
		return nil, custom_errors.NewAppError(err)
	}

	// Invalidate the cache after creating a new tender
	t.redis.Del(context.Background(), "tenders_cache")

	return tender, nil
}

// ValidateCreateTender validates the input for creating a tender.
func validateCreateTender(req *request_model.CreateTenderReq) *custom_errors.AppError {
	if req.Title == "" {
		return custom_errors.NewBadRequestError("Invalid input")
	}
	if req.Deadline.Before(time.Now()) {
		return custom_errors.NewBadRequestError("Invalid tender data")
	}

	if req.Budget <= 0 {
		return custom_errors.NewBadRequestError("Invalid tender data")
	}

	return nil
}

// GetTenderById retrieves a tender by its ID.
func (t *TenderService) GetTenderById(id int64) (*model.Tender, *custom_errors.AppError) {
	var tender model.Tender

	// Try fetching the tender from the database
	if err := t.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("Tender not found or access denied")
		}
		return nil, custom_errors.NewAppError(err)
	}

	return &tender, nil
}

// GetTenders retrieves all tenders from the cache or database.
func (t *TenderService) GetTenders() ([]model.Tender, error) {
	// Redis context
	ctx := context.Background()

	// Redis key for caching tenders
	cacheKey := "tenders_cache"

	// Try fetching the tenders from Redis
	cachedTenders, err := t.redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedTenders != "" {
		// If cached data exists, unmarshal it and return
		var tenders []model.Tender
		if err := json.Unmarshal([]byte(cachedTenders), &tenders); err == nil {
			return tenders, nil
		}
	}

	// If cache miss or unmarshal error, fetch from the database
	var tenders []model.Tender
	if err := t.db.Find(&tenders).Error; err != nil {
		return nil, err
	}

	// Marshal the tenders to JSON and store in Redis
	tendersJSON, err := json.Marshal(tenders)
	if err == nil {
		// Cache the tenders for 10 minutes
		_ = t.redis.Set(ctx, cacheKey, tendersJSON, 10*time.Minute).Err()
	}

	return tenders, nil
}

// UpdateTender updates the tender with the given ID.
func (t *TenderService) UpdateTender(tenderID, clientID int64, req *request_model.UpdateTenderReq) (*model.Tender, *custom_errors.AppError) {
	// Validate that the tender belongs to the client
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("Tender not found or access denied")
		}
		return nil, custom_errors.NewAppError(err)
	}

	// Validate the update request
	if err := ValidateTenderUpdate(tender.Status, req.Status); err != nil {
		return nil, err
	}

	tender.Status = req.Status

	// Save the updated tender to the database
	if err := t.db.Save(&tender).Error; err != nil {
		return nil, custom_errors.NewAppError(err)
	}

	// Invalidate the cache after updating the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return &tender, nil
}

func ValidateTenderUpdate(existingStatus, newStatus string) *custom_errors.AppError {
	// Reject updates if the existing status is not "open"
	if existingStatus != "open" {
		return custom_errors.NewAppError(errors.New("updates are only allowed for tenders with 'open' status"))
	}

	// Reject updates if the new status is "awarded"
	if newStatus == "awarded" {
		return custom_errors.NewAppError(errors.New("status cannot be updated to 'awarded'"))
	}

	return nil
}

// DeleteTender deletes a tender by its ID.
func (t *TenderService) DeleteTender(tenderID, clientID int64) *custom_errors.AppError {
	// Validate that the tender belongs to the client
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	// Perform the deletion
	if err := t.db.Delete(&model.Tender{}, tenderID).Error; err != nil {
		return custom_errors.NewAppError(err)
	}

	// Invalidate the cache after deleting the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return nil
}

// ValidateTenderBelongsToUser ensures that a tender belongs to a specific client.
func (t *TenderService) ValidateTenderBelongsToUser(tenderID, clientID int64) *custom_errors.AppError {
	notFoundError := custom_errors.NewNotFoundError("Tender not found or access denied")

	var tender model.Tender

	// Fetch the tender from the database
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFoundError
		}
		return custom_errors.NewAppError(err)
	}

	// Check if the tender belongs to the client
	if tender.ClientID != clientID {
		return notFoundError
	}

	return nil
}

func (t *TenderService) AwardTender(tenderID, clientID, bidID int64) *custom_errors.AppError {
	// Validate that the tender belongs to the user.
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	// Validate that the bid belongs to the tender.
	if err := t.ValidateBidBelongsToTender(bidID, tenderID); err != nil {
		return err
	}

	// Update the tender status to "awarded", and set the awarded contractor ID.
	if err := t.db.Model(&model.Tender{}).Where("id = ?", tenderID).Updates(map[string]interface{}{
		"status":                "awarded",
		"awarded_contractor_id": bidID,
	}).Error; err != nil {
		return custom_errors.NewAppError(err)
	}

	ntf, _ := t.NotificationService.CreateNotification(&request_model.CreateNotificationReq{
		UserID:  clientID,
		Message: fmt.Sprintf("Your bid for Tender(with id: %s) has been awarded", tenderID),
	})

	if ntf != nil {
		queueReq, _ := proto.Marshal(&gen_proto.Notification{
			Id:      ntf.ID,
			UserId:  ntf.UserID,
			Message: ntf.Message,
		})

		rabbit_mq.Publish("notifications", queueReq)
	}

	return nil
}

func (t *TenderService) ValidateBidBelongsToTender(bidID, tenderID int64) *custom_errors.AppError {
	notFoundError := custom_errors.NewNotFoundError("Bid not found or access denied")

	var bid model.Bid

	if err := t.db.First(&bid, bidID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFoundError
		}
		return custom_errors.NewAppError(err)
	}

	if bid.TenderID != tenderID {
		return notFoundError
	}

	return nil
}
