package main

import (
	"GruzowikiRoutesGenerator/db"
	"os"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	pgr := db.PGRoutingQueries{}
	pgr.EstablishConnection(os.Getenv("CONN_STRING"))
	defer pgr.FuckingDestroyConnection()

	app := fiber.New();

	// setup middleware
	app.Use(func(c *fiber.Ctx) error {
        c.Locals("pgr", pgr)
        return c.Next()
    })

	app.Use(
		logger.New(), // add Logger middleware
	)

	app.Post("/api/create_route", OnCreateRouteHandler)

	app.Listen(":8080")
}