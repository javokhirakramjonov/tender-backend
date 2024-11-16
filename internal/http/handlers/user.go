package handlers

import (
	"net/http"
	"strconv"
	request_model "tender-backend/model/request"

	"github.com/gin-gonic/gin"
)

// UpdateUser godoc
// @Summary Update user by ID
// @Description Updates a user's information by their ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body request_model.UpdateUserReq true "User update request"
// @Success 200 {object} model.User "User updated successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /users/{id} [PUT]
func (h *HTTPHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req request_model.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedUser, err := h.UserService.UpdateUser(&req, int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Deletes a user by their ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {object} string "User deleted successfully"
// @Failure 400 {object} string "Invalid user ID"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /users/{id} [DELETE]
func (h *HTTPHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UserService.DeleteUser(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
