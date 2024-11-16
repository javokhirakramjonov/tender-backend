package api

import (
	_ "tender-backend/docs"
	"tender-backend/internal/http/handlers"
	"tender-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	// User routes
	userGroup := router.Group("/users").Use(middleware.JWTMiddleware())
	userGroup.GET("", defHandler)
	userGroup.PUT("/:id", defHandler)
	userGroup.DELETE("/:id", defHandler)

	// Tenders routes
	tenderGroup := router.Group("/tenders")
	tenderGroup.Use(middleware.JWTMiddleware())

	tenderGroup.POST("", defHandler)
	tenderGroup.GET("/:id", defHandler)
	tenderGroup.GET("", defHandler)
	tenderGroup.PUT("/:id", defHandler)
	tenderGroup.DELETE("/:id", defHandler)

	// Bids routes
	bidGroup := tenderGroup.Group("/:id/bids")
	bidGroup.POST("", defHandler)
	bidGroup.GET("", defHandler)
	bidGroup.GET("/:id", defHandler)

	// Awards routes
	awardGroup := tenderGroup.Group("/:id/awards")
	awardGroup.POST("", defHandler)

	return router
}
