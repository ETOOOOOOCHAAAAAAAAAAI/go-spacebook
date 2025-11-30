package handlers

import (
	"SpaceBookProject/internal/repository"
	"net/http"
	"strconv"

	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	svc *services.SpaceService
}

func NewSpaceHandler(svc *services.SpaceService) *SpaceHandler {
	return &SpaceHandler{svc: svc}
}

func (h *SpaceHandler) ListSpaces(c *gin.Context) {
	q := c.Query("q")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	minAreaStr := c.Query("min_area")
	maxAreaStr := c.Query("max_area")

	var f repository.SpaceFilter

	if q != "" {
		f.Query = &q
	}
	if minPriceStr != "" {
		if v, err := strconv.Atoi(minPriceStr); err == nil {
			f.MinPrice = &v
		}
	}
	if maxPriceStr != "" {
		if v, err := strconv.Atoi(maxPriceStr); err == nil {
			f.MaxPrice = &v
		}
	}
	if minAreaStr != "" {
		if v, err := strconv.ParseFloat(minAreaStr, 64); err == nil {
			f.MinArea = &v
		}
	}
	if maxAreaStr != "" {
		if v, err := strconv.ParseFloat(maxAreaStr, 64); err == nil {
			f.MaxArea = &v
		}
	}

	spaces, err := h.svc.ListSpaces(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load spaces"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": spaces})
}

func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var req domain.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	rawID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	ownerID, ok := rawID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}

	space, err := h.svc.CreateSpace(ownerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create space"})
		return
	}

	c.JSON(http.StatusCreated, space)
}
