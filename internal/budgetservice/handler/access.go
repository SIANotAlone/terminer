package handler

import (
	"net/http"
	"terminer/internal/budgetservice/models"
	"terminer/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ShareBudget(c *gin.Context) {
	var input models.ShareBudgetInput
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("ShareBudget: bind error: %v", err)
        NewErrorResponse(c, http.StatusBadRequest, "invalid input")
        return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}
	var accessID uuid.UUID
	if access_id, err := h.services.Access.ShareBudget(user_id, input.BudgetID, input.TargetUser); err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id":   user_id,
            "budget_id": input.BudgetID,
        }).WithError(err).Error("failed to share budget")
        NewErrorResponse(c, http.StatusInternalServerError, "could not share budget")
        return
	} else {
		accessID = access_id
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": input.BudgetID,
		"access_id": accessID,
	}).Info("budget shared successfully")

	c.JSON(200, map[string]interface{}{
		"message":   "ok",
		"access_id": accessID,
	})

}

func (h *Handler) RevokeAccess(c *gin.Context) {
	var input models.RevokeAccessInput
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("RevokeAccess: bind error: %v", err)
        NewErrorResponse(c, http.StatusBadRequest, "invalid input")
        return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Access.RevokeAccess(user_id, input.AccessID); err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id":   user_id,
            "access_id": input.AccessID,
        }).WithError(err).Error("failed to revoke access")
        NewErrorResponse(c, http.StatusInternalServerError, "failed to revoke access")
        return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"access_id": input.AccessID,
	}).Info("access to budget revoked successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetBudgetAccessList(c *gin.Context) {
	var budgetID models.BudgetID
	paramID := c.Query("budget_id")
	if paramID == "" {
		h.logger.Debugf("GetBudgetAccessList: invalid uuid: %s", paramID)
        NewErrorResponse(c, http.StatusBadRequest, "budget_id must be a valid UUID")
        return
	}

	// 2. Объявляем переменную для ID
	// Используем := для создания новой переменной err
	parsedID, err := uuid.Parse(paramID)
	if err != nil {
		NewErrorResponse(c, 400, "budget_id must be a valid UUID")
		return
	}
	budgetID.ID = parsedID

	user_id, err := getUserId(c)
	if err != nil {
			h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	accessList, err := h.services.Access.GetBudgetAccessList(user_id, budgetID.ID)
	if err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id":   user_id,
            "budget_id": budgetID,
        }).WithError(err).Error("failed to fetch access list")
        NewErrorResponse(c, http.StatusInternalServerError, "failed to load access list")
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
		h.logger.WithError(err).Error("failed to fetch all users")
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
		h.logger.Debugf("GetBudgetUsers: invalid budget uuid: %s", idParam)
		NewErrorResponse(c, 400, "id must be a valid UUID")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	users, err := h.services.Access.GetBudgetAccessList(user_id, budgetID)
	if err != nil {
		h.logger.WithError(err).Error("failed to fetch budget users")
		NewErrorResponse(c, 500, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"users":   users,
	})
}
