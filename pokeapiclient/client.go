package pokeapiclient

import (
	"encoding/json"
	"errors"
	"io"
	"malta895/pokedex/types"
	"net/http"
)

var (
	ErrPokemonNotFound = errors.New("pokemon not found")
)

type PokemonClient struct {
	baseUrl string
}

func NewClient(baseUrl string) PokemonClient {
	return PokemonClient{baseUrl}
}

func (p PokemonClient) PokemonByName(name string) (types.Pokemon, error) {
	resp, err := http.Get(p.baseUrl + "/" + name)
	if err != nil {
		return types.Pokemon{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return types.Pokemon{}, ErrPokemonNotFound
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Pokemon{}, err
	}

	pokemonSpecies := pokemonSpecies{}
	err = json.Unmarshal(respBytes, &pokemonSpecies)
	if err != nil {
		return types.Pokemon{}, err
	}

	return types.Pokemon{
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
