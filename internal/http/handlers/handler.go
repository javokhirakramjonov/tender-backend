package handlers

import (
	"tender-backend/server"

	"gorm.io/gorm"
)

type HTTPHandler struct {
	UserService *server.UserService
}

func NewHandler(db *gorm.DB) *HTTPHandler {
	return &HTTPHandler{
		UserService: server.NewUserService(db),
	}
}
