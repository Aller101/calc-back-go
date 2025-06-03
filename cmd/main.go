package main

import (
	"log"

	"github.com/Aller101/calc-back-go/internal/config"
	"github.com/Aller101/calc-back-go/internal/db"
	"github.com/Aller101/calc-back-go/internal/handlers"
	"github.com/Aller101/calc-back-go/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	cfg := config.MustLoad()

	datab, err := db.InitDB(*cfg)
	if err != nil {
		log.Fatalf("Could not con to db: %v", err)
	}

	calcRepo := service.NewCalcRepository(datab)
	calcServ := service.NewCalculationService(calcRepo)
	calcHandl := handlers.NewCalcHandler(calcServ)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", calcHandl.GetCalcs)
	e.POST("/calculations", calcHandl.PostCalc)
	e.PATCH("/calculations/:id", calcHandl.PatchCalc)
	e.DELETE("/calculations/:id", calcHandl.DeleteCalc)

	e.Start(cfg.Address)

}
