package handlers

import (
	"net/http"
	"strconv"

	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

type BookingHistoryHandler struct {
	service *services.BookingService
}

func NewBookingHistoryHandler(service *services.BookingService) *BookingHistoryHandler {
	return &BookingHistoryHandler{service: service}
}

func (h *BookingHistoryHandler) GetHistory(c *gin.Context) {
	// 1. Получаем userID из контекста (кладёт AuthMiddleware)
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "user not authenticated",
		})
		return
	}

	// 2. Получаем bookingID из URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid booking id",
		})
		return
	}

	// 3. Вызываем сервис
	history, err := h.service.BookingHistory(id, userID.(int))
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "forbidden",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get booking history",
		})
		return
	}

	// 4. Отдаём результат
	c.JSON(http.StatusOK, history)
}
