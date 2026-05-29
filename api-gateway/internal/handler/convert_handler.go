package handler

import (
	"net/http"
	"strconv"

	"api-gateway/helper"
	pbconvert "api-gateway/proto/convert"

	"github.com/labstack/echo/v5"
)

type ConvertHandler struct {
	client pbconvert.ConvertClient
}

func NewConvertHandler(client pbconvert.ConvertClient) *ConvertHandler {
	return &ConvertHandler{client: client}
}

// ConvertToIDR godoc
// @Summary      Convert kg CO₂ to IDR value
// @Tags         emissions
// @Produce      json
// @Security     BearerAuth
// @Param        kg   query     number  true  "Emission in kg CO₂"  example(100.0)
// @Success      200  {object}  SuccessResponse{data=ConvertData}
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /emissions/convert [get]
func (h *ConvertHandler) ConvertToIDR(c *echo.Context) error {
	kgStr := c.QueryParam("kg")
	if kgStr == "" {
		return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", "kg query parameter is required"))
	}

	kg, err := strconv.ParseFloat(kgStr, 64)
	if err != nil || kg < 0 {
		return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", "kg must be a non-negative number"))
	}

	resp, err := h.client.ConvertToIDR(c.Request().Context(), &pbconvert.ConvertToIDRRequest{
		EmissionKgCo2: kg,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to convert emission"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]interface{}{
		"emission_kg_co2":        kg,
		"price_per_ton_usd":      resp.PricePerTonUsd,
		"exchange_rate_usd_idr":  resp.ExchangeRateUsdIdr,
		"total_idr":              resp.TotalIdr,
	}))
}
