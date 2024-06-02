package main

import (
	"context"
	"fmt"
	"log"
	"malta895/pokedex/apiclients/funtranslations"
	"malta895/pokedex/apiclients/pokeapi"
	"malta895/pokedex/pokemonmux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := log.Default()
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3000"
		logger.Printf("HTTP_PORT not set, defaulting to %s", httpPort)
	}

	pokeapiClient := pokeapi.NewClient()
	funtranslationsClient := funtranslations.NewClient()
	pokemonMux := pokemonmux.New(
		logger,
		pokeapiClient,
		funtranslationsClient,
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", httpPort),
		Handler: pokemonMux,
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s", err)
		}
	}()
	logger.Printf("Server started on port %s", httpPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %s", err)
	}

	logger.Println("Server shut down")
}
