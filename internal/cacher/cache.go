package cacher //nolint:golint,stylecheck

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Demacr/image_previewer/internal/previewer"
	"github.com/pkg/errors"
)

type Cache interface {
	GetImage(url string, width int, height int) (*DownloadedImage, error)
	set(key string, value *DownloadedImage) error // Добавить значение в кэш по ключу
	get(key string) (*DownloadedImage, bool)      // Получить значение из кэша по ключу
	clear()                                       // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[string]cacheItem
	mutex    sync.Mutex
}

type DownloadedImage struct {
	Buffer io.ReadCloser
	Header http.Header
}

func (lc *lruCache) GetImage(url string, width int, height int) (*DownloadedImage, error) {
	nameBytes := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", url, width, height)))
	key := base32.StdEncoding.EncodeToString(nameBytes[:])
	if result, ok := lc.get(key); ok {
		log.Println("get image from cache")
		return result, nil
	}
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during creating request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error during getting image")
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		return nil, errors.New("wrong Content-Type")
	}

	previwedImage, err := previewer.Preview(resp.Body, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "error during image previewing")
	}

	result := &DownloadedImage{
		Buffer: previwedImage,
		Header: resp.Header,
	}

	if err = lc.set(key, result); err != nil {
		return nil, errors.Wrap(err, "error during saving previewed image")
	}
	result, _ = lc.get(key)
	return result, nil
}

func (lc *lruCache) set(key string, value *DownloadedImage) error {
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
	}
	return nil
}

func (lc *lruCache) get(key string) (value *DownloadedImage, ok bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	item, ok := lc.items[key]
	if ok {
		value = item.Value.(queueItem).value
		lc.queue.MoveToFront(item)
		// fs
		fd, err := os.Open("cache/" + key)
		if err != nil {
			return // todo: return error
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
	value *DownloadedImage
	key   string
}

func NewCache(capacity int) (Cache, error) {
	if err := os.Mkdir("cache", os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "error during creating cache folder")
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    map[string]cacheItem{},
	}, nil
}
