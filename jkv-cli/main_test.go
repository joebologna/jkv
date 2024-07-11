package main

import (
	"testing"

	"github.com/panduit-joeb/jkv/store/fs"
	"github.com/stretchr/testify/assert"
)

func TestHGET(t *testing.T) {
	t.Run("Test HGET", func(t *testing.T) {
		f := fs.NewJKVClient()
		f.Open()
		value, err := f.HGET("other", "one")
		assert.Nil(t, err)
		assert.Equal(t, "value", value)
	})
}
