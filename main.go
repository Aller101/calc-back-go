package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Calculation struct {
	Id         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

var calculations = []Calculation{}

func calcExpr(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	res, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", res), nil
}

func getCalcs(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}

func postCalc(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	res, err := calcExpr(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}
	calc := Calculation{
		Id:         uuid.NewString(),
		Expression: req.Expression,
		Result:     res,
	}
	calculations = append(calculations, calc)
	return c.JSON(http.StatusCreated, calc)
}

func patchCalc(c echo.Context) error {
	id := c.Param("id")
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	res, err := calcExpr(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	for i, v := range calculations {
		if v.Id == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = res
			return c.JSON(http.StatusOK, calculations[i])
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcilation not found"})
}

func deleteCalc(c echo.Context) error {
	id := c.Param("id")
	for i, v := range calculations {
		if v.Id == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcilation not found"})

}

func main() {
	fmt.Println("------------------")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalcs)
	e.POST("/calculations", postCalc)
	e.PATCH("/calculations/:id", patchCalc)
	e.DELETE("/calculations/:id", deleteCalc)

	e.Start("localhost:8080")

}
