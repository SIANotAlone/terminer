package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateRecord(c *gin.Context) {
	var input models.NewRecord
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.UserID = user_id
	var id uuid.UUID

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if id, err = h.services.Record.CreateRecord(input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "ok",
		"record_id": id,
	})

}
