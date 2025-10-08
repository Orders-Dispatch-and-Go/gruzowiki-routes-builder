package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"GruzowikiRoutesGenerator/utils/dto"
	"GruzowikiRoutesGenerator/db"
)

var OnCreateRouteHandler =
	func(c *fiber.Ctx) error {
		c.Accepts("html")             // "html"
		c.Accepts("text/html")        // "text/html"
		c.Accepts("json", "text")     // "json"
		c.Accepts("application/json") // "application/json"

		createRouteDTO := new(dto.CreateRouteDTO)

		if err := c.BodyParser(createRouteDTO); err != nil {
            return err
        }

		pgr := c.Locals("pgr").(db.PGRoutingQueries)

		routeDto, err := pgr.BuildRout(createRouteDTO)
		if err != nil && !routeDto.IsDtoValid() {
			log.Error("QueryRow failed: %v\n", err)
			return c.Status(fiber.StatusTeapot).JSON(routeDto)
		} else if err != nil {
			log.Error("QueryRow failed: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Status(fiber.StatusOK).JSON(routeDto)
	}