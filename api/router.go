package api

import (
	"tender-backend/config"
	token "tender-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "tender-backend/docs"
)

// NewGinRouter creates a new gin router with the specified RpcClientManager
// @Title RYG API Gateway
// @Version 1.0
// @Description This is the API Gateway for the RYG project.
// @SecurityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
func NewGinRouter(cfg *config.Config) *gin.Engine {
	
	router := gin.Default()
	swaggerUrl := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, swaggerUrl))

	defHandler := func(c *gin.Context){}

	// Auth routes
	router.POST("/login", )
	router.POST("/register", defHandler)

	// User routes
	userGroup := router.Group("/users").Use(token.JWTMiddleware(cfg))
	userGroup.GET("", defHandler)
	userGroup.PUT("", defHandler)
	userGroup.DELETE("", defHandler)

	// Tenders routes
	tendergroup := router.Group("/tenders")
	tendergroup.Use(token.JWTMiddleware(cfg))

	tendergroup.POST("", defHandler)
	tendergroup.GET("/:challenge_id", defHandler)
	tendergroup.GET("", defHandler)
	tendergroup.PUT("/:id", defHandler)
	tendergroup.DELETE("/:id", defHandler)

	// Bids routes
	bidGroup := tendergroup.Group("/:id/bids")
	bidGroup.POST("", defHandler)
	bidGroup.GET("", defHandler)
	bidGroup.GET("/:id", defHandler)

	// Awards routes
	awardGroup := tendergroup.Group("/:id/awards:bid_id")
	awardGroup.POST("", defHandler)




	return router
}
