package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateService(c *gin.Context) {
	var input models.NewService
	userId, err := getUserId(c)
	input.Service.PerformerID = userId
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Offering.CreateService(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

func (h *Handler) UpdateService(c *gin.Context) {
	var input models.ServiceUpdate

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Offering.UpdateService(input); err != nil {

	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) DeleteService(c *gin.Context) {
	var input models.ServiceDelete
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Offering.DeleteService(input.UUID); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (h *Handler) GetTypes(c *gin.Context) {

	types, err := h.services.Offering.GetTypes()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, types)

}
