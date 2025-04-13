package employee

import (
	"encoding/json"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/common"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	mock.Mock
}

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

	t.Run("should return created employee id", func(t *testing.T) {
		var app = fiber.New()
		var svc = new(MockService)
		var controller = NewController(app, svc)
		controller.RegisterRoutes()
		var body = strings.NewReader("{\"name\": \"john doe\"}")
		var req = httptest.NewRequest(fiber.MethodPost, "/api/v1/employees", body)
		req.Header.Set("Content-Type", "application/json")

		svc.On("CreateEmployee", mock.AnythingOfType("CreateRequest")).Return(int64(123), nil)

		resp, err := app.Test(req)

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
