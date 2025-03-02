package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Створення нового запису
// @Description  Хендлер для створення нового запису. Приймає дані для запису та ідентифікує користувача на основі токена.
// @Tags         Запис на послугу
// @Accept       json
// @Produce      json
// @Param        input  body     models.NewRecord  true  "Дані для створення нового запису"
// @Success      200    {object} map[string]interface{}  "Ідентифікатор створеного запису"  {"message": "ok", "record_id": "your_record_id_here"}
// @Failure      400    {object} map[string]string       "Невірні дані запиту"
// @Failure      500    {object} map[string]string       "Помилка сервера"
// @Router       /api/record/create [post]
// @Security     ApiKeyAuth
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

// @Summary      Завершення запису
// @Description  Хендлер для завершення запису. Приймає ID запису та позначає його як завершений для користувача.
// @Tags         Запис на послугу
// @Accept       json
// @Produce      json
// @Param        input  body     models.DoneRecord  true  "Дані для позначення запису як завершеного (ID запису)"
// @Success      200    {object} map[string]interface{}  "Повідомлення про успішне завершення запису"  {"message": "ok"}
// @Failure      400    {object} map[string]string       "Невірні дані запиту"
// @Failure      500    {object} map[string]string       "Помилка сервера"
// @Router       /api/record/done [post]
// @Security     ApiKeyAuth
func (h *Handler) DoneRecord(c *gin.Context) {
	var input models.DoneRecord

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Record.DoneRecord(input.ID, user_id); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})

}

// @Summary      Підтвердження запису
// @Description  Хендлер для підтвердження запису. Приймає ID запису та підтверджує його для користувача.
// @Tags         Запис на послугу
// @Accept       json
// @Produce      json
// @Param        input  body     models.DoneRecord  true  "Дані для підтвердження запису (ID запису)"
// @Success      200    {object} map[string]interface{}  "Повідомлення про успішне підтвердження запису"  {"message": "ok"}
// @Failure      400    {object} map[string]string       "Невірні дані запиту"
// @Failure      500    {object} map[string]string       "Помилка сервера"
// @Router       /api/record/confirm [post]
// @Security     ApiKeyAuth
func (h Handler) ConfirmRecord(c *gin.Context) {
	var input models.DoneRecord
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Record.ConfirmRecord(input.ID, user_id); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
