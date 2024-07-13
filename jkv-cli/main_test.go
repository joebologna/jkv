package main

import (
	"context"
	"testing"

	"github.com/panduit-joeb/jkv/store/fs"
	"github.com/stretchr/testify/assert"
)

func TestHGET(t *testing.T) {
	t.Run("Test HGET", func(t *testing.T) {
		ctx := context.Background()
		f := fs.NewClient(&fs.Options{})
		f.Open()
		f.FlushDB(ctx)
		f.HSet(ctx, "other", "one", "value")
		rec := f.HGet(ctx, "other", "one")
		assert.Nil(t, rec.Err())
		assert.Equal(t, "value", rec.Val())
	})
}
