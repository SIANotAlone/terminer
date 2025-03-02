package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary      Реєстрація нового користувача
// @Description  Хендлер для реєстрації нового користувача. Приймає дані користувача в форматі JSON та створює нового користувача.
// @Tags         Користувач
// @Accept       json
// @Produce      json
// @Param        input  body     models.UserRegistration  true  "Дані для реєстрації нового користувача"
// @Success      200    {object} map[string]interface{}  "Ідентифікатор нового користувача" {"id": 1}
// @Failure      400    {object} map[string]string       "Помилка запиту (неправильні дані)"
// @Failure      500    {object} map[string]string       "Помилка сервера"
// @Router       /auth/sign-up [post]
// @Security     []
func (h *Handler) SignUp(c *gin.Context) {
	var input models.UserRegistration
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// @Summary      Авторизація користувача
// @Description  Хендлер для авторизації користувача. Приймає email та пароль для генерації токена.
// @Tags         Користувач
// @Accept       json
// @Produce      json
// @Param        input  body     models.UserSignIn  true  "Дані для авторизації користувача (email, пароль)"
// @Success      200    {object} map[string]interface{}  "Токен доступу користувача"  {"token": "your_token_here"}
// @Failure      400    {object} map[string]string       "Невірні дані авторизації"
// @Failure      500    {object} map[string]string       "Помилка сервера"
// @Router       /auth/sign-in [post]
// @Security     []
func (h *Handler) SignIn(c *gin.Context) {
	var input models.UserSignIn
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
