package cacher //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		for i := 0; i < 6; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
		// First added element should pop-out
		_, ok := c.Get("0")
		require.False(t, ok)

		c.Clear()
		// Last added element should disappear
		_, ok = c.Get("5")
		require.False(t, ok)

		// c is empty cache
		c.Set(Key("a"), 1)
		c.Set(Key("b"), 1)
		c.Set(Key("c"), 1)
		c.Set(Key("d"), 1)
		c.Set(Key("e"), 1) // [e d c b a]
		c.Set(Key("a"), 2) // [a e d c b]
		c.Set(Key("d"), 2) // [d a e c b]
		c.Set(Key("f"), 2) // [f d a e c]
		c.Set(Key("g"), 2) // [g f d a e]

		_, ok = c.Get("b")
		require.False(t, ok)
		_, ok = c.Get("c")
		require.False(t, ok)
	})

	t.Run("weird tests", func(t *testing.T) {
		c := NewCache(0)
		c.Set(Key("test"), 100)
		_, ok := c.Get("test")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
