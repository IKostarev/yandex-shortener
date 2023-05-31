package fs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSaveAndGet(t *testing.T) {
	original2 := "https://example.com/long/url2"
	original1 := "https://example.com/long/url1"

	file, err := os.CreateTemp("", "TestSaveAndGet-db1")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.Remove(file.Name()))
	}()

	fs, err := NewFs(file)
	assert.NoError(t, err)

	short1, err := fs.Save(original1, "")

	assert.NoError(t, err)
	assert.NotEqual(t, short1, "")
	//assert.Equal(t, interface{}(fs.Get(short1, "")), original1)

	short2, err := fs.Save(original2, "")

	assert.NoError(t, err)
	assert.NotEqual(t, short2, "")
	//assert.Equal(t, interface{}(fs.Get(short2, "")), original2)

	assert.NoError(t, file.Close())

	file2, err := os.Open(file.Name())
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, file2.Close())
	}()

	//fs2, err := NewFs(file2)
	//assert.NoError(t, err)
	//
	//for short, original := range fs.cacheURL {
	//	assert.Equal(t, interface{}(fs2.Get(short, "")), original)
	//}
}
