package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	handler "github.com/ravip18596/order-food-online/internal/handler"
)

func main() {
	// New router
	r := mux.NewRouter()

	// Middleware
	r.Use(handlers.RecoveryHandler())
	r.Use(handlers.CompressHandler)

	// Init handler
	h := handler.NewHandler()
	h.RegisterRoutes(r)

	// HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Starting server
	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
