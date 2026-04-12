package handler

import (
	"net/http"
	"strconv"
	"terminer/internal/budgetservice/models"
	"terminer/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBudget(c *gin.Context) {
	var input models.NewBudget
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("CreateBudget: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid input")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := h.services.Budget.CreateBudget(user_id, input)
	if err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
			"name":    input.Name,
		}).WithError(err).Error("failed to create budget")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": id,
	}).Info("budget created successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	var input models.UpdateBudget
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("UpdateBudget: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid input")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Budget.UpdateBudget(user_id, input); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id":   user_id,
			"budget_id": input.ID,
		}).WithError(err).Error("failed to update budget")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": input.ID,
	}).Info("budget updated successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	var deleteID models.BudgetID
	if err := c.BindJSON(&deleteID); err != nil {
		h.logger.Debugf("DeleteBudget: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid id")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Budget.DeleteBudget(user_id, deleteID.ID); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id":   user_id,
			"budget_id": deleteID.ID,
		}).WithError(err).Error("failed to delete budget")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": deleteID.ID,
	}).Info("budget archived")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetAvailableBudgets(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.Error("DeleteBudget: user context missing")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	// Читаем параметры из запроса
	// Пример: /api/budget?limit=10&offset=0&archived=true
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	archived, _ := strconv.ParseBool(c.DefaultQuery("archived", "false"))

	// Вызываем обновленный сервис
	budgets, err := h.services.Budget.GetAvailableBudgets(user_id, archived, limit, offset)
	if err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id":  user_id,
			"limit":    limit,
			"offset":   offset,
			"archived": archived,
		}).WithError(err).Error("failed to fetch budgets")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, budgets)
}

func (h *Handler) ArchiveBudget(c *gin.Context) {
	var budgetID models.BudgetID
	if err := c.BindJSON(&budgetID); err != nil {
		h.logger.Debugf("ArchiveBudget: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid id")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Budget.ArchiveBudget(user_id, budgetID.ID); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id":   user_id,
			"budget_id": budgetID.ID,
		}).WithError(err).Error("failed to archive budget")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": budgetID.ID,
	}).Info("budget archived")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) UnArchiveBudget(c *gin.Context) {
	var budgetID models.BudgetID
	if err := c.BindJSON(&budgetID); err != nil {
		h.logger.Debugf("UnArchiveBudget: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid id")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Budget.UnArchiveBudget(user_id, budgetID.ID); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id":   user_id,
			"budget_id": budgetID.ID,
		}).WithError(err).Error("failed to unarchive budget")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id":   user_id,
		"budget_id": budgetID.ID,
	}).Info("budget unarchived")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetBudgetTypes(c *gin.Context) {
	budgetTypes, err := h.services.Budget.GetBudgetTypes()
	if err != nil {
		h.logger.WithError(err).Error("failed to get budget types")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, budgetTypes)
}

func (h *Handler) GetCurrencies(c *gin.Context) {
	currencies, err := h.services.Budget.GetCurrencies()
	if err != nil {
		h.logger.WithError(err).Error("failed to get currencies")
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, currencies)
}
