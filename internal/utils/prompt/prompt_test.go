package prompt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrompt_Ask(t *testing.T) {
	t.Run("Should run command without panics", func(t *testing.T) {
		assert.NotPanics(t, func() {
			_, _ = NewPrompt().Ask("What's your name?", "GoGenerator")
		})
	})
}