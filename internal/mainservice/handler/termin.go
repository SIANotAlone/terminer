package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Отримання всіх термінів для виконавця
// @Description  Хендлер для отримання всіх термінів для виконавця, що належать користувачу.
// @Tags         Реакції
// @Produce      json
// @Param        userId  header  string  true  "Ідентифікатор користувача"
// @Success      200     {array}  models.Termin  "Список термінів для виконавця"
// @Failure      400     {object} map[string]string      "Невірні дані запиту"
// @Failure      500     {object} map[string]string      "Помилка сервера"
// @Router       /api/termin/getallperformertermins [get]
// @Security     ApiKeyAuth
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

// @Summary      Отримання всіх термінів користувача
// @Description  Хендлер для отримання всіх термінів, що належать користувачу.
// @Tags         Реакції
// @Produce      json
// @Success      200     {array}  models.Termin   "Список термінів користувача"
// @Failure      400     {object} map[string]string   "Невірні дані запиту"
// @Failure      500     {object} map[string]string   "Помилка сервера"
// @Router       /api/termin/getallusertermins [get]
// @Security     ApiKeyAuth
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
