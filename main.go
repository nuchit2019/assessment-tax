package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nuchit2019/assessment-tax/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	cfg,err:=config.New()
	if err!=nil{
		log.Fatal("error initializing configuration: %v", err)
	}

	defer cfg.Close()

	// e.Logger.Fatal(e.Start(":1323"))
	apiPort := cfg.Port
	fmt.Printf("Server is running on port %s\n", apiPort)
	e.Logger.Fatal(e.Start(":" + apiPort))

}
