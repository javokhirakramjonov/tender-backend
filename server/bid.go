package server

import (
	"errors"
	"fmt"
	"tender-backend/model"
	request_model "tender-backend/model/request"

	"gorm.io/gorm"
)

type BidService struct {
	db *gorm.DB
}

func NewBidService(db *gorm.DB) *BidService {
	return &BidService{db: db}
}

func (s *BidService) CreateBid(req *request_model.CreateBidReq, tenderID int64, contractorID int64) (*model.Bid, error) {
	if err := s.validateCreateBidRequest(req); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := s.db.First(&tender, "id = ?", tenderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tender with ID %d not found", tenderID)
		}
		return nil, err
	}

	if tender.Status != "open" {
		return nil, fmt.Errorf("cannot place a bid on a tender that is not open")
	}

	newBid := model.Bid{
		TenderID:     tenderID,
		ContractorID: contractorID,
		Price:        req.Price,
		DeliveryTime: req.DeliveryTime,
		Comments:     req.Comments,
		Status:       "pending",
	}

	if err := s.db.Create(&newBid).Error; err != nil {
		return nil, err
	}

	return &newBid, nil
}

func (s *BidService) GetBidByID(id int64) (*model.Bid, error) {
	var bid model.Bid
	if err := s.db.First(&bid, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, fmt.Errorf("failed to retrieve bid:" + err.Error())
	}

	return &bid, nil
}

func (s *BidService) GetAllBids() ([]model.Bid, error) {
	var bids []model.Bid
	if err := s.db.Find(&bids).Error; err != nil {
		return nil, err
	}

	return bids, nil
}

func (s *BidService) validateCreateBidRequest(req *request_model.CreateBidReq) error {
	if req.DeliveryTime < 0 {
		return errors.New("delivery time must be greater than zero")
	}

	if req.Price < 0 {
		return errors.New("price must be greater than zero")
	}

	return nil
}
