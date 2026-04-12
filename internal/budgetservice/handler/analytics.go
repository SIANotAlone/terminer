package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"terminer/pkg/logger"
)

func (h *Handler) GetDashboardData(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}
	idParam := c.Param("budgetid")

	budgetID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Debugf("GetDashboardData: invalid budget uuid: %s", idParam)
		NewErrorResponse(c, 400, "id must be a valid UUID")
		return
	}

	data, err := h.services.Analytics.GetDashboardData(budgetID, user_id)
	if err != nil {
		// Логируем со всеми контекстными полями для дебага
		h.logger.WithFields(logger.Fields{
			"user_id":   user_id,
			"budget_id": budgetID,
		}).WithError(err).Error("failed to retrieve dashboard data")

		NewErrorResponse(c, 500, "could not load analytics")
		return
	}
	c.JSON(200, data)
}
