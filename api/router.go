package api

import (
	_ "tender-backend/docs"
	"tender-backend/internal/http/handlers"
	"tender-backend/internal/http/middleware"

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

	// Auth routes
	router.POST("/login", h.Login)
	router.POST("/register", h.Register)
	router.GET("/users/:user_id", h.GetUserByID)

	// User routes (protected)
	userGroup := router.Group("/users").Use(middleware.JWTMiddleware())
	userGroup.PUT("", h.UpdateUser)
	userGroup.DELETE("", h.DeleteUser)

	// Tenders routes
	tenderGroup := router.Group("/api/client/tenders")

	// Unprotected GET routes for tenders
	tenderGroup.GET("/:tender_id", h.GetTender) // View a specific tender
	tenderGroup.GET("", h.GetTenders)           // List all tenders

	// Protected POST, PUT, DELETE routes for tenders
	protectedTenderGroup := tenderGroup.Use(middleware.JWTMiddleware(), middleware.ClientMiddleware())
	protectedTenderGroup.POST("", h.CreateTender)
	protectedTenderGroup.PUT("/:tender_id", h.UpdateTender)
	protectedTenderGroup.DELETE("/:tender_id", h.DeleteTender)

	// Bids routes
	bidGroup := router.Group("/api/contractor/tenders/:tender_id/bid")

	// Unprotected GET routes for bids
	bidGroup.GET("/:bid_id", h.GetBid)

	clientBidsGroup := router.Group("/api/client/tenders/:tender_id/bids")
	clientBidsGroup.Use(middleware.JWTMiddleware(), middleware.ClientMiddleware())
	clientBidsGroup.GET("", h.GetBids)

	// Protected POST routes for bids
	protectedBidGroup := bidGroup.Use(middleware.JWTMiddleware(), middleware.ContractorMiddleware())
	protectedBidGroup.POST("", h.CreateBid)

	contractorBidGroup := router.Group("/api/contractor/bids")
	contractorBidGroup.Use(middleware.JWTMiddleware(), middleware.ContractorMiddleware())
	contractorBidGroup.GET("", h.GetContractorBids)
	contractorBidGroup.DELETE("/:bid_id", h.DeleteBid)

	// Awards routes
	awardGroup := tenderGroup.Group("/:tender_id/award")
	awardGroup.POST("/:bid_id", h.AwardTender)

	return router
}
