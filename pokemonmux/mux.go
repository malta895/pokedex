package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
	"malta895/pokedex/apiclients/funtranslations"
	"malta895/pokedex/apiclients/pokeapi"
	"malta895/pokedex/types"
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
		buildPokemonHandler(logger, pokeAPIClient, funtranslationsClient, false),
	)

	// Endpoint 2: Translated Pokemon Description
	serveMux.HandleFunc(
		fmt.Sprintf("GET /pokemon/translated/{%s}", pokemonNamePathWildcard),
		buildPokemonHandler(logger, pokeAPIClient, funtranslationsClient, true),
	)

	return serveMux
}

func buildPokemonHandler(
	logger *log.Logger,
	pokeAPIClient pokeapi.Client,
	funtranslationsClient funtranslations.Client,
	translateDescription bool,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pokemonName := r.PathValue(pokemonNamePathWildcard)

		pokemon, err := pokeAPIClient.PokemonByName(pokemonName)
		if err != nil {
			handlePokemonError(logger, w, "error retrieving pokemon", err)
			return
		}

		if translateDescription {
			translatePokemonDescription(logger, pokemon, funtranslationsClient)
		}

		writeResponse(logger, w, http.StatusOK, pokemon)
	}
}

func translatePokemonDescription(
	logger *log.Logger,
	pokemon *types.Pokemon,
	funtranslationsClient funtranslations.Client,
) {
	translatorType := funtranslations.TranslatorShakespeare
	if pokemon.IsLegendary || pokemon.Habitat == "cave" {
		translatorType = funtranslations.TranslatorYoda
	}
	translatedDesc, err := funtranslationsClient.FunTranslate(translatorType, pokemon.Description)
	if err != nil {
		logger.Printf("error translating description for pokemon %s: %v", pokemon.Name, err)
		return
	}
	pokemon.Description = translatedDesc
}

func handlePokemonError(
	logger *log.Logger,
	w http.ResponseWriter,
	message string,
	err error,
) {
	logger.Printf("%s: %v", message, err)
	if err == pokeapi.ErrPokemonNotFound {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func writeResponse(
	logger *log.Logger,
	w http.ResponseWriter,
	statusCode int,
	data interface{},
) {
	respBody, err := json.Marshal(data)
	if err != nil {
		handlePokemonError(logger, w, "error marshalling response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}
