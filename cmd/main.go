package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"syscall"
	"time"
)

// @title IDM API documentation
// @BasePath /api/v1/
func main() {
	// Перенесли сюда из функции build() парсинг конфига
	var cfg = common.GetConfig(".env")

	// Создаем логгер
	var logger = common.NewLogger(cfg)
	// Отложенный вызов записи сообщений из буфера в лог. Необходимо вызывать перед выходом из приложения
	defer func() {
		err := logger.Sync()
		fmt.Printf("logger synchronization error: %v", err)
	}()
	var database = common.ConnectDbWithCfg(cfg)

	// передаём конфиг, логгер и соединение к базе данных в функцию создания сервера
	var server = build(cfg, logger, database)

	// Запускаем сервер в отдельной горутине
	go func() {
		var err = server.App.Listen(":8080")
		if err != nil {
			logger.Panic("http server error: %s", zap.Error(err))
		}
	}()

	// Создаем канал для ожидания сигнала завершения работы сервера
	var done = make(chan bool, 1)
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, done, logger)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	<-done
	logger.Info("graceful shutdown complete")
}

// gracefulShutdown - функция "элегантного" завершения работы сервера по сигналу от операционной системы
func gracefulShutdown(
	server *web.Server,
	done chan bool,
	logger *common.Logger,
) {
	// Создаём контекст, который слушает сигнал прерывания от ОС.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Слушаем сигнал прерывания от ОС
	<-ctx.Done()
	// заменили отладочную печать на логирование
	logger.Info("shutting down gracefully")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		// Запись ошибки в лог
		logger.Error("Server forced to shutdown with error", zap.Error(err))
	}
	logger.Info("Server exiting")
	// Уведомить основную горутину о завершении работы
	done <- true
}

// build - функция сборки приложения
func build(
	cfg common.Config,
	logger *common.Logger,
	database *sqlx.DB,
) *web.Server {
	var server = web.NewServer()
	var vld = validator.New()

	var employeeRepo = employee.NewRepository(database)
	var employeeService = employee.NewService(employeeRepo, vld)
	// передаём логгер в конструктор контроллера
	var employeeController = employee.NewController(server, employeeService, logger)
	employeeController.RegisterRoutes()

	var infoController = info.NewController(server, cfg)
	infoController.RegisterRoutes()

	return server
}
