package handlers

import (
	"strconv"
	request_model "tender-backend/model/request"

	"github.com/gin-gonic/gin"
)

// CreateTender godoc
// @Security BearerAuth
// @Summary Create a new tender
// @Description Create a new tender
// @Tags tender
// @Accept json
// @Produce json
// @Param tender body request_model.CreateTenderReq true "Tender information"
// @Success 201 {object} model.Tender
// @Router /tenders [post]
func (h *HTTPHandler) CreateTender(ctx *gin.Context) {
	req := request_model.CreateTenderReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.TenderService.CreateTender(&req)

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, res)
}

// GetTender godoc
// @Security BearerAuth
// @Summary Get a tender by ID
// @Description Get a tender by ID
// @Tags tender
// @Produce json
// @Param tender_id path int true "Tender ID"
// @Success 200 {object} model.Tender
// @Router /tenders/{tender_id} [get]
func (h *HTTPHandler) GetTender(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid tender ID"})
		return
	}

	res, err := h.TenderService.GetTenderById(int64(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// GetTenders godoc
// @Security BearerAuth
// @Summary Get all tenders
// @Description Get all tenders
// @Tags tender
// @Produce json
// @Success 200 {object} []model.Tender
// @Router /tenders [get]
func (h *HTTPHandler) GetTenders(ctx *gin.Context) {
	res, err := h.TenderService.GetTenders()

	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// UpdateTender godoc
// @Security BearerAuth
// @Summary Update a tender by ID
// @Description Update a tender by ID
// @Tags tender
// @Accept json
// @Produce json
// @Param tender_id path int true "Tender ID"
// @Param tender body request_model.UpdateTenderReq true "Tender information"
// @Success 200 {object} model.Tender
// @Router /tenders/{tender_id} [put]
func (h *HTTPHandler) UpdateTender(ctx *gin.Context) {
	// Get tender ID from the path
	tenderID, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid tender ID"})
		return
	}

	userID  := ctx.GetInt64("userID")

	// Bind the request JSON to the UpdateTenderReq struct
	req := request_model.UpdateTenderReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Call the service method
	res, err := h.TenderService.UpdateTender(int64(tenderID), userID, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated tender
	ctx.JSON(200, res)
}


// DeleteTender godoc
// @Security BearerAuth
// @Summary Delete a tender by ID
// @Description Delete a tender by ID
// @Tags tender
// @Param tender_id path int true "Tender ID"
// @Success 204
// @Router /tenders/{tender_id} [delete]
func (h *HTTPHandler) DeleteTender(ctx *gin.Context) {
	// Get tender ID from the path
	tenderID, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid tender ID"})
		return
	}

	userID := ctx.GetInt64("userID")
	err = h.TenderService.DeleteTender(int64(tenderID), userID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with no content
	ctx.Status(204)
}
