package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Отримання всіх користувачів
// @Description  Хендлер для отримання списку всіх користувачів з бази даних.
// @Tags         Користувачі
// @Produce      json
// @Success      200    {array}   models.User    "Список користувачів"
// @Failure      500    {object} map[string]string "Помилка сервера"
// @Router       /api/user/getallusers [get]
// @Security     ApiKeyAuth
func (h *Handler) GetAllUsers(c *gin.Context) {

	users, err := h.services.User.GetAllUsers()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}
