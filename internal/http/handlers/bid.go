package handlers

import (
	"net/http"
	"strconv"
	request_model "tender-backend/model/request"

	"github.com/gin-gonic/gin"
)

// CreateBid godoc
// @Summary Create a new bid
// @Description Creates a new bid. Example time: 2024-11-16T15:00:00Z
// @Tags Bid
// @Accept json
// @Produce json
// @Param bid body request_model.CreateBidReq true "Bid creation request"
// @Success 201 {object} model.Bid "Bid created successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /tenders/{tender_id}/bids [POST]
func (h *HTTPHandler) CreateBid(c *gin.Context) {
	var req request_model.CreateBidReq
	tenderIdStr := c.Param("tender_id")
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	contractorId := int64(c.GetInt("user_id"))

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.DeliveryTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery time"})
		return
	}

	createdBid, err := h.BidService.CreateBid(&req, int64(tenderId), contractorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bid"})
		return
	}

	c.JSON(http.StatusCreated, createdBid)
}

// GetBidByID godoc
// @Summary Get Bid by ID
// @Description Retrieves a bid by its ID.
// @Tags Bid
// @Accept json
// @Produce json
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} model.Bid "Bid retrieved successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Bid not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /tenders/{tender_id}/bids/{bid_id} [GET]
func (h *HTTPHandler) GetBid(c *gin.Context) {
	idStr := c.Param("bid_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	bid, err := h.BidService.GetBidByID(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bid"})
		return
	}

	if bid == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// GetAllBids godoc
// @Summary Get all Bids
// @Description Retrieves all Bids for the authenticated user.
// @Tags Bid
// @Accept json
// @Produce json
// @Success 200 {object} []model.Bid "All bids retrieved successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /tenders/{tender_id}/bids [get]
func (h *HTTPHandler) GetBids(c *gin.Context) {
	bids, err := h.BidService.GetAllBids()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}
