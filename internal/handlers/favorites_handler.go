package handlers

import (
	"net/http"
	"strconv"

	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

type FavoritesHandler struct {
	svc *services.FavoritesService
}

func NewFavoritesHandler(svc *services.FavoritesService) *FavoritesHandler {
	return &FavoritesHandler{svc: svc}
}

func getUserID(c *gin.Context) (int, bool) {
	keys := []string{"user_id", "userID", "userId", "id", "user"}
	for _, k := range keys {
		if v, ok := c.Get(k); ok {
			switch t := v.(type) {
			case int:
				return t, true
			case int64:
				return int(t), true
			case float64:
				return int(t), true
			case string:
				n, err := strconv.Atoi(t)
				if err == nil {
					return n, true
				}
			case domain.User:
				return t.ID, true
			case *domain.User:
				return t.ID, true
			}
		}
	}
	return 0, false
}

func (h *FavoritesHandler) AddFavorite(c *gin.Context) {
	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil || spaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.svc.Add(c.Request.Context(), userID, spaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *FavoritesHandler) RemoveFavorite(c *gin.Context) {
	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil || spaceID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid space id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.svc.Remove(c.Request.Context(), userID, spaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *FavoritesHandler) ListFavorites(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	spaces, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list favorites"})
		return
	}

	c.JSON(http.StatusOK, spaces)
}
