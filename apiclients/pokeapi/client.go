package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"malta895/pokedex/types"
	"net/http"
	"net/url"
)

var (
	ErrPokemonNotFound = errors.New("pokemon not found")
	ErrUnknown         = errors.New("cannot retrieve pokemon due to unknown error")
)

const (
	// pokeAPIBaseURL is the base URL of the `pokeapi.co` v2 APIs
	pokeAPIBaseURL = "https://pokeapi.co/api/v2"

	// pokemonSpeciesPath is the complete URL of the pokemon-species pokeapi endpoint
	//
	// Reference: https://pokeapi.co/docs/v2#pokemon-species
	pokemonSpeciesPath = "/pokemon-species"
)

type client struct {
	baseUrl string
}

type Client interface {
	PokemonByName(name string) (*types.Pokemon, error)
}

func NewClient() Client {
	return &client{pokeAPIBaseURL}
}

func (p *client) PokemonByName(name string) (*types.Pokemon, error) {
	resUrl, err := url.JoinPath(p.baseUrl, pokemonSpeciesPath, name)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(resUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := mapStatusToErr(resp.StatusCode); err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, err)
	}

	pokemonSpecies := pokemonSpecies{}
	if err := json.Unmarshal(respBytes, &pokemonSpecies); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, err)
	}

	return &types.Pokemon{
		Name:        pokemonSpecies.Name,
		Description: retrieveFirstEnglishDescription(pokemonSpecies),
		Habitat:     pokemonSpecies.Habitat.Name,
		IsLegendary: pokemonSpecies.IsLegendary,
	}, nil
}

func mapStatusToErr(statusCode int) error {
	if statusCode == http.StatusNotFound {
		return ErrPokemonNotFound
	}
	if statusCode != http.StatusOK {
		return ErrUnknown
	}
	return nil
}

func retrieveFirstEnglishDescription(ps pokemonSpecies) string {
	for _, flavorTextEntry := range ps.FlavorTextEntries {
		if flavorTextEntry.Language.Name == "en" {
			return flavorTextEntry.FlavorText
		}
	}
	return ""
}
