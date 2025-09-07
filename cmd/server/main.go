package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	db "github.com/ravip18596/order-food-online/internal/database"
	handler "github.com/ravip18596/order-food-online/internal/handler"
	repo "github.com/ravip18596/order-food-online/internal/repository"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Load coupon codes
	set := make(map[string][]int)
	err := loadSetFromFile("couponbase1", set, 1)
	if err != nil {
		log.Fatalf("Failed to load couponbase1: %v", err)
	}
	err = loadSetFromFile("couponbase2", set, 2)
	if err != nil {
		log.Fatalf("Failed to load couponbase2: %v", err)
	}
	err = loadSetFromFile("couponbase3", set, 3)
	if err != nil {
		log.Fatalf("Failed to load couponbase3: %v", err)
	}
	fmt.Println("Coupon codes loaded into sets")

	// Initialize repositories
	productRepo := repo.NewProductRepository(db.DB)
	orderRepo := repo.NewOrderRepository(db.DB)

	// Initialize handler with repositories
	h := handler.NewHandler(productRepo, orderRepo, set)

	// Create a new router
	r := mux.NewRouter()

	// Middleware
	r.Use(handlers.RecoveryHandler())
	r.Use(handlers.CompressHandler)

	// Register routes
	h.RegisterRoutes(r)

	// HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
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

func loadSetFromFile(path string, uniqueWords map[string][]int, fileNo int) error {
	// A mutex to protect the uniqueWords map from concurrent access by multiple goroutines.
	var mu sync.Mutex
	var wg sync.WaitGroup

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// A buffered channel
	lineChan := make(chan string, 100)

	// --- PRODUCER GOROUTINE ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(lineChan) // Close the channel when reading is done

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineChan <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Scanner error: %v", err)
		}
	}()

	// --- CONSUMER GOROUTINES ---
	numWorkers := 3 // Adjust based on your system's CPU cores
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Process lines from the channel until it's closed
			for word := range lineChan {
				// using a mutex to prevent race conditions
				mu.Lock()
				if _, ok := uniqueWords[word]; !ok {
					uniqueWords[word] = append(uniqueWords[word], fileNo)
				} else {
					fileNos, _ := uniqueWords[word]
					if !containsInt(fileNos, fileNo) {
						uniqueWords[word] = append(uniqueWords[word], fileNo)
					}
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Set loaded from path:" + path)
	return nil
}

func containsInt(slice []int, v int) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}
