package api

import (
	_ "tender-backend/docs"
	"tender-backend/internal/http/handlers"
	"tender-backend/internal/http/middleware"
	"time"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @tag.name Authentication
// @tag.description User registration and login methods

// @tag.name Tender
// @tag.description Tender CRUDs

// @tag.name Bid
// @tag.description Bid methods

// NewGinRouter godoc
// @Title Tender API Gateway
// @Version 1.0
// @Description This is the API Gateway for the Tender project.
// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
func NewGinRouter(h *handlers.HTTPHandler) *gin.Engine {
	router := gin.Default()

	swaggerUrl := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, swaggerUrl))

	// Authentication Routes (Public)
	authGroup := router.Group("/")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/register", h.Register)
		authGroup.GET("/users/:user_id", h.GetUserByID)
	}

	// User Routes (Protected with JWT Middleware)
	userGroup := router.Group("/users").Use(middleware.JWTMiddleware())
	{
		userGroup.PUT("", h.UpdateUser)
		userGroup.DELETE("", h.DeleteUser)
	}

	// Tender Routes
	tenderGroup := router.Group("/tenders")
	{
		tenderGroup.GET("/:tender_id", h.GetTender)
		tenderGroup.GET("", h.GetTenders)

		protectedTenderGroup := tenderGroup.Use(middleware.JWTMiddleware(), middleware.ClientMiddleware())
		protectedTenderGroup.POST("", h.CreateTender)
		protectedTenderGroup.PUT("/:tender_id", h.UpdateTender)
		protectedTenderGroup.DELETE("/:tender_id", h.DeleteTender)
	}

	// Bid Routes
	bidGroup := router.Group("/tenders/:tender_id/bids")
	{
		bidSubmissionRateLimit := middleware.RateLimitMiddleware(
			2,           // Max 5 requests
			time.Minute, // Per minute
		)

		bidGroup.GET("/:bid_id", h.GetBid)
		bidGroup.GET("", h.GetBids)

		protectedBidGroup := bidGroup.Use(middleware.JWTMiddleware(), middleware.ContractorMiddleware())
		protectedBidGroup.POST("", bidSubmissionRateLimit, h.CreateBid)
	}

	return router
}
