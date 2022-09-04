package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetHealthInfo(t *testing.T) {
	t.Run("should be able to fetch the server health status", func(t *testing.T) {
		client := NewClient(
			"http://localhost:8153/go",
			"",
			"",
			"info",
			nil,
			"info",
		)

		healthStatus, err := client.GetHealthInfo()
		assert.NoError(t, err)
		log.Println(healthStatus)
		assert.Equal(t, 5, len(healthStatus))
	})
}
