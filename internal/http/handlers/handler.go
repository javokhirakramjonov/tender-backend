package handlers

import (
	"tender-backend/server"

	"gorm.io/gorm"
)

type HTTPHandler struct {
	UserService *server.UserService
	TenderService *server.TenderService
}

func NewHttpHandler(db *gorm.DB) *HTTPHandler {
	return &HTTPHandler{
		UserService: server.NewUserService(db),
		TenderService: server.NewTenderService(db),

	}
}
