package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateComment(c *gin.Context) {

	var input models.Comment
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.UserID = user_id

	id, err := h.services.Comment.CreateComment(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"id":      id,
	})
}

func (h *Handler) UpdateComment(c *gin.Context) {
	var input models.UpdateComment
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.UserID = user_id
	if err := h.services.Comment.UpdateComment(input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteComment(c *gin.Context) {
	var commentID models.DeleteComment
	if err := c.BindJSON(&commentID); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Comment.DeleteComment(commentID.ID, user_id); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetCommentsOnRecord(c *gin.Context) {
	var input models.GetCommentsInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	comments, err := h.services.Comment.GetCommentsOnRecord(input.RecordID, user_id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, comments)
}
