package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/database"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/repository"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/route"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/service"
)

func main() {
	database.InitDB()

	profileRepository := repository.NewProfileRepository(database.DB)
	profileService := service.NewProfileService(profileRepository)
	profileRouter := route.SetupProfileRoutes(profileService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      profileRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	srvErrChan := make(chan error, 1)

	go func() {
		log.Printf("Server is starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- err
		}
	}()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutting down...")
	case err := <-srvErrChan:
		log.Printf("HTTP server setup error: %v", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	log.Println("Shutdown completed gracefully.")
}
