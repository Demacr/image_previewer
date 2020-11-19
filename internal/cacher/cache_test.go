package cacher //nolint:golint,stylecheck

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	domain "github.com/Demacr/image_previewer/internal/domain/types"
	"github.com/stretchr/testify/require"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
	if err := os.Chdir("../../"); err != nil {
		panic(err)
	}
}

func getImage(t *testing.T, path string) io.ReadCloser {
	source, err := os.Open(path)
	require.NoError(t, err)
	return source
}

func isEqualReadClosers(t *testing.T, rc1, rc2 io.ReadCloser) bool {
	buf1 := bytes.Buffer{}
	buf2 := bytes.Buffer{}
	_, err1 := io.Copy(&buf1, rc1)
	_, err2 := io.Copy(&buf2, rc2)
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, rc1.Close())
	require.NoError(t, rc2.Close())
	return bytes.Equal(buf1.Bytes(), buf2.Bytes())
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c, err := NewCache(10)
		require.NoError(t, err)

		_, err = c.Get("aaa")
		require.NoError(t, err)

		_, err = c.Get("bbb")
		require.NoError(t, err)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5)
		require.NoError(t, err)

		err = c.Set("aaa", &domain.DownloadedImage{
			Buffer: getImage(t, "test/unit/source/image.jpg"),
			Header: nil,
		})
		require.NoError(t, err)

		err = c.Set("bbb", &domain.DownloadedImage{
			Buffer: getImage(t, "test/unit/source/image.jpg"),
			Header: nil,
		})
		require.NoError(t, err)

		val, err := c.Get("aaa")
		require.NoError(t, err)
		require.NotNil(t, val)
		require.True(t, isEqualReadClosers(t, val.Buffer, getImage(t, "test/unit/source/image.jpg")))

		val, err = c.Get("bbb")
		require.NoError(t, err)
		require.NotNil(t, val)
		require.True(t, isEqualReadClosers(t, val.Buffer, getImage(t, "test/unit/source/image.jpg")))

		val, err = c.Get("ccc")
		require.NoError(t, err)
		require.Nil(t, val)
	})
}
