package handler

import (
	"terminer/internal/budgetservice/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ShareBudget(c *gin.Context) {
	var input models.ShareBudgetInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Access.ShareBudget(user_id, input.BudgetID, input.TargetUser); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})

}

func (h *Handler) RevokeAccess(c *gin.Context) {
	var input models.RevokeAccessInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	if err := h.services.Access.RevokeAccess(user_id, input.AccessID); err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetBudgetAccessList(c *gin.Context) {
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

	accessList, err := h.services.Access.GetBudgetAccessList(user_id, budgetID.ID)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":     "ok",
		"access_list": accessList,
	})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.services.Access.GetAllUsers()
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"users":   users,
	})
}

func (h *Handler) GetBudgetUsers(c *gin.Context) {
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

	users, err := h.services.Access.GetBudgetAccessList(user_id, budgetID)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"users":   users,
	})
}
