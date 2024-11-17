package handlers

import (
	"tender-backend/server"

	"github.com/redis/go-redis/v9" // Use v9 Redis package
	"gorm.io/gorm"
)

type HTTPHandler struct {
	UserService        *server.UserService
	BidService         *server.BidService
	TenderService      *server.TenderService
	RedisClient        *redis.Client // v9 Redis client
	NotificationServer *server.Server
}

func NewHttpHandler(db *gorm.DB, RedisClient *redis.Client) *HTTPHandler {
	return &HTTPHandler{
		UserService:        server.NewUserService(db),
		BidService:         server.NewBidService(db, RedisClient),
		TenderService:      server.NewTenderService(db, RedisClient),
		RedisClient:        RedisClient,
		NotificationServer: server.NewNotificationServer(),
	}
}
