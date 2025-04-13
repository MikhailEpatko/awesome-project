package tests

import (
	"idm/inner/employee"
)

type Fixture struct {
	employees *employee.Repository
}

func NewFixture(employees *employee.Repository) *Fixture {
	return &Fixture{employees}
}

func (f *Fixture) Employee(name string) int64 {
	var entity = employee.Entity{
		Name: name,
	}
	var newId, err = f.employees.Save(&entity)
	if err != nil {
		panic(err)
	}
	return newId
}
