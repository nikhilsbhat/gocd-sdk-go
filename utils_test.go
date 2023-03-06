package gocd_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func Test_getSLice(t *testing.T) {
	t.Run("", func(t *testing.T) {
		input := []interface{}{"one", "two", "three", "four"}

		expected := []string{"one", "two", "three", "four"}

		actual := gocd.GetSLice(input)
		assert.Equal(t, expected, actual)
	})
}
