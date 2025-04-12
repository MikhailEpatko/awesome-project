package employee

import (
	"errors"
	"github.com/gofiber/fiber"
	"idm/inner/common"
)

type Controller struct {
	app             *fiber.App
	employeeService Svc
}

type Svc interface {
	FindById(id int64) (Response, error)
	CreateEmployee(request CreateRequest) (int64, error)
}

func NewController(app *fiber.App, employeeService Svc) *Controller {
	return &Controller{
		app:             app,
		employeeService: employeeService,
	}
}

func (c *Controller) RegisterRoutes() {
	api := c.app.Group("/api")
	api.Post("/v1/employees", c.CreateEmployee)
}

func (c *Controller) CreateEmployee(ctx *fiber.Ctx) {
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		return
	}
	var newEmployeeId, err = c.employeeService.CreateEmployee(request)
	if err != nil {
		switch {
		case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
			_ = common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
		default:
			_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, err.Error())
		}
		return
	}
	if err = common.OkResponse(ctx, newEmployeeId); err != nil {
		_ = common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
		return
	}
}
