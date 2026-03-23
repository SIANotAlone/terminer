package handler

import (
	"net/http"
	"terminer/internal/budgetservice/models"
	"terminer/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateGoal(c *gin.Context) {
	var input models.NewGoal
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("CreateGoal: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid input")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := h.services.Goal.CreateGoal(user_id, input)
	if err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
			"name":    input.TargetName,
		}).WithError(err).Error("failed to create goal")
		NewErrorResponse(c, 500, "could not create goal")
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id": user_id,
		"goal_id": id,
	}).Info("goal created successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateGoal(c *gin.Context) {
	var input models.UpdateGoal
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("UpdateGoal: bind error: %v", err)
		NewErrorResponse(c, 400, "invalid input")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Goal.UpdateGoal(user_id, input); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
			"name":    input.TargetName,
		}).WithError(err).Error("failed to update goal")
		NewErrorResponse(c, 500, "could not update goal")
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id": user_id,
		"goal_id": input.ID,
	}).Info("goal created successfully")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteGoal(c *gin.Context) {
	var deleteID models.GoalID
	if err := c.BindJSON(&deleteID); err != nil {
		h.logger.Debugf("DeleteGoal: bad input: %v", err)
		NewErrorResponse(c, 400, "invalid id format")
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.services.Goal.DeleteGoal(user_id, deleteID.ID); err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
			"goal_id": deleteID.ID,
		}).WithError(err).Error("delete goal failed")
		NewErrorResponse(c, 500, "internal error")
		return
	}

	h.logger.WithFields(logger.Fields{
		"user_id": user_id,
		"goal_id": deleteID.ID,
	}).Info("goal deleted")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetAvailableGoals(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	goals, err := h.services.Goal.GetAvailableGoals(user_id)
	if err != nil {

		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
		}).WithError(err).Error("failed to fetch goals")

		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, goals)
}

func (h *Handler) GetGoalsTransactions(c *gin.Context) {
	Id_param := c.Param("goalid")
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}
	var GoalID uuid.UUID
	GoalID, err = uuid.Parse(Id_param)
	if err != nil {
		h.logger.Debugf("GetGoalsTransactions: invalid uuid: %s", Id_param)
		NewErrorResponse(c, 400, "invalid goal id")
		return
	}
	transactions, err := h.services.Goal.GetGoalsTransactions(user_id, GoalID)
	if err != nil {
		h.logger.WithFields(logger.Fields{
			"user_id": user_id,
			"goal_id": GoalID,
		}).WithError(err).Error("failed to fetch transactions for goal")
		NewErrorResponse(c, 500, "database error")
		return
	}

	c.JSON(200, transactions)
}
