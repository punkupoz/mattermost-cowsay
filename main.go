package main

import (
	"context"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	PORT = 8080
)

type server struct {
	repo *Repo
	router *chi.Mux
	config *Config
}

func main() {
	r := chi.NewRouter()
	repo, err := newRepo("sqlite3", "data.db")
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	config, err := loadConfigFile()
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	s := &server{
		repo: repo,
		router: r,
		config: config,
	}
	s.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: s.router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	waitForShutdown(srv)
}

// gracefully shutdown the server
func waitForShutdown(srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server Shutdown: ", err)
	}

	log.Println("server exiting")

	os.Exit(0)
}
