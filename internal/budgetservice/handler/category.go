package handler

import (
	"terminer/internal/budgetservice/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var input models.NewCategory
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	id, err := h.services.Category.CreateCategory(user_id, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var input models.UpdateCategory
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	err = h.services.Category.UpdateCategory(user_id, input)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	var deleteID models.CategoryID
	if err := c.BindJSON(&deleteID); err != nil {
		NewErrorResponse(c, 400, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	err = h.services.Category.DeleteCategory(user_id, deleteID.ID)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetAvaliableCategories(c *gin.Context) {
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	categories, err := h.services.Category.GetAvaliableCategories(user_id)
	if err != nil {
		NewErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, map[string]interface{}{
		"message":    "ok",
		"categories": categories,
	})
}
