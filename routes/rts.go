package routes

import (
	"github.com/culturadevops/drs/handlers"

	"github.com/gofiber/fiber/v2"
)

func Ejemplo(e *fiber.App, Handler *handlers.Ejemplo) {
	r := e.Group("/support-lima-dynamo")
	r.Get("/list/:tabla", Handler.List())

	/*r.Get("/show/:id", Handler.Show())
	r.Post("/list", Handler.Add())
	r.Put("/update/:id", Handler.Update())
	r.Delete("/del/:id", Handler.Del())
	*/
}
