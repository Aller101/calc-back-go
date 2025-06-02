package handlers

import (
	"net/http"

	"github.com/Aller101/calc-back-go/internal/service"
	"github.com/labstack/echo/v4"
)

type CalcHandler struct {
	service service.CalcilationService
}

func NewCalcHandler(s service.CalcilationService) *CalcHandler {
	return &CalcHandler{service: s}
}

func (h *CalcHandler) GetCalcs(c echo.Context) error {
	calcs, err := h.service.GetAllCalculations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calcs"})
	}
	return c.JSON(http.StatusOK, calcs)
}

func (h *CalcHandler) PostCalc(c echo.Context) error {
	var req service.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create calc"})
	}
	return c.JSON(http.StatusCreated, calc)
}

func (h *CalcHandler) PatchCalc(c echo.Context) error {
	id := c.Param("id")
	var req service.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	calc, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calc"})
	}
	return c.JSON(http.StatusOK, calc)
}

func (h *CalcHandler) DeleteCalc(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteCalculation(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calc"})
	}
	return c.NoContent(http.StatusNoContent)
}
