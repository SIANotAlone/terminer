package handler

import (
	"fmt"
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateService створює послугу
// @Summary      Створити послугу
// @Description  Тут створюється нова послуга
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        input  body     models.NewService  true  "Данные для создания услуги"
// @Success      200  {object}  map[string]interface{}
// @Router       /api/service/create [post]
// @Security     ApiKeyAuth
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

// @Summary      Створення промо-послуги
// @Description  Хендлер для створення промо-послуги. Приймає структуру NewPromoService і повертає ідентифікатор нової послуги.
// @Tags         Промокод
// @Accept       json
// @Produce      json
// @Param        input  body     models.NewPromoService  true  "Дані для створення промо-послуги"
// @Success      200    {object}  map[string]interface{}  "ID створеної промо-послуги"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/create_promo [post]
func (h *Handler) CreatePromoService(c *gin.Context) {
	var input models.NewPromoService
	userId, err := getUserId(c)
	input.PromoService.PerformerID = userId
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	info, err := h.services.Offering.CreatePromoService(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"service_info": info,
	})
}

// @Summary      Валідація промокоду
// @Description  Хендлер для перевірки дійсності промокоду. Приймає структуру PromocodeValidationInput і повертає результат перевірки.
// @Tags         Промокод
// @Accept       json
// @Produce      json
// @Param        input  body     models.PromocodeValidationInput  true  "Дані для перевірки промокоду"
// @Success      200    {object}  models.PromocodeValidation  "Результат перевірки промокоду"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Router       /api/service/validate_promo [post]
func (h *Handler) ValidatePromoCode(c *gin.Context) {
	var input models.PromocodeValidationInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var valid models.PromocodeValidation
	valid, err := h.services.Offering.ValidatePromoCode(input.Promocode)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, valid)
}

// @Summary      Оновлення послуги
// @Description  Хендлер для оновлення даних послуги. Приймає структуру ServiceUpdate.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        input  body     models.ServiceUpdate  true  "Дані для оновлення послуги"
// @Success      200    {object}  map[string]interface{}  "Повідомлення про успішне оновлення"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Security     BearerAuth
// @Router       /api/service/update [post]
// @Security     ApiKeyAuth
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

// @Summary      Видалення послуги
// @Description  Хендлер для видалення послуги. Приймає структуру ServiceDelete. Перевіряє, чи користувач є адміністратором або власником послуги.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        input  body     models.ServiceDelete  true  "Дані для видалення послуги"
// @Success      200    {object}  map[string]interface{}  "Повідомлення про успішне видалення або помилку"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/delete [post]
// @Security     ApiKeyAuth
func (h *Handler) DeleteService(c *gin.Context) {
	var input models.ServiceDelete
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	isAdmin, err := h.services.User.IsAdmin(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	owner, err := h.services.GetServiceOwner(input.UUID)
	fmt.Println("from handler", owner)
	if isAdmin == true || owner == userId {
		if err := h.services.Offering.DeleteService(input.UUID); err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "ok",
			"status":  "deleted",
		})

	} else {

		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "User not owner or admin",
			"status":  "Not deleted",
		})

	}

}

// @Summary      Активація промокоду
// @Description  Хендлер для активації промокоду. Приймає промокод та активує його для поточного користувача.
// @Tags         Промокод
// @Accept       json
// @Produce      json
// @Param        input  body     models.PromocodeActivationInput  true  "Дані для активації промокоду"
// @Success      200    {object}  map[string]string  "Повідомлення про успішну активацію"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      401    {object}  map[string]string  "Користувач не авторизований"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/activate_promo [post]
func (h *Handler) ActivatePromoCode(c *gin.Context) {
	var input models.PromocodeActivationInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user_id, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Offering.ActivatePromoCode(input.Promocode, user_id); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "promocode activated",
		"status":  "ok",
	})
}

// @Summary      Отримання актуальних послуг користувача
// @Description  Хендлер для отримання списку актуальних послуг, створених поточним користувачем.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Success      200    {array}   models.Service  "Список актуальних послуг"
// @Failure      401    {object}  map[string]string  "Користувач не авторизований"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/getmyactualservices [get]
func (h *Handler) GetMyActualServices(c *gin.Context) {

	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	services, err := h.services.Offering.GetMyActualServices(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, services)
}

// @Summary      Отримання історії послуг користувача
// @Description  Хендлер для отримання історії послуг поточного користувача з пагінацією.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        input  body     models.MyHistoryServiceInput  true  "Дані для отримання історії послуг"
// @Success      200    {array}   models.Service  "Список історичних послуг"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      401    {object}  map[string]string  "Користувач не авторизований"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/getmyhistory [post]
func (h *Handler) GetMyHistory(c *gin.Context) {
	var input models.MyHistoryServiceInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	services, err := h.services.Offering.GetHistoryMyServices(userId, input.Limit, input.Offset)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, services)
}

// @Summary      Отримання типів послуг
// @Description  Хендлер для отримання всіх типів послуг.
// @Tags         Послуга
// @Produce      json
// @Success      200    {array}   models.ServiceType  "Список типів послуг"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/gettypes [get]
// @Security     ApiKeyAuth
func (h *Handler) GetTypes(c *gin.Context) {

	types, err := h.services.Offering.GetTypes()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, types)

}

// @Summary      Створення типу послуги
// @Description  Хендлер для створення нового типу послуги. Приймає структуру ServiceType.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        input  body     models.ServiceType  true  "Дані для створення типу послуги"
// @Success      200    {object}  map[string]interface{}  "Повідомлення про успішне створення типу послуги"
// @Failure      400    {object}  map[string]string  "Помилка запиту"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/createservicetype [post]
// @Security     ApiKeyAuth
func (h *Handler) CreateServiceType(c *gin.Context) {
	var input models.ServiceType
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Offering.CreateServiceType(input); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"status":  "created",
	})
}

// @Summary      Отримання послуг користувача
// @Description  Хендлер для отримання всіх послуг, що належать поточному користувачу. Використовує авторизацію для визначення користувача.
// @Tags         Послуга
// @Produce      json
// @Success      200    {array}   models.Service  "Список послуг користувача"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/getmyservices [get]
// @Security     ApiKeyAuth
func (h *Handler) GetMyServices(c *gin.Context) {

	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	services, err := h.services.Offering.GetMyServices(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, services)

}

// @Summary      Отримання доступних послуг
// @Description  Хендлер для отримання доступних послуг для користувача. Використовує авторизацію для визначення користувача.
// @Tags         Послуга
// @Produce      json
// @Success      200    {array}   models.Service  "Список доступних послуг"
// @Failure      500    {object}  map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/available [get]
// @Security     ApiKeyAuth
func (h *Handler) GetAvailableServices(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	services, err := h.services.Offering.GetAvailableService(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, services)
}

// @Summary      Отримання доступного часу для послуги
// @Description  Хендлер для отримання доступного часу для конкретної послуги. Приймає структуру з ID послуги в тілі запиту.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        serviceID  body     models.ServiceAvailableTimeInput  true  "Дані для отримання доступного часу послуги"
// @Success      200        {array}  map[string]interface{}  "Список доступного часу для послуги"
// @Failure      400        {object} map[string]string  "Помилка запиту"
// @Failure      500        {object} map[string]string  "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/availabletime [post]
// @Security     ApiKeyAuth
func (h *Handler) GetAvailableTime(c *gin.Context) {
	var serviceID models.ServiceAvailableTimeInput
	if err := c.BindJSON(&serviceID); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	available_times, err := h.services.Offering.GetAvailableTime(serviceID.ID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, available_times)
}
