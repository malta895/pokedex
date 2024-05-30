package pokeapiclient

import (
	"encoding/json"
	"io"
	"malta895/pokedex/types"
	"net/http"
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
		return types.Pokemon{}, nil
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Pokemon{}, nil
	}

	pokemonSpecies := pokemonSpecies{}
	err = json.Unmarshal(respBytes, &pokemonSpecies)
	if err != nil {
		return types.Pokemon{}, nil
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
