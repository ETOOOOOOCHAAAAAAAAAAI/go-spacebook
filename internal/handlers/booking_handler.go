package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	booking, err := h.bookingService.CreateBooking(userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create booking: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) MyBookings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	bookings, err := h.bookingService.ListMyBookings(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch bookings",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) OwnerBookings(c *gin.Context) {
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	bookings, err := h.bookingService.ListOwnerBookings(ownerID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch bookings",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid booking ID",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var req domain.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		// Причина опциональна, поэтому не требуем её
		req.Reason = nil
	}

	if err := h.bookingService.CancelBooking(bookingID, userID.(int), req.Reason); err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to cancel this booking",
			})
		case services.ErrAlreadyStarted:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot cancel booking that has already started",
			})
		case services.ErrWrongStatus:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot cancel booking with current status",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to cancel booking: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Booking cancelled successfully",
	})
}

func (h *BookingHandler) ApproveBooking(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid booking ID",
		})
		return
	}

	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var req domain.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		// Причина опциональна
		req.Reason = nil
	}

	if err := h.bookingService.ApproveBooking(bookingID, ownerID.(int), req.Reason); err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to approve this booking",
			})
		case services.ErrWrongStatus:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot approve booking with current status",
			})
		case services.ErrOverlappingBooking:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Cannot approve booking due to overlapping with another approved booking",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to approve booking: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Booking approved successfully",
	})
}

func (h *BookingHandler) RejectBooking(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid booking ID",
		})
		return
	}

	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	var req domain.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		// Причина опциональна
		req.Reason = nil
	}

	if err := h.bookingService.RejectBooking(bookingID, ownerID.(int), req.Reason); err != nil {
		switch err {
		case services.ErrForbidden:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to reject this booking",
			})
		case services.ErrWrongStatus:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot reject booking with current status",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to reject booking: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Booking rejected successfully",
	})
}

func (h *BookingHandler) GetBookingHistory(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid booking ID",
		})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	roleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User role not found",
		})
		return
	}

	userID := userIDRaw.(int)
	userRole := domain.UserRole(roleRaw.(string))

	history, err := h.bookingService.GetBookingHistory(bookingID, userID, userRole)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have access to this booking",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get booking history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking_id": bookingID,
		"history":    history,
		"count":      len(history),
	})
}
