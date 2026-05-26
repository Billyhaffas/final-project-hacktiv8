package http

import (
	"net/http"

	"convert-emission-service/internal/domain"

	"github.com/labstack/echo/v5"
)

type ConversionHandler struct {
	usecase domain.ConversionUsecase
}

func NewConversionHandler(u domain.ConversionUsecase) *ConversionHandler {
	return &ConversionHandler{usecase: u}
}

// HandleDaily godoc
// @Summary      Convert daily carbon emissions to currency valuation
// @Description  Accepts raw user daily emissions in kg CO2 and translates them into financial costs (USD and local currency) using country carbon tax pricing data.
// @Tags         conversion
// @Accept       json
// @Produce      json
// @Param        country_code  query     string                         true  "3-letter country code (e.g., IDN, ARG)"
// @Param        payload       body      domain.UserDailyEmission       true  "Daily emission data"
// @Success      200           {object}  domain.UserDailyCostResponse
// @Failure      400           {object}  map[string]string              "Example: {'error': 'invalid country_code query param format'}"
// @Failure      500           {object}  map[string]string              "Example: {'error': 'internal database failure'}"
// @Router       /api/v1/convert/daily [post]
func (h *ConversionHandler) HandleDaily(e *echo.Context) error {
	countryCode := e.QueryParam("country_code")
	if len(countryCode) != 3 {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid country_code query param format"})
	}

	var payload domain.UserDailyEmission
	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "failed decoding request payload"})
	}

	res, err := h.usecase.ConvertDailyEmission(e.Request().Context(), countryCode, payload)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return e.JSON(http.StatusOK, res)
}

// HandleMonthly godoc
// @Summary      Convert monthly carbon emissions to currency valuation
// @Description  Accepts monthly aggregate emissions and breakdown lists, returning mapped cost breakdowns per day alongside totals.
// @Tags         conversion
// @Accept       json
// @Produce      json
// @Param        country_code  query     string                         true  "3-letter country code (e.g., IDN, ARG)"
// @Param        payload       body      domain.UserMonthlyEmission     true  "Monthly emission dataset"
// @Success      200           {object}  domain.UserMonthlyCostResponse
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /api/v1/convert/monthly [post]
func (h *ConversionHandler) HandleMonthly(e *echo.Context) error {
	countryCode := e.QueryParam("country_code")
	if len(countryCode) != 3 {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid country_code query param format"})
	}

	var payload domain.UserMonthlyEmission
	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "failed decoding request payload"})
	}

	res, err := h.usecase.ConvertMonthlyEmission(e.Request().Context(), countryCode, payload)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return e.JSON(http.StatusOK, res)
}

// HandleYearly godoc
// @Summary      Convert yearly carbon emissions to currency valuation
// @Description  Accepts structural yearly datasets and returns complex valuation arrays grouped by calendar month.
// @Tags         conversion
// @Accept       json
// @Produce      json
// @Param        country_code  query     string                         true  "3-letter country code (e.g., IDN, ARG)"
// @Param        payload       body      domain.UserYearlyEmission      true  "Yearly emission dataset"
// @Success      200           {object}  domain.UserYearlyCostResponse
// @Failure      400           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /api/v1/convert/yearly [post]
func (h *ConversionHandler) HandleYearly(e *echo.Context) error {
	countryCode := e.QueryParam("country_code")
	if len(countryCode) != 3 {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid country_code query param format"})
	}

	var payload domain.UserYearlyEmission
	if err := e.Bind(&payload); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "failed decoding request payload"})
	}

	res, err := h.usecase.ConvertYearlyEmission(e.Request().Context(), countryCode, payload)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return e.JSON(http.StatusOK, res)
}
