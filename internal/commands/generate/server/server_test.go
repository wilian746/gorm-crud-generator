package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/wilian746/go-generator/internal/enums/database"
	"testing"
)

func TestServer_CreateFoldersAndFiles(t *testing.T) {
	t.Run("Create default folders without error", func(t *testing.T) {
		path := "/home/wilian/go/src/github.com/wilian746/go-generator/tpm"
		module := "github.com/wilian746/go-generator"
		err := NewServer().CreateFoldersAndFiles(path, module, database.Gorm)
		assert.NoError(t, err)
	})
}
