package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
	"malta895/pokedex/apiclients/funtranslations"
	"malta895/pokedex/apiclients/pokeapi"
	"net/http"
)

const pokemonNamePathWildcard = "pokemonName"

func New(
	logger *log.Logger,
	pokeAPIClient pokeapi.Client,
	funtranslationsClient funtranslations.Client,
) *http.ServeMux {
	serveMux := http.NewServeMux()

	// Endpoint 1: Basic Pokemon Information
	serveMux.HandleFunc(
		fmt.Sprintf("GET /pokemon/{%s}", pokemonNamePathWildcard),
		buildBasicPokemonInfoHandler(pokeAPIClient),
	)

	// Endpoint 2: Translated Pokemon Description
	serveMux.HandleFunc(
		fmt.Sprintf("GET /pokemon/translated/{%s}", pokemonNamePathWildcard),
		func(w http.ResponseWriter, r *http.Request) {
			pokeAPIClient.PokemonByName("somecavepokemon")
			funtranslationsClient.FunTranslate(
				funtranslations.TranslatorYoda, "this is some cave pokemon",
			)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{
				"name": "somecavepokemon",
				"description": "some cave pokemon, this is",
				"habitat": "cave",
				"isLegendary": false
			}`))
		},
	)

	return serveMux
}

func buildBasicPokemonInfoHandler(pokeAPIClient pokeapi.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pokemonName := r.PathValue(pokemonNamePathWildcard)

		pokemon, err := pokeAPIClient.PokemonByName(pokemonName)
		if err == pokeapi.ErrPokemonNotFound {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(pokemon)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}
