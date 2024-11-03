package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllPerformerTermins(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	performerTermins, err := h.services.Termin.GetAllPerformerTermins(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, performerTermins)
}

func (h *Handler) GetAllUserTermins(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userTermins, err := h.services.Termin.GetAllUserTermins(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, userTermins)
}
