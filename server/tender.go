package server

import (
	"errors"
	"time"
	"tender-backend/model"
	request_model "tender-backend/model/request"

	"gorm.io/gorm"
)

type TenderService struct {
	db *gorm.DB
}

func NewTenderService(db *gorm.DB) *TenderService {
	return &TenderService{
		db: db,
	}
}

// Create a new tender
func (s *TenderService) CreateTender(req *request_model.CreateTenderRequest) (*model.Tender, error) {
	// Parse the deadline string to a time.Time object
	deadline, err := time.Parse("2006-01-02", req.Deadline) // Adjust format as needed
	if err != nil {
		return nil, errors.New("invalid deadline format, expected YYYY-MM-DD")
	}

	tender := &model.Tender{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
		Budget:      req.Budget,
		Status:      req.Status,
	}

	// Save the tender to the database
	if err := s.db.Create(tender).Error; err != nil {
		return nil, err
	}
	return tender, nil
}

// Get a tender by its ID
func (s *TenderService) GetTenderById(id int64) (*model.Tender, error) {
	var tender model.Tender
	if err := s.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}
	return &tender, nil
}

// Get all tenders
func (s *TenderService) GetTenders() ([]model.Tender, error) {
	var tenders []model.Tender
	if err := s.db.Find(&tenders).Error; err != nil {
		return nil, err
	}
	return tenders, nil
}

// Update an existing tender by its ID
func (s *TenderService) UpdateTender(id uint, req *request_model.UpdateTenderRequest) (*model.Tender, error) {
	// Fetch the existing tender
	var tender model.Tender
	if err := s.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	// Parse the deadline string to time.Time
	deadline, err := time.Parse("2006-01-02", req.Deadline) // Adjust format as needed
	if err != nil {
		return nil, errors.New("invalid deadline format, expected YYYY-MM-DD")
	}

	// Update the fields with the new data
	tender.Title = req.Title
	tender.Description = req.Description
	tender.Deadline = deadline
	tender.Budget = req.Budget
	tender.Status = req.Status

	// Save the updated tender to the database
	if err := s.db.Save(&tender).Error; err != nil {
		return nil, err
	}

	return &tender, nil
}


// Delete a tender by its ID
func (s *TenderService) DeleteTender(id int64) error {
	if err := s.db.Delete(&model.Tender{}, id).Error; err != nil {
		return err
	}
	return nil
}
