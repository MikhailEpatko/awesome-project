package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDivide(t *testing.T) {
	a := assert.New(t)

	t.Run("divide common", func(t *testing.T) {
		want := 5
		got, err := Divide(10, 2)
		a.Nil(err)
		a.Equal(want, got)
	})

	t.Run("divide zero", func(t *testing.T) {
		want := 0
		got, err := Divide(0, 2)
		a.Nil(err)
		a.Equal(want, got)
	})

	t.Run("divide by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		a.NotNil(err)
	})
}
