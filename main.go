package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/komron-dev/hrms-fiber/database"
	"github.com/komron-dev/hrms-fiber/handlers"
	"log"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Get("/employee", handlers.ListEmployees)
	app.Get("/employee/:id", handlers.GetEmployeeById)
	app.Post("/employee", handlers.CreateEmployee)
	app.Put("/employee/:id", handlers.UpdateEmployee)
	app.Delete("/employee/:id", handlers.DeleteEmployee)

	log.Fatal(app.Listen(":3000"))
}
