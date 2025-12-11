package handlers

import (
    "net/http"
    "SpaceBookProject/internal/services"
    "github.com/gin-gonic/gin"
)

type NotificationHandler struct {
    svc services.NotificationService
}

func NewNotificationHandler(svc services.NotificationService) *NotificationHandler {
    return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c *gin.Context) {
    raw, ok := c.Get("userID")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID := raw.(int)

    items, err := h.svc.ListUserNotifications(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notifications"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"items": items})
}
