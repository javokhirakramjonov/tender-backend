package handlers

import (
	"fmt"
	"tender-backend/config"
	"tender-backend/internal/http/token"
	request_model "tender-backend/model/request"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags 01-Authentication
// @Accept json
// @Produce json
// @Param user body pb.UserCreateReqForSwagger true "User registration request"
// @Success 201 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /register [post]
func (h *HTTPHandler) Register(c *gin.Context) {
	var req request_model.CreateUserReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if !config.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	hashedPassword, err := config.HashPassword(req.Password)
	if err != nil {
		fmt.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}

	req.Password = hashedPassword
	user, err := h.UserService.CreateUser(&req)
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}

	tokens := token.GenerateJWTToken(config.GlobalConfig, int64(user.ID))

	fmt.Println("New account registered to the system: ", req.Email)
	c.JSON(http.StatusCreated, tokens)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags 01-Authentication
// @Accept json
// @Produce json
// @Param credentials body pb.LoginReq true "User login credentials"
// @Success 200 {object} token.Tokens "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid email or password"
// @Router /login [post]
func (h *HTTPHandler) Login(c *gin.Context) {
	req := request_model.LoginUserReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	user, err := h.UserService.GetUserByEmail(req.Email)
	if err != nil {
		fmt.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User registered with this email not found"})
		return
	}

	if !config.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := token.GenerateJWTToken(config.GlobalConfig, int64(user.ID))

	c.JSON(http.StatusOK, token)
}
