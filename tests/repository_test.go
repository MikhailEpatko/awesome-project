package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/database"
	"idm/inner/employee"
	"testing"
)

func TestRepository(t *testing.T) {
	a := assert.New(t)
	var db = database.Connect()
	var clearDatabase = func() {
		db.MustExec("delete from employee")
	}
	defer func() {
		if r := recover(); r != nil {
			clearDatabase()
		}
	}()
	var employeeRepository = employee.NewRepository(db)
	var fixture = NewFixture(employeeRepository)

	t.Run("find an employee by id", func(t *testing.T) {
		var newEmployeeId = fixture.Employee("Test Name")

		got, err := employeeRepository.FindById(newEmployeeId)

		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal("Test Name", got.Name)
		clearDatabase()
	})
}
