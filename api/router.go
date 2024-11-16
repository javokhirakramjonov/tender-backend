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

	defHandler := func(c *gin.Context) {}

	// Auth routes
	router.POST("/login", h.Login)
	router.POST("/register", h.Register)

	// User routes (protected)
	userGroup := router.Group("/users").Use(middleware.JWTMiddleware())
	userGroup.PUT("/:id", h.UpdateUser)
	userGroup.DELETE("/:id", h.DeleteUser)

	// Tenders routes
	tenderGroup := router.Group("/tenders")

	// Unprotected GET routes for tenders
	tenderGroup.GET("/:tender_id", defHandler) // View a specific tender
	tenderGroup.GET("", defHandler)            // List all tenders

	// Protected POST, PUT, DELETE routes for tenders
	protectedTenderGroup := tenderGroup.Use(middleware.JWTMiddleware(), middleware.ContractorMiddleware())
	protectedTenderGroup.POST("", defHandler)
	protectedTenderGroup.PUT("/:tender_id", defHandler)
	protectedTenderGroup.DELETE("/:tender_id", defHandler)

	// Bids routes
	bidGroup := router.Group("/tenders/:tender_id/bids")

	// Unprotected GET routes for bids
	bidGroup.GET("/:bid_id", h.GetBid) // View a specific bid
	bidGroup.GET("", h.GetBids)        // List all bids for a tender

	// Protected POST routes for bids
	protectedBidGroup := bidGroup.Use(middleware.JWTMiddleware(), middleware.ClientMiddleware())
	protectedBidGroup.POST("", h.CreateBid)

	// Awards routes
	awardGroup := tenderGroup.Group("/:tender_id/awards").Use(middleware.JWTMiddleware())
	awardGroup.POST("", defHandler)

	return router
}
