package employee

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Объявляем структуру мока сервиса employee.Service
type MockService struct {
	mock.Mock
}

// Реализуем функции мок-сервиса
func (svc *MockService) FindById(id int64) (Response, error) {
	args := svc.Called(id)
	return args.Get(0).(Response), args.Error(1)
}

func (svc *MockService) CreateEmployee(request CreateRequest) (int64, error) {
	args := svc.Called(request)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateEmployee(t *testing.T) {
	var a = assert.New(t)

	// тестируем положительный сценарий: работника создали и получили его id
	t.Run("should return created employee id", func(t *testing.T) {
		// Готовим тестовое окружение
		server := web.NewServer()
		var svc = new(MockService)
		var cfg = common.Config{
			LogDevelopMode: true,
			LogLevel:       "debug",
		}
		var logger = common.NewLogger(cfg)
		var controller = NewController(server, svc, logger)
		controller.RegisterRoutes()
		// Готовим тестовое окружение
		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var req = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		req.Header.Set("Content-Type", "application/json")

		// Настраиваем поведение мока в тесте
		svc.On("CreateEmployee", mock.AnythingOfType("CreateRequest")).Return(int64(123), nil)

		// Отправляем тестовый запрос на веб сервер
		resp, err := server.App.Test(req)

		// Выполняем проверки полученных данных
		a.Nil(err)
		a.NotEmpty(resp)
		a.Equal(http.StatusOK, resp.StatusCode)
		bytesData, err := io.ReadAll(resp.Body)
		a.Nil(err)
		var responseBody common.Response
		err = json.Unmarshal(bytesData, &responseBody)
		a.Nil(err)
		a.Equal(float64(123), responseBody.Data)
		a.True(responseBody.Success)
		a.Empty(responseBody.Message)
	})
}
