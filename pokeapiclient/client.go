package pokeapiclient

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

type Client struct {
	baseUrl string
}

func New() *Client {
	return &Client{pokeAPIBaseURL}
}

func (p *Client) PokemonByName(name string) (*types.Pokemon, error) {
	resUrl, err := url.JoinPath(p.baseUrl, pokemonSpeciesPath, name)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(resUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPokemonNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnknown
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, err)
	}

	pokemonSpecies := pokemonSpecies{}
	err = json.Unmarshal(respBytes, &pokemonSpecies)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknown, err)
	}

	return &types.Pokemon{
		Name:        pokemonSpecies.Name,
		Description: retrieveEnglishDescription(pokemonSpecies),
		Habitat:     pokemonSpecies.Habitat.Name,
		IsLegendary: pokemonSpecies.IsLegendary,
	}, nil
}

func retrieveEnglishDescription(ps pokemonSpecies) string {
	for _, flavorTextEntry := range ps.FlavorTextEntries {
		if flavorTextEntry.Language.Name == "en" {
			return flavorTextEntry.FlavorText
		}
	}
	return ""
}
