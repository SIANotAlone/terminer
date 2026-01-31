package handler

import (
	"terminer/internal/budgetservice/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBudget(c *gin.Context) {
	var input models.NewBudget
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	id, err := h.services.Budget.CreateBudget(user_id, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	var input models.UpdateBudget
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Budget.UpdateBudget(user_id, input); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	var deleteID models.BudgetID
	if err := c.BindJSON(&deleteID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Budget.DeleteBudget(user_id, deleteID.ID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
}

func (h *Handler) GetAvailableBudgets(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	budgets, err := h.services.Budget.GetAvailableBudgets(user_id)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, budgets)
}

func (h *Handler) ArchiveBudget(c *gin.Context) {
	var budgetID models.BudgetID
	if err := c.BindJSON(&budgetID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Budget.ArchiveBudget(user_id, budgetID.ID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) UnArchiveBudget(c *gin.Context) {
	var budgetID models.BudgetID
	if err := c.BindJSON(&budgetID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Budget.UnArchiveBudget(user_id, budgetID.ID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetBudgetTypes(c *gin.Context) {
	budgetTypes, err := h.services.Budget.GetBudgetTypes()
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, budgetTypes)
}
