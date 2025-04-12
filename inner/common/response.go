package common

import "github.com/gofiber/fiber"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"error"`
	Data    any    `json:"data"`
}

func ErrResponse(
	c *fiber.Ctx,
	code int,
	message string,
) error {
	return c.Status(code).JSON(&Response{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func OkResponse(
	c *fiber.Ctx,
	data any,
) error {
	return c.JSON(&Response{
		Success: true,
		Data:    data,
	})
}
