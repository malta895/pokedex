package pokemon

import (
	"encoding/json"
	"io"
	"net/http"
)

type PokemonClient struct {
	baseUrl string
}

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"isLegendary"`
}

func NewClient(baseUrl string) PokemonClient {
	return PokemonClient{baseUrl}
}

func (p PokemonClient) PokemonByName(name string) (Pokemon, error) {
	resp, err := http.Get(p.baseUrl + "/" + name)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonSpecies := pokemonSpecies{}
	err = json.Unmarshal(respBytes, &pokemonSpecies)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
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
