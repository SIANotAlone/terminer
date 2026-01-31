package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetGaveDoneRecordsStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	statistic, err := h.services.Statistic.GetProvidedDoneRecordsStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetGavedRecordsProMonthStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetProvidedRecordsProMonthStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetRecievedRecordsProMonthStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetRecievedRecordsProMonthStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetMassageProTypeStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetMassageProTypeStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetResievedMassageProTypeStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetResievedMassageProTypeStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetRecievedServicesByTypesStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetResievedServicesTypes(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetProvidetServicesByTypesStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetProvidedServicesTypes(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}

func (h *Handler) GetAvailableYears(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	years, err := h.services.Statistic.GetAvailableYears(user_id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, years)
}

func (h *Handler) GetMainStatistic(c *gin.Context) {
	var user_id, err = getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	statistic, err := h.services.Statistic.GetMainStatistic(user_id, yearInt)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statistic)
}
