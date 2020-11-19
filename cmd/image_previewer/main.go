package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Demacr/image_previewer/internal/cacher"
	"github.com/Demacr/image_previewer/internal/config"
	"github.com/Demacr/image_previewer/internal/domain/previewer"
	"github.com/Demacr/image_previewer/internal/http/server"
)

func main() {
	fc, err := cacher.NewCache(10)
	if err != nil {
		log.Fatal(err)
	}
	if fc == nil {
		log.Fatal("empty filecache")
	}
	p := previewer.NewPreviewer(fc)
	wg := sync.WaitGroup{}
	quitCh := make(chan interface{})

	cfg, _ := config.Configure()

	router := server.NewRouter(p)
	server := &http.Server{
		Addr:         cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler:      router.RootHandler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	go func() {
		<-quitCh
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("error during shutdown HTTP server:", err)
		}
	}()

	<-sigc
	log.Println("get signal to shutdown service")
	close(quitCh)

	wg.Wait()
	log.Println("successfully shutdown")
}
