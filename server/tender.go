package server

import (
	"errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"

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
	}

	// Save the tender to the database.
	if err := t.db.Create(tender).Error; err != nil {
		return nil, err
	}

	return tender, nil
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
	// Validate that the tender belongs to the user.
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

	// Update the fields with the new data.
	if req.Title != "" {
		tender.Title = req.Title
	}
	if req.Description != "" {
		tender.Description = req.Description
	}
	if !req.Deadline.IsZero() { // Check if the deadline is set (not the zero value of time.Time).
		tender.Deadline = req.Deadline
	}
	if req.Budget > 0 {
		tender.Budget = req.Budget
	}
	if req.Status != "" {
		tender.Status = req.Status
	}

	// Save the updated tender to the database.
	if err := t.db.Save(&tender).Error; err != nil {
		return nil, err
	}

	return &tender, nil
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
