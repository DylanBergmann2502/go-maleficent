// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maleficent/go-maleficent/configs"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Go Maleficent API is running",
		})
	})

	address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	e.Logger.Fatal(e.Start(address))
}
