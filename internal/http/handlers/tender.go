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
// @Param tender body request_model.CreateTenderRequest true "Tender information"
// @Success 201 {object} model.Tender
// @Router /tenders [post]
func (h *HTTPHandler) CreateTender(ctx *gin.Context) {
	req := request_model.CreateTenderRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Tender.CreateTender(&req)

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

	res, err := h.Tender.GetTenderById(int64(id))
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
	res, err := h.Tender.GetTenders()

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
// @Param tender body request_model.UpdateTenderRequest true "Tender information"
// @Success 200 {object} model.Tender
// @Router /tenders/{tender_id} [put]
func (h *HTTPHandler) UpdateTender(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid tender ID"})
		return
	}

	req := request_model.UpdateTenderRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Tender.UpdateTender(uint(id), &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

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
	id, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid tender ID"})
		return
	}

	err = h.Tender.DeleteTender(int64(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(204)
}
