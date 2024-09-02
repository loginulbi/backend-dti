package main

import (

	"login-service/config"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
}
