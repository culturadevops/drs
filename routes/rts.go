package routes

import (
	"fmt"
	"os"

	"github.com/culturadevops/drs/handlers"

	"github.com/gofiber/fiber/v2"
)

func Ejemplo(e *fiber.App, Handler *handlers.Ejemplo) {
	key := "PATH_BASE"
	value, defined := os.LookupEnv(key)
	if !defined {
		fmt.Println("falta variable de entorno" + key)
		os.Exit(1)
	}
	r := e.Group("/" + value)
	r.Get("/list/:tabla", Handler.List())

	/*r.Get("/show/:id", Handler.Show())
	r.Post("/list", Handler.Add())
	r.Put("/update/:id", Handler.Update())
	r.Delete("/del/:id", Handler.Del())
	*/
}
