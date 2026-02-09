package main

import (
	"log"

	"logforge/db"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	ServiceURL  string `json:"service_url"`
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Post("/register/service", func(c *fiber.Ctx) error {
		var service Service

		if err := c.BodyParser(&service); err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid request body"})
		}

		err := db.DB.QueryRow(
			c.Context(),
			`INSERT INTO service (service_name, service_url)
			VALUES ($1, $2)
			RETURNING service_id`,
			service.ServiceName,
			service.ServiceURL,
		).Scan(&service.ServiceID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "failed to register service"})
		}

		return c.Status(fiber.StatusCreated).JSON(service)
	})

	log.Fatal(app.Listen(":3000"))
}
