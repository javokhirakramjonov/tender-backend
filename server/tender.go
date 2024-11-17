package server

import (
	"errors"
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
func (t *TenderService) CreateTender(req *request_model.CreateTenderReq) (*model.Tender, error) {
	tender := &model.Tender{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline, // Use the time.Time value directly.
		Budget:      req.Budget,
		Status:      "open",
	}

	if err := validateCreateTender(req); err != nil {
		return nil, err
	}

	// Save the tender to the database.
	if err := t.db.Create(tender).Error; err != nil {
		return nil, err
	}

	return tender, nil
}

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

	if err := t.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
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

func (t *TenderService) UpdateTender(tenderID, clientID int64, req *request_model.UpdateTenderReq) (*model.Tender, error) {
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

	if err := ValidateTenderUpdate(tender.Status, req.Status, req.Deadline, req.Budget); err != nil {
		return nil, err
	}

	tender.Title = req.Title
	tender.Description = req.Description
	tender.Deadline = req.Deadline
	tender.Budget = req.Budget
	tender.Status = req.Status

	if err := t.db.Save(&tender).Error; err != nil {
		return nil, err
	}

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
	// Validate that the tender belongs to the user.
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	// Perform the deletion.
	if err := t.db.Delete(&model.Tender{}, tenderID).Error; err != nil {
		return err
	}

	return nil
}

// ValidateTenderBelongsToUser ensures that a tender belongs to a specific client.
func (t *TenderService) ValidateTenderBelongsToUser(tenderID, clientID int64) error {
	var tender model.Tender

	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tender not found")
		}
		return err
	}

	if tender.ClientID != clientID {
		return errors.New("tender does not belong to the user")
	}

	return nil
}
