package handlers

import (
	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc *services.BookingService
}

func NewBookingHandler(svc *services.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := uidVal.(int)

	booking, err := h.svc.CreateBooking(tenantID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) MyBookings(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := uidVal.(int)

	items, err := h.svc.ListMyBookings(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BookingHandler) OwnerBookings(c *gin.Context) {
	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := uidVal.(int)

	items, err := h.svc.ListOwnerBookings(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load owner bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := uidVal.(int)

	err = h.svc.CancelBooking(id, tenantID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "you can cancel only your own booking"})
		case errors.Is(err, services.ErrAlreadyStarted):
			c.JSON(http.StatusBadRequest, gin.H{"error": "booking already started"})
		case errors.Is(err, services.ErrWrongStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": "booking cannot be cancelled in this status"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel booking"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled",
	})
}

func (h *BookingHandler) ApproveBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := uidVal.(int)

	err = h.svc.ApproveBooking(id, ownerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "this booking does not belong to your spaces"})
		case errors.Is(err, services.ErrWrongStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": "only pending bookings can be approved"})
		case errors.Is(err, services.ErrOverlappingBooking):
			c.JSON(http.StatusConflict, gin.H{"error": "booking overlaps with existing approved booking"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve booking"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking approved",
	})
}

func (h *BookingHandler) RejectBooking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	uidVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ownerID := uidVal.(int)

	err = h.svc.RejectBooking(id, ownerID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "this booking does not belong to your spaces"})
		case errors.Is(err, services.ErrWrongStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": "only pending bookings can be rejected"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject booking"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking rejected",
	})
}
