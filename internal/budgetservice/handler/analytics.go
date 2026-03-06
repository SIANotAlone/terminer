package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func (h *Handler) GetDashboardData(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	idParam := c.Param("budgetid")

	budgetID, err := uuid.Parse(idParam)
	if err != nil {
		NewErrorResponse(c, 400, "id must be a valid UUID")
		return
	}

	data, err := h.services.Analytics.GetDashboardData(budgetID, user_id)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, data)
}