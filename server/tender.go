package server

import (
	"errors"
	"tender-backend/custom_errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"time"

	"gorm.io/gorm"
)

type TenderService struct {
	db *gorm.DB
}

// NewTenderService initializes a new TenderService with the database connection.
func NewTenderService(db *gorm.DB) *TenderService {
	return &TenderService{
		db: db,
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
		Deadline:    req.Deadline, // Use the time.Time value directly.
		Budget:      req.Budget,
		Status:      "open",
	}

	// Save the tender to the database.
	if err := t.db.Create(tender).Error; err != nil {
		return nil, custom_errors.NewAppError(err)
	}

	return tender, nil
}

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

	if err := t.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("Tender not found or access denied")
		}
		return nil, custom_errors.NewAppError(err)
	}

	return &tender, nil
}

// GetTenders retrieves all tenders from the database.
func (t *TenderService) GetTenders() ([]model.Tender, error) {
	var tenders []model.Tender

	if err := t.db.Find(&tenders).Error; err != nil {
		return nil, err
	}

	return tenders, nil
}

func (t *TenderService) UpdateTender(tenderID, clientID int64, req *request_model.UpdateTenderReq) (*model.Tender, *custom_errors.AppError) {
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

	if err := ValidateTenderUpdate(tender.Status, req.Status); err != nil {
		return nil, err
	}

	tender.Status = req.Status

	if err := t.db.Save(&tender).Error; err != nil {
		return nil, custom_errors.NewAppError(err)
	}

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
	// Validate that the tender belongs to the user.
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	// Perform the deletion.
	if err := t.db.Delete(&model.Tender{}, tenderID).Error; err != nil {
		return custom_errors.NewAppError(err)
	}

	return nil
}

// ValidateTenderBelongsToUser ensures that a tender belongs to a specific client.
func (t *TenderService) ValidateTenderBelongsToUser(tenderID, clientID int64) *custom_errors.AppError {
	notFoundError := custom_errors.NewNotFoundError("Tender not found or access denied")

	var tender model.Tender

	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFoundError
		}
		return custom_errors.NewAppError(err)
	}

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
