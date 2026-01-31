package handler

import (
	"terminer/internal/budgetservice/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateGoal(c *gin.Context) {
	var input models.NewGoal
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	id, err := h.services.Goal.CreateGoal(user_id, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateGoal(c *gin.Context) {
	var input models.UpdateGoal
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Goal.UpdateGoal(user_id, input); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteGoal(c *gin.Context) {
	var deleteID models.GoalID
	if err := c.BindJSON(&deleteID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Goal.DeleteGoal(user_id, deleteID.ID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetAvailableGoals(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	goals, err := h.services.Goal.GetAvailableGoals(user_id)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, goals)
}
