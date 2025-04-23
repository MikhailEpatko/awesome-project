package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	App *fiber.App
	// группа публичного API
	GroupApiV1 fiber.Router
	// группа непубличного API
	GroupInternal fiber.Router
}

func NewServer() *Server {
	app := fiber.New()
	registerMiddleware(app)
	groupInternal := app.Group("/internal")
	groupApi := app.Group("/api")
	groupApiV1 := groupApi.Group("/v1")
	return &Server{
		App:           app,
		GroupApiV1:    groupApiV1,
		GroupInternal: groupInternal,
	}
}

func registerMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(requestid.New())
}
