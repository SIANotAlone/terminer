package handler

import (
	"net/http"
	"terminer/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Description  Хендлер для видалення послуги. Приймає структуру ServiceDelete. Перевіряє, чи є користувачвласником послуги.
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

	if err := h.services.Offering.DeleteService(input.UUID, userId); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})

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
// @Success      200    {array}   models.UserServiceHistory  "Історія моїх послуг"
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

// @Summary      Отримання типів масажу
// @Description  Хендлер для отримання списку всіх доступних типів масажу.
// @Tags         Послуга
// @Produce      json
// @Success      200    {array}   models.MassageType  "Список типів масажу"
// @Failure      500    {object}  map[string]string   "Помилка сервера"
// @Security     BearerAuth
// @Router       /api/service/getmassagetypes [get]
// @Security     ApiKeyAuth
func (h *Handler) GetMassageTypes(c *gin.Context) {
	types, err := h.services.Offering.GetMassageTypes()
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

// @Summary      Отримання повної інформації про послугу
// @Description  Хендлер для отримання детальної інформації про конкретну послугу за її ID. ID послуги передається як параметр маршруту.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        serviceid  path     string  true  "ID послуги (UUID)"
// @Success      200        {object} models.FullServiceInformation  "Повна інформація про послугу"
// @Failure      400        {object} map[string]string  "Некоректний формат ID послуги (невалідна UUID)"
// @Failure      500        {object} map[string]string  "Помилка сервера (наприклад, помилка отримання даних)"
// @Security     BearerAuth
// @Router       /api/service/{serviceid} [get]
func (h *Handler) GetFullServiceInfo(c *gin.Context) {
	serviceID := c.Param("serviceid")
	parsedUUID, err := uuid.Parse(serviceID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	service, err := h.services.Offering.GetFullServiceInfo(parsedUUID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, service)
}

// @Summary      Редагування існуючої послуги
// @Description  Хендлер для оновлення детальної інформації про існуючу послугу. Приймає повну структуру послуги (models.ServiceInformation) в тілі запиту. Користувач повинен бути автентифікований, а його ID використовується для перевірки прав доступу.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        service  body     models.ServiceInformation  true  "Дані послуги для оновлення. Повинен містити ID послуги."
// @Success      200      {object} map[string]string  "Послугу успішно оновлено"
// @Failure      400      {object} map[string]string  "Помилка вхідних даних (невірний формат JSON)"
// @Failure      500      {object} map[string]string  "Помилка сервера (неможливо отримати ID користувача, або помилка оновлення послуги в базі даних)"
// @Security     BearerAuth
// @Router       /api/service/edit [put]
func (h *Handler) EditService(c *gin.Context) {
	var service models.ServiceInformation
	if err := c.BindJSON(&service); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var userID uuid.UUID
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Offering.EditService(service, userID); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

// @Summary      Додавання нового доступного часу/вікна для послуги
// @Description  Хендлер для створення нового запису про доступний час (вікно) для певної послуги. Приймає структуру з даними про час і послугу. ID користувача використовується для перевірки прав доступу/власності.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        availableFor  body     models.NewAvailableFor  true  "Дані про новий доступний час (наприклад, дата, час початку/кінця, ID послуги)"
// @Success      200           {object} map[string]interface{} "Успішне створення запису"
// @Failure      400           {object} map[string]string  "Помилка вхідних даних (невірний формат JSON)"
// @Failure      500           {object} map[string]string  "Помилка сервера (неможливо отримати ID користувача, або помилка при створенні запису в базі даних)"
// @Security     BearerAuth
// @Router       /api/service/availablefor [post]
func (h *Handler) NewAvailableFor(c *gin.Context) {
	var availableFor models.NewAvailableFor
	if err := c.BindJSON(&availableFor); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var userID uuid.UUID
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	new_id, err := h.services.Offering.NewAvailableFor(availableFor, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"id":      new_id,
	})
}

// @Summary      Видалення запису про доступний час/вікно
// @Description  Хендлер для видалення існуючого запису про доступний час для певної послуги. Приймає структуру з ID запису, який потрібно видалити, у тілі запиту. ID користувача використовується для перевірки прав доступу/власності.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        availableFor  body     models.DeleteAvailableFor  true  "Дані для видалення доступного часу (повинно містити ID запису)"
// @Success      200           {object} map[string]string "Запис про доступний час успішно видалено"
// @Failure      400           {object} map[string]string  "Помилка вхідних даних (невірний формат JSON)"
// @Failure      500           {object} map[string]string  "Помилка сервера (неможливо отримати ID користувача, помилка видалення запису в базі даних або відсутність прав)"
// @Security     BearerAuth
// @Router       /api/service/availablefor [delete]
func (h *Handler) DeleteAvailableFor(c *gin.Context) {
	var availableFor models.DeleteAvailableFor
	if err := c.BindJSON(&availableFor); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var userID uuid.UUID
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Offering.DeleteAvailableFor(availableFor, userID); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

// @Summary      Додавання нового доступного часу для послуги (одиничний слот)
// @Description  Хендлер для створення нового запису про доступний час (одиничний слот) для певної послуги. Приймає структуру з даними слоту. ID користувача використовується для перевірки прав доступу/власності на послугу.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        availableTime  body     models.NewAvailableTime  true  "Дані про новий доступний часовий слот (наприклад, ID послуги, час початку/кінця)"
// @Success      200            {object} map[string]interface{} "Успішне створення часового слоту"
// @Failure      400            {object} map[string]string  "Помилка вхідних даних (невірний формат JSON)"
// @Failure      500            {object} map[string]string  "Помилка сервера (неможливо отримати ID користувача, або помилка при створенні запису в базі даних)"
// @Security     BearerAuth
// @Router       /api/service/availabletime [post]
func (h *Handler) NewAvailableTime(c *gin.Context) {
	var availableTime models.NewAvailableTime
	if err := c.BindJSON(&availableTime); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var userID uuid.UUID
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	new_id, err := h.services.Offering.NewAvailableTime(availableTime, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"id":      new_id,
	})
}

// @Summary      Видалення запису про доступний часовий слот
// @Description  Хендлер для видалення існуючого одиничного часового слоту (available time). Приймає структуру з ID слоту, який потрібно видалити, у тілі запиту. ID користувача використовується для перевірки прав доступу/власності на послугу.
// @Tags         Послуга
// @Accept       json
// @Produce      json
// @Param        availableTime  body     models.DeleteAvailableTime  true  "Дані для видалення доступного часового слоту (повинно містити ID слоту)"
// @Success      200            {object} map[string]string "Запис про доступний час успішно видалено"
// @Failure      400            {object} map[string]string  "Помилка вхідних даних (невірний формат JSON)"
// @Failure      500            {object} map[string]string  "Помилка сервера (неможливо отримати ID користувача, помилка видалення запису в базі даних або відсутність прав)"
// @Security     BearerAuth
// @Router       /api/service/availabletime [delete]
func (h *Handler) DeleteAvailableTime(c *gin.Context) {
	var availableTime models.DeleteAvailableTime
	if err := c.BindJSON(&availableTime); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var userID uuid.UUID
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if err := h.services.Offering.DeleteAvailableTime(availableTime, userID); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
