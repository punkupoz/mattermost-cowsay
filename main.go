package main

import (
	"context"
	"io/ioutil"

	//"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v2"
	//"github.com/jinzhu/gorm"
)

type server struct {
	//db *gorm.DB
	router *chi.Mux
	config *Config
}

type Config struct {
	Mattermost struct{
		Token string
	}
}

// loadConfig receive a path and return config object
func loadConfig(path string) (*Config, error) {
	var config Config
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(configFile), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func newServer() *server {
	r := chi.NewRouter()
	config, err := loadConfig("./conf.yaml")
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	srv := server{
		router: r,
		config: config,
	}

	srv.routes()

	return &srv
}

func main() {
	s := newServer()
	srv := &http.Server{
		Addr:    ":8080",
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
