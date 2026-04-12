package handler

import (
	"net/http"
	"terminer/internal/budgetservice/models"
	"terminer/pkg/logger"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var input models.NewCategory
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("CreateCategory: bind error: %v", err)
        NewErrorResponse(c, 400, "invalid input")
        return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := h.services.Category.CreateCategory(user_id, input)
	if err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id": user_id,
            "name":    input.Name,
        }).WithError(err).Error("database error: failed to create category")
        NewErrorResponse(c, 500, "internal server error")
        return
	}

	h.logger.WithFields(logger.Fields{
        "user_id":     user_id,
        "category_id": id,
    }).Info("category created")	

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var input models.UpdateCategory
	if err := c.BindJSON(&input); err != nil {
		h.logger.Debugf("UpdateCategory: bind error: %v", err)
        NewErrorResponse(c, 400, "invalid input")
        return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	err = h.services.Category.UpdateCategory(user_id, input)
	if err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id":     user_id,
            "category_id": input.ID,
        }).WithError(err).Error("failed to update category")
        NewErrorResponse(c, 500, "internal server error")
        return
	}

	h.logger.WithFields(logger.Fields{
        "user_id":     user_id,
        "category_id": input.ID,
    }).Info("category updated")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	var deleteID models.CategoryID
	if err := c.BindJSON(&deleteID); err != nil {
		h.logger.Debugf("DeleteCategory: bad request: %v", err)
        NewErrorResponse(c, 400, "invalid id format")
        return
	}
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	err = h.services.Category.DeleteCategory(user_id, deleteID.ID)
	if err != nil {
		h.logger.WithFields(logger.Fields{
            "user_id":     user_id,
            "category_id": deleteID.ID,
        }).WithError(err).Error("failed to delete category")
        NewErrorResponse(c, 500, "internal server error")
        return
	}

	h.logger.WithFields(logger.Fields{
        "user_id":     user_id,
        "category_id": deleteID.ID,
    }).Info("category deleted")

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetAvaliableCategories(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		h.logger.WithError(err).Error("CRITICAL: user context lost in handler")
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	categories, err := h.services.Category.GetAvaliableCategories(user_id)
	if err != nil {
		h.logger.WithFields(logger.Fields{"user_id": user_id}).WithError(err).Error("GetCategories: db query failed")
        NewErrorResponse(c, 500, "internal server error")
        return
	}

	c.JSON(200, map[string]interface{}{
		"message":    "ok",
		"categories": categories,
	})
}
