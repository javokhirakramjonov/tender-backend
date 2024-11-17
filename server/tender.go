package server

import (
	"context"
	"encoding/json"
	"errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TenderService struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewTenderService initializes a new TenderService with the database connection.
func NewTenderService(db *gorm.DB, redisClient *redis.Client) *TenderService {
	return &TenderService{
		db:    db,
		redis: redisClient,
	}
}

// CreateTender creates a new tender in the database.
func (t *TenderService) CreateTender(req *request_model.CreateTenderReq) (*model.Tender, error) {
	tender := &model.Tender{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
		Status:      "open",
	}

	// Validate the tender data
	if err := validateCreateTender(req); err != nil {
		return nil, err
	}

	// Save the tender to the database
	if err := t.db.Create(tender).Error; err != nil {
		return nil, err
	}

	// Invalidate the cache after creating a new tender
	t.redis.Del(context.Background(), "tenders_cache")

	return tender, nil
}

// ValidateCreateTender validates the input for creating a tender.
func validateCreateTender(req *request_model.CreateTenderReq) error {
	if req.Deadline.Before(time.Now()) {
		return errors.New("deadline must be a future date and time")
	}

	if req.Budget <= 0 {
		return errors.New("budget must be greater than zero")
	}

	return nil
}

// GetTenderById retrieves a tender by its ID.
func (t *TenderService) GetTenderById(id int64) (*model.Tender, error) {
	var tender model.Tender

	// Try fetching the tender from the database
	if err := t.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
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
func (t *TenderService) UpdateTender(tenderID, clientID int64, req *request_model.UpdateTenderReq) (*model.Tender, error) {
	// Validate that the tender belongs to the client
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	// Validate the update request
	if err := ValidateTenderUpdate(tender.Status, req.Status, req.Deadline, req.Budget); err != nil {
		return nil, err
	}

	// Update the tender fields
	tender.Title = req.Title
	tender.Description = req.Description
	tender.Deadline = req.Deadline
	tender.Budget = req.Budget
	tender.Status = req.Status

	// Save the updated tender to the database
	if err := t.db.Save(&tender).Error; err != nil {
		return nil, err
	}

	// Invalidate the cache after updating the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return &tender, nil
}

func ValidateTenderUpdate(existingStatus, newStatus string, deadline time.Time, budget float64) error {
	// Reject updates if the existing status is not "open"
	if existingStatus != "open" {
		return errors.New("updates are only allowed for tenders with 'open' status")
	}

	// Reject updates if the new status is "awarded"
	if newStatus == "awarded" {
		return errors.New("status cannot be updated to 'awarded'")
	}

	// Ensure the deadline is a future date
	if !deadline.After(time.Now()) {
		return errors.New("deadline must be a future date and time")
	}

	// Ensure the budget is greater than 0
	if budget <= 0 {
		return errors.New("budget must be greater than zero")
	}

	return nil
}

// DeleteTender deletes a tender by its ID.
func (t *TenderService) DeleteTender(tenderID, clientID int64) error {
	// Validate that the tender belongs to the client
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	// Perform the deletion
	if err := t.db.Delete(&model.Tender{}, tenderID).Error; err != nil {
		return err
	}

	// Invalidate the cache after deleting the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return nil
}

// ValidateTenderBelongsToUser ensures that a tender belongs to a specific client.
func (t *TenderService) ValidateTenderBelongsToUser(tenderID, clientID int64) error {
	var tender model.Tender

	// Fetch the tender from the database
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tender not found")
		}
		return err
	}

	// Check if the tender belongs to the client
	if tender.ClientID != clientID {
		return errors.New("tender does not belong to the user")
	}

	return nil
}
