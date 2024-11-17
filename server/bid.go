package server

import (
	"errors"
	"fmt"
	"tender-backend/custom_errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"

	"gorm.io/gorm"
)

type BidService struct {
	db            *gorm.DB
	tenderService *TenderService
}

func NewBidService(db *gorm.DB) *BidService {
	return &BidService{
		db:            db,
		tenderService: NewTenderService(db),
	}
}

func (s *BidService) CreateBid(req *request_model.CreateBidReq, tenderID int64, contractorID int64) (*model.Bid, *custom_errors.AppError) {
	if err := s.validateCreateBidRequest(req); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := s.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, custom_errors.NewNotFoundError("Tender not found")
		}
		return nil, custom_errors.NewAppError(err)
	}

	if tender.Status != "open" {
		return nil, custom_errors.NewBadRequestError("Tender is not open for bids")
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
		return nil, custom_errors.NewAppError(err)
	}

	return &newBid, nil
}

func (s *BidService) GetBidByID(bidID, tenderID int64) (*model.Bid, error) {
	var bid model.Bid
	if err := s.db.Where("id = ? AND tender_id = ?", bidID, tenderID).First(&bid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, fmt.Errorf("failed to retrieve bid: %s", err.Error())
	}

	return &bid, nil
}

func (s *BidService) GetAllBids(tenderID int64) ([]model.Bid, *custom_errors.AppError) {
	_, err := s.tenderService.GetTenderById(tenderID)
	if err != nil {
		return nil, err
	}

	var bids []model.Bid
	if err := s.db.Where("tender_id = ?", tenderID).Find(&bids).Error; err != nil {
		return nil, custom_errors.NewAppError(err)
	}

	return bids, nil
}

func (s *BidService) validateCreateBidRequest(req *request_model.CreateBidReq) *custom_errors.AppError {
	err := custom_errors.NewBadRequestError("Invalid bid data")

	if req.DeliveryTime <= 0 {
		return err
	}

	if req.Price <= 0 {
		return err
	}

	return nil
}

func (s *BidService) GetContractorBids(contractorID int64) ([]model.Bid, error) {
	var bids []model.Bid
	if err := s.db.Where("contractor_id = ?", contractorID).Find(&bids).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve bids: %s", err.Error())
	}

	return bids, nil
}

func (s *BidService) DeleteBid(bidID, contractorID int64) *custom_errors.AppError {
	var bid model.Bid
	if err := s.db.Where("id = ? AND contractor_id = ?", bidID, contractorID).First(&bid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return custom_errors.NewNotFoundError("Bid not found or access denied")
		}
		return custom_errors.NewAppError(err)
	}

	if err := s.db.Delete(&bid).Error; err != nil {
		return custom_errors.NewAppError(err)
	}

	return nil
}
