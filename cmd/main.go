package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title IDM API documentation
// @Version 0.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("error closing db: %v", zap.Error(err))
		}
	}()

	// передаём конфиг, логгер и соединение к базе данных в функцию создания сервера
	var server = build(cfg, logger, database)

	// Запускаем сервер в отдельной горутине
	go func() {
		// загружаем сертификаты
		cer, err := tls.LoadX509KeyPair(cfg.SslSert, cfg.SslKey)
		if err != nil {
			logger.Panic("failed certificate loading: %s", zap.Error(err))
		}
		// создаём конфигурацию TLS сервера
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
		// создаём слушателя https соединения
		ln, err := tls.Listen("tcp", ":8080", tlsConfig)
		if err != nil {
			logger.Panic("failed TLS listener creating: %s", zap.Error(err))
		}
		// запускаем веб-сервер
		err = server.App.Listener(ln)
		if err != nil {
			logger.Panic("http server error: %s", zap.Error(err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg, logger)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	logger.Info("graceful shutdown complete")
}

// gracefulShutdown - функция "элегантного" завершения работы сервера по сигналу от операционной системы
func gracefulShutdown(
	server *web.Server,
	wg *sync.WaitGroup,
	logger *common.Logger,
) {
	// Уведомить основную горутину о завершении работы
	defer wg.Done()
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
}

// build - функция сборки приложения
func build(
	cfg common.Config,
	logger *common.Logger,
	database *sqlx.DB,
) *web.Server {
	// Создаём сервер
	var server = web.NewServer()
	// Регистрируем middleware
	server.App.Use("/swagger/*", swagger.HandlerDefault)
	server.App.Use(requestid.New())
	server.App.Use(recover.New())
	server.GroupApi.Use(web.AuthMiddleware(logger))

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
