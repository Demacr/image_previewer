package cacher //nolint:golint,stylecheck

import (
	"io"
	"os"
	"sync"

	domain "github.com/Demacr/image_previewer/internal/domain/types"
	"github.com/pkg/errors"
)

type Cache interface {
	Set(key string, value *domain.DownloadedImage) error // Добавить значение в кэш по ключу
	Get(key string) (*domain.DownloadedImage, error)     // Получить значение из кэша по ключу
	clear()                                              // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[string]cacheItem
	mutex    sync.Mutex
}

func (lc *lruCache) Set(key string, value *domain.DownloadedImage) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	item, was := lc.items[key]
	if was {
		lc.queue.MoveToFront(item)
		lc.queue.Front().Value = queueItem{value, key}
	} else {
		lc.items[key] = lc.queue.PushFront(queueItem{value, key})
		if lc.queue.Len() >= lc.capacity {
			delete(lc.items, lc.queue.Back().Value.(queueItem).key)
			lc.queue.Remove(lc.queue.Back())
		}

		// fs

		fd, err := os.Create("cache/" + key)
		if err != nil {
			return errors.Wrap(err, "error during create cached image")
		}
		defer fd.Close()

		if _, err = io.Copy(fd, value.Buffer); err != nil {
			return errors.Wrap(err, "error during saving image")
		}
		value.Buffer.Close()
	}
	return nil
}

func (lc *lruCache) Get(key string) (value *domain.DownloadedImage, err error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	item, ok := lc.items[key]
	if ok {
		value = item.Value.(queueItem).value
		lc.queue.MoveToFront(item)
		// fs
		fd, err := os.Open("cache/" + key)
		if err != nil {
			return value, errors.Wrap(err, "error during open image file")
		}
		value.Buffer = fd
	}
	return
}

func (lc *lruCache) clear() {
	lc.queue = NewList()
	lc.items = map[string]cacheItem{}
}

type cacheItem *listItem
type queueItem struct {
	value *domain.DownloadedImage
	key   string
}

func NewCache(capacity int) (Cache, error) {
	if err := os.MkdirAll("cache", os.ModePerm); os.IsNotExist(err) {
		return nil, errors.Wrap(err, "error during creating cache folder")
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    map[string]cacheItem{},
	}, nil
}
