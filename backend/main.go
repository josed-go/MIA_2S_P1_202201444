package main

import (
	"backend/analizador"
	"backend/utilidades"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Entrada struct {
	Texto string `json:"entrada"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Post("/analizar", func(c *fiber.Ctx) error {

		var entrada Entrada

		c.BodyParser(&entrada)

		analizador.Analyze(entrada.Texto)

		return c.JSON(&fiber.Map{
			"response": utilidades.ObtenerRespuestas(),
		})

	})

	app.Listen(":3000")
	fmt.Println("Servidor en el puerto 3000")

}
