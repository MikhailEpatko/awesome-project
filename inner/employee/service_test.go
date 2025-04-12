package employee

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"idm/inner/validator"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindById(id int64) (employee Entity, err error) {
	args := m.Called(id)
	return args.Get(0).(Entity), args.Error(1)
}

func (m *MockRepo) Save(entity Entity) (int64, error) {
	args := m.Called(entity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) BeginTransaction() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockRepo) FindByNameTx(tx *sqlx.Tx, name string) (bool, error) {
	args := m.Called(tx, name)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockRepo) SaveTx(tx *sqlx.Tx, entity Entity) (int64, error) {
	args := m.Called(tx, entity)
	return args.Get(0).(int64), args.Error(1)
}

func TestFindById(t *testing.T) {
	var a = assert.New(t)

	t.Run("should return found employee", func(t *testing.T) {
		var repo = new(MockRepo)
		var vld = validator.New()
		var svc = NewService(repo, vld)
		var entity = Entity{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var want = entity.toResponse()

		repo.On("FindById", int64(1)).Return(entity, nil)

		var got, err = svc.FindById(1)

		a.Nil(err)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("should return wrapped error", func(t *testing.T) {
		var repo = new(MockRepo)
		var vld = validator.New()
		var svc = NewService(repo, vld)
		var entity = Entity{}
		var err = errors.New("database error")
		var want = fmt.Errorf("error finding employee with id 1: %w", err)

		repo.On("FindById", int64(1)).Return(entity, err)

		var response, got = svc.FindById(1)

		a.Empty(response)
		a.NotNil(got)
		a.Equal(want, got)
		a.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})
}
