package main

import (
	osrmapi "GruzowikiRoutesGenerator/osrm_api"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	//pgr := db.PGRoutingQueries{}
	//pgr.EstablishConnection(os.Getenv("CONN_STRING"))
	//defer pgr.FuckingDestroyConnection()

	osrm := osrmapi.OSRMQueries{}
	osrm.ConfigureOSRMQueries(os.Getenv("CONN_STRING"),
								os.Getenv("OVERVIEW"),
								os.Getenv("ALTERNATIVES"),
								os.Getenv("STEPS"))

	app := fiber.New();

	// setup middleware
	//app.Use(func(c *fiber.Ctx) error {
    //    c.Locals("pgr", pgr)
    //    return c.Next()
    //})

	app.Use(func(c *fiber.Ctx) error {
        c.Locals("osrm", osrm)
        return c.Next()
    })

	app.Use(
		logger.New(), // add Logger middleware
	)

	app.Post("/api/create_route", OnCreateRouteHandlerWithOSRMApi)

	app.Listen(":8080")
}