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
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("error initializing configuration: %v", err)
	}
	defer cfg.Close()

	e.Use(middleware.Logger())
	setupRoutes(e, cfg)
	startServer(e, cfg)
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func setupRoutes(e *echo.Echo, cfg *config.Config) {
	setupBasicRoutes(e)
	setupTaxRoutes(e, cfg)
	setupAdminRoutes(e, cfg)
}

func setupBasicRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
}

func setupTaxRoutes(e *echo.Echo, cfg *config.Config) {
	taxController := controller.New(cfg)
	tax := e.Group("/tax")
	tax.POST("/calculations", taxController.TaxCalculateController)
	tax.POST("/calculations/upload-csv", taxController.TaxCalculateFormCsvController)
}

func setupAdminRoutes(e *echo.Echo, cfg *config.Config) {
	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == cfg.Admin && password == cfg.AdminPassword, nil
	}))

	taxController := controller.New(cfg)
	admin.POST("/deductions/:deductType", taxController.UpdatePersonalDeductionController)
}

func startServer(e *echo.Echo, cfg *config.Config) {
	apiPort := cfg.Port
	go func() {
		if err := e.Start(":" + apiPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	fmt.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
}
