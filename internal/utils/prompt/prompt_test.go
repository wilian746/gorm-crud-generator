package prompt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPrompt(t *testing.T) {
	t.Run("Should run command without panics", func(t *testing.T) {
		assert.NotNil(t, NewPrompt())
	})
}
