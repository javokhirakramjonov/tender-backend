package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BidService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewBidService(db *gorm.DB, redisClient *redis.Client) *BidService {
	return &BidService{
		db:    db,
		redis: redisClient,
	}
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

	// Clear relevant cache for this tender's bids
	s.clearBidsCache(tenderID)

	return &newBid, nil
}

func (s *BidService) GetBidByID(bidID, tenderID int64) (*model.Bid, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("bid_%d_tender_%d", bidID, tenderID)

	// Check Redis cache
	cachedBid, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit: return cached bid
		var bid model.Bid
		if err := json.Unmarshal([]byte(cachedBid), &bid); err == nil {
			return &bid, nil
		}
	}

	// Cache miss: query from DB
	var bid model.Bid
	if err := s.db.Where("id = ? AND tender_id = ?", bidID, tenderID).First(&bid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, fmt.Errorf("failed to retrieve bid: %s", err.Error())
	}

	// Cache the result
	bidJSON, _ := json.Marshal(bid)
	_ = s.redis.Set(ctx, cacheKey, bidJSON, 10*time.Minute).Err()

	return &bid, nil
}

func (s *BidService) GetAllBids(tenderID int64) ([]model.Bid, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("bids_tender_%d", tenderID)

	// Check Redis cache
	cachedBids, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit: return cached bids
		var bids []model.Bid
		if err := json.Unmarshal([]byte(cachedBids), &bids); err == nil {
			return bids, nil
		}
	}

	// Cache miss: query from DB
	var bids []model.Bid
	if err := s.db.Where("tender_id = ?", tenderID).Find(&bids).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve bids: %s", err.Error())
	}

	// Cache the result
	bidsJSON, _ := json.Marshal(bids)
	_ = s.redis.Set(ctx, cacheKey, bidsJSON, 10*time.Minute).Err()

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

// clearBidsCache clears cached bids for a specific tender.
func (s *BidService) clearBidsCache(tenderID int64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("bids_tender_%d", tenderID)
	_ = s.redis.Del(ctx, cacheKey).Err()
}
