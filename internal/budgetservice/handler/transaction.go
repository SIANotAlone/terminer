package handler

import (
	"net/http"
	"terminer/internal/budgetservice/models"
	"terminer/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	var input models.NewTransaction

	// ВАЖНО: Заменяем c.BindJSON(&input) на ShouldBindBodyWith
	// binding.JSON нужно импортировать из "github.com/gin-gonic/gin/binding"
	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		h.logger.Debugf("bind json error: %v", err)
		NewErrorResponse(c, 400, "invalid input data")
		return
	}

	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := h.services.Transaction.CreateTransaction(user_id, input)
	if err != nil {
		// Логируем с контекстом, чтобы понять ПОЧЕМУ упало
		h.logger.WithFields(logger.Fields{
			"user":      user_id,
			"budget_id": input.BudgetID,
		}).WithError(err).Error("transaction creation failed")

		NewErrorResponse(c, 500, "could not create transaction")
		return
	}

	// 3. Успех - важная веха для бизнеса
	h.logger.WithFields(logger.Fields{
		"user_id":        user_id,
		"transaction_id": id,
	}).Info("transaction created successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	var input models.UpdateTransaction
	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		h.logger.Debugf("UpdateTransaction: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid request body")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Transaction.UpdateTransaction(user_id, input); err != nil {
		// Логируем ошибку с деталями, чтобы можно было разобраться в причинах
		h.logger.WithFields(logger.Fields{
			"user_id":        user_id,
			"transaction_id": input.TransactionID,
		}).WithError(err).Error("failed to update transaction in database")

		NewErrorResponse(c, 500, "could not update transaction")
		return
	}

	// Лог успеха для аудита (audit log)
	h.logger.WithFields(logger.Fields{
		"user_id":        user_id,
		"transaction_id": input.TransactionID,
		"amount":         input.Amount,
	}).Info("transaction updated successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	// Читаем ID из URL: /api/transactions/:id
	idStr := c.Param("id")
	transactionID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Debugf("delete request with bad uuid: %s", idStr)
		NewErrorResponse(c, 400, "invalid transaction id format")
		return
	}

	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Transaction.DeleteTransaction(user_id, transactionID); err != nil {
		// Тут можно сделать проверку на "Not Found"
		h.logger.WithFields(logger.Fields{
			"user":           user_id,
			"transaction_id": transactionID,
		}).WithError(err).Error("delete operation failed")

		NewErrorResponse(c, 500, "failed to delete")
		return
	}

	// Фиксируем успех
	h.logger.WithFields(logger.Fields{
		"user_id":        user_id,
		"transaction_id": transactionID,
	}).Info("transaction deleted")
	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetTransactionsByBudget(c *gin.Context) {
	idParam := c.Param("budgetid")

	budgetID, err := uuid.Parse(idParam)
	if err != nil {
		// Минимальный лог для отладки фронта
		h.logger.Debugf("GetTransactions: invalid budget uuid: %s", idParam)
		NewErrorResponse(c, 400, "id must be a valid UUID")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	transactions, err := h.services.Transaction.GetTransactionsByBudget(user_id, budgetID)
	if err != nil {
		// Критично: логируем детали сбоя
		h.logger.WithFields(logger.Fields{
			"user":      user_id,
			"budget_id": budgetID,
		}).WithError(err).Error("database query failed")

		NewErrorResponse(c, 500, "could not retrieve transactions")
		return
	}
	c.JSON(200, map[string]interface{}{
		"message":      "ok",
		"transactions": transactions,
	})
}
