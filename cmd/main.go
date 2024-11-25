package main

import (
	"context"
	"log"
	"net/http"
	"online-offline-indicator/handlers"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
)

const serverPort = ":8080"

func main() {
	// Setup Redis DB
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		PoolSize:     100,
		MinIdleConns: 5,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init handlers
	h := handlers.NewHandler(rdb)

	// Routes
	r := http.DefaultServeMux
	r.HandleFunc("GET /users", h.GetUserStatusByBatch)
	r.HandleFunc("POST /heartbeat", h.SetUserStatus)

	// Server configuration
	srv := &http.Server{
		Addr:    serverPort,
		Handler: r,
	}

	// Start server
	go func() {
		log.Println("Starting HTTP server on port", serverPort)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	<-exit

	log.Println("Shutting down the server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}

}
