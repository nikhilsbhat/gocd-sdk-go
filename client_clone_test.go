package gocd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_client_clone(t *testing.T) {
	t.Run("should return copied client", func(t *testing.T) {
		conf := &client{}

		actual := conf.clone()
		assert.NotNil(t, actual)
		assert.NotSame(t, conf, actual)
	})

	t.Run("should panic on copy error", func(t *testing.T) {
		expectedErr := errors.New("copy failed")
		originalCopyClientWithOption := copyClientWithOption
		copyClientWithOption = func(_, _ any) error {
			return expectedErr
		}
		t.Cleanup(func() {
			copyClientWithOption = originalCopyClientWithOption
		})

		assert.PanicsWithError(t, expectedErr.Error(), func() {
			(&client{}).clone()
		})
	})
}
