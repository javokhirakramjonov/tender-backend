package handlers

import (
	"tender-backend/server"

	"gorm.io/gorm"
)

type HTTPHandler struct {
	UserService *server.UserService
	BidService  *server.BidService
}

func NewHttpHandler(db *gorm.DB) *HTTPHandler {
	return &HTTPHandler{
		UserService: server.NewUserService(db),
		BidService:  server.NewBidService(db),
	}
}
