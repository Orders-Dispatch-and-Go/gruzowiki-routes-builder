package main

import (
	"GruzowikiRoutesGenerator/db"
	osrmapi "GruzowikiRoutesGenerator/osrm_api"
	"GruzowikiRoutesGenerator/utils/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var OnCreateRouteHandlerWithPGRouting =
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

		routeDto, err := pgr.BuildRoute(createRouteDTO)
		if err != nil && !routeDto.IsDtoValid() {
			log.Error("QueryRow failed: %v\n", err)
			return c.Status(fiber.StatusTeapot).JSON(routeDto)
		} else if err != nil {
			log.Error("QueryRow failed: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Status(fiber.StatusOK).JSON(routeDto)
	}

var OnCreateRouteHandlerWithOSRMApi =
	func(c *fiber.Ctx) error {
		c.Accepts("html")             // "html"
		c.Accepts("text/html")        // "text/html"
		c.Accepts("json", "text")     // "json"
		c.Accepts("application/json") // "application/json"

		createRouteDTO := new(dto.CreateRouteDTO)

		if err := c.BodyParser(createRouteDTO); err != nil {
            return err
        }

		osrm := c.Locals("osrm").(osrmapi.OSRMQueries)

		routeDto, err := osrm.BuildRoute(createRouteDTO)
		if err != nil && !routeDto.IsDtoValid() {
			log.Error("QueryRow failed: %v\n", err)
			return c.Status(fiber.StatusTeapot).JSON(routeDto)
		} else if err != nil {
			log.Error("QueryRow failed: %v\n", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Status(fiber.StatusOK).JSON(routeDto)
	}