package cacher //nolint:golint,stylecheck

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

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

		_, err = c.get("aaa")
		require.NoError(t, err)

		_, err = c.get("bbb")
		require.NoError(t, err)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5)
		require.NoError(t, err)

		err = c.set("aaa", &DownloadedImage{
			Buffer: getImage(t, "test/unit/source/image.jpg"),
			Header: nil,
		})
		require.NoError(t, err)

		err = c.set("bbb", &DownloadedImage{
			Buffer: getImage(t, "test/unit/source/image.jpg"),
			Header: nil,
		})
		require.NoError(t, err)

		val, err := c.get("aaa")
		require.NoError(t, err)
		require.NotNil(t, val)
		require.True(t, isEqualReadClosers(t, val.Buffer, getImage(t, "test/unit/source/image.jpg")))

		val, err = c.get("bbb")
		require.NoError(t, err)
		require.NotNil(t, val)
		require.True(t, isEqualReadClosers(t, val.Buffer, getImage(t, "test/unit/source/image.jpg")))

		val, err = c.get("ccc")
		require.NoError(t, err)
		require.Nil(t, val)
	})

	// t.Run("purge logic", func(t *testing.T) {
	// 	c := NewCache(5)

	// 	for i := 0; i < 6; i++ {
	// 		c.Set(Key(strconv.Itoa(i)), i)
	// 	}
	// 	// First added element should pop-out
	// 	_, ok := c.get("0")
	// 	require.False(t, ok)

	// 	c.Clear()
	// 	// Last added element should disappear
	// 	_, ok = c.get("5")
	// 	require.False(t, ok)

	// 	// c is empty cache
	// 	c.Set(Key("a"), 1)
	// 	c.Set(Key("b"), 1)
	// 	c.Set(Key("c"), 1)
	// 	c.Set(Key("d"), 1)
	// 	c.Set(Key("e"), 1) // [e d c b a]
	// 	c.Set(Key("a"), 2) // [a e d c b]
	// 	c.Set(Key("d"), 2) // [d a e c b]
	// 	c.Set(Key("f"), 2) // [f d a e c]
	// 	c.Set(Key("g"), 2) // [g f d a e]

	// 	_, ok = c.Get("b")
	// 	require.False(t, ok)
	// 	_, ok = c.Get("c")
	// 	require.False(t, ok)
	// })

	// t.Run("weird tests", func(t *testing.T) {
	// 	c := NewCache(0)
	// 	c.Set(Key("test"), 100)
	// 	_, ok := c.Get("test")
	// 	require.False(t, ok)
	// })
}

// func TestCacheMultithreading(t *testing.T) {
// 	c := NewCache(10)
// 	wg := &sync.WaitGroup{}
// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < 1_000_000; i++ {
// 			c.Set(Key(strconv.Itoa(i)), i)
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < 1_000_000; i++ {
// 			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
// 		}
// 	}()

// 	wg.Wait()
// }
