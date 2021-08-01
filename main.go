package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yanpozka/airvet-jwt/api"
	"github.com/yanpozka/airvet-jwt/dao"
)

const (
	dbPath = "users.db"

	shutdownTimeout   = 7 * time.Second
	readTimeout       = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	writeTimeout      = 15 * time.Second
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	d, err := dao.NewDAO(dbPath)
	if err != nil {
		log.Panic(err)
	}
	err = d.InitDB(context.Background())
	if err != nil {
		log.Panic(err)
	}

	a := api.NewAPI(d)

	addr := ":" + getEnvStr("PORT", "8080")
	srv := &http.Server{
		Addr:    addr,
		Handler: a.GetRoutes(),

		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)

	go func() {
		log.Printf("Serving on %q ...", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// blocks until we get a kill/interrupt OS signal
	select {
	case osSignal := <-ch:
		log.Printf("Got OS signal: '%v', shuting down the server with timeout: %v", osSignal, shutdownTimeout)
	}

	log.Println("Closing db ...")
	d.Close()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	srv.SetKeepAlivesEnabled(false)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown the server: %v", err)
	}
}

func getEnvStr(name, defaultVal string) string {
	envVal := os.Getenv(name)
	if envVal == "" {
		return defaultVal
	}
	return envVal
}
