package main

import (
	"fmt"
	"idm/inner/common"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/validator"
	"idm/inner/web"
)

func main() {
	var server = build()
	var err = server.App.Listen(":8080")
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func build() *web.Server {
	var cfg = common.GetConfig(".env")
	var server = web.NewServer()
	var database = common.ConnectDbWithCfg(cfg)
	var vld = validator.New()

	var employeeRepo = employee.NewRepository(database)
	var employeeService = employee.NewService(employeeRepo, vld)
	var employeeController = employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()

	var infoController = info.NewController(server, cfg)
	infoController.RegisterRoutes()

	return server
}
