package handlers

import (
	"net/http"

	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	spaceService *services.SpaceService
}

func NewSpaceHandler(spaceService *services.SpaceService) *SpaceHandler {
	return &SpaceHandler{spaceService: spaceService}
}

func (h *SpaceHandler) ListSpaces(c *gin.Context) {
	spaces, err := h.spaceService.ListActiveSpaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list spaces"})
		return
	}
	c.JSON(http.StatusOK, spaces)
}

func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	ownerID := userIDAny.(int)

	var req domain.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}

	space, err := h.spaceService.CreateSpace(ownerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create space"})
		return
	}

	c.JSON(http.StatusCreated, space)
}

func (h *SpaceHandler) ListMySpaces(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not authenticated"})
		return
	}
	ownerID := userIDAny.(int)

	spaces, err := h.spaceService.ListOwnerSpaces(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to list owner spaces"})
		return
	}

	c.JSON(http.StatusOK, spaces)
}
