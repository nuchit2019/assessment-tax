package main

import (
	"fmt"
	"log"
	"net/http"

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


	 // Tax Handler
	 taxController := controller.New(cfg)
	 e.POST("/tax/calculations", taxController.TaxCalculate)

	apiPort := cfg.Port
	if apiPort == "" {
		apiPort = "8080"
		log.Fatal("PORT environment variable not set")
	}

	address := fmt.Sprintf(":%s", apiPort)
	fmt.Printf("Server is running on port %s\n", apiPort)
	if err := e.Start(address); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
