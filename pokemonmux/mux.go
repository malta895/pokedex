package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
	"malta895/pokedex/pokeapi"
	"net/http"
)

const pokemonNamePathWildcard = "pokemonName"

func New(
	logger *log.Logger,
	pokeAPIClient pokeapi.Client,
) *http.ServeMux {
	serveMux := http.NewServeMux()

	// Endpoint 1: Basic Pokemon Information
	serveMux.HandleFunc(
		fmt.Sprintf("GET /pokemon/{%s}", pokemonNamePathWildcard),
		func(w http.ResponseWriter, r *http.Request) {
			pokemonName := r.PathValue(pokemonNamePathWildcard)
			pokemon, err := pokeAPIClient.PokemonByName(pokemonName)
			if err == pokeapi.ErrPokemonNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`Not Found`))
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`Internal Server Error`))
				return
			}
			respBody, err := json.Marshal(pokemon)
			if err != nil {

			}
			w.Write(respBody)
		})

	return serveMux
}
