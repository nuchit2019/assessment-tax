package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nuchit2019/assessment-tax/config"
	"github.com/nuchit2019/assessment-tax/controller"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("error initializing configuration: %v", err)
	}
	defer cfg.Close()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	taxController := controller.New(cfg)

	tax := e.Group("/tax")
	tax.POST("/calculations", taxController.TaxCalculateController)
	tax.POST("/calculations/upload-csv", taxController.TaxCalculateFormCsvController) // TODO

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == cfg.Admin && password == cfg.AdminPassword, nil
	}))

	admin.POST("/deductions/:deductType", taxController.UpdatePersonalDeductionController)

	apiPort := cfg.Port 

	go func() {
		if err := e.Start(":" + apiPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(e.Start(":" + apiPort))
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown

	// Print "shutting down the server"
	fmt.Println("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
