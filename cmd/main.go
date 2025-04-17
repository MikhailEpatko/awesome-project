package main

import (
	"context"
	"fmt"
	"idm/inner/common"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var server = build()
	// Создаем канал для ожидания сигнала завершения работы сервера
	var done = make(chan bool, 1)
	// Запускаем сервер в отдельной горутине
	go func() {
		var err = server.App.Listen(":8080")
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, done)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	<-done
	fmt.Println("Graceful shutdown complete.")
}

// Функция "элегантного" завершения работы сервера по сигналу от операционной системы
func gracefulShutdown(server *web.Server, done chan bool) {
	// Создаём контекст, который слушает сигнал прерывания от ОС.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Слушаем сигнал прерывания от ОС
	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server forced to shutdown with error: %v\n", err)
	}
	fmt.Println("Server exiting")
	// Уведомить основную горутину о завершении работы
	done <- true
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
