package handler

import (
	"terminer/internal/budgetservice/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	var input models.NewTransaction
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	id, err := h.services.Transaction.CreateTransaction(user_id, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	var input models.UpdateTransaction
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Transaction.UpdateTransaction(user_id, input); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	var deleteID models.TransactionID
	if err := c.BindJSON(&deleteID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Transaction.DeleteTransaction(user_id, deleteID.ID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetTransactionsByBudget(c *gin.Context) {
	idParam := c.Param("budgetid")

	budgetID, err := uuid.Parse(idParam)
	if err != nil {
		NewErrorResponse(c, 400, "id must be a valid UUID")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	transactions, err := h.services.Transaction.GetTransactionsByBudget(user_id, budgetID)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":      "ok",
		"transactions": transactions,
	})
}
