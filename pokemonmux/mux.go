package pokemonmux

import (
	"log"
	"malta895/pokedex/pokeapi"
	"net/http"
)

func New(
	logger *log.Logger,
	pokeAPIClient pokeapi.Client,
) *http.ServeMux {
	serveMux := http.NewServeMux()

	// Endpoint 1: Basic Pokemon Information
	// serveMux.HandleFunc("GET /pokemon/{pokemonName}", func(w http.ResponseWriter, r *http.Request) {})

	return serveMux
}
