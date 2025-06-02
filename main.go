package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Calculation struct {
	Id         string `gorm:"primaryKey" json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

// var calculations = []Calculation{}

func initDB() {
	dsn := "host=localhost user=postgres password=1233 dbname=postgres port=5432 sslmode=disable"

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not con to db: %v", err)
	}
	if err := db.AutoMigrate(&Calculation{}); err != nil {
		log.Fatalf("Migrate: %v", err)
	}
}

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
	var calcs []Calculation
	if err := db.Find(&calcs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calcs"})
	}
	return c.JSON(http.StatusOK, calcs)
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
	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calc"})
	}
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

	var calc Calculation
	if err := db.First(&calc, "id=?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not fund expr"})
	}
	calc.Expression = req.Expression
	calc.Expression = res

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calc"})
	}

	return c.JSON(http.StatusOK, calc)
}

func deleteCalc(c echo.Context) error {
	id := c.Param("id")

	if err := db.Delete(&Calculation{}, "id=?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calcilation not found"})
	}
	return c.NoContent(http.StatusNoContent)
}

func main() {
	fmt.Println("------------------")
	initDB()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalcs)
	e.POST("/calculations", postCalc)
	e.PATCH("/calculations/:id", patchCalc)
	e.DELETE("/calculations/:id", deleteCalc)

	e.Start("localhost:8080")

}
