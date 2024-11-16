package request_model

import (
	"tender-backend/model"
	"time"
)

type CreateUserReq struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type CreateBidReq struct {
	TenderID     int64     `json:"tender_id"`
	ContractorID int64     `json:"contractor_id"`
	Price        float64   `json:"price"`
	DeliveryTime time.Time `json:"delivery_time"`
	Comments     string    `json:"comments"`
	Status       string    `json:"status"`
}

type GetAllBidsRes struct {
	Bids []model.Bid `json:"bids"`
}
