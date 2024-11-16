package server

import (
	"errors"
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

func (s *BidService) CreateBid(req *request_model.CreateBidReq) (*model.Bid, error) {
	newBid := model.Bid{
		TenderID:     req.TenderID,
		ContractorID: req.ContractorID,
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

func (s *BidService) GetBidByID(id uint) (*model.Bid, error) {
	var bid model.Bid
	if err := s.db.First(&bid, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, nil
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
