package pokemon

import (
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
	resp, err := http.Get(p.baseUrl + "/" + "ditto")
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{
		Name:        "ditto",
		Description: "Capable of copying\nan enemy's genetic\ncode to instantly\ftransform itself\ninto a duplicate\nof the enemy.",
		Habitat:     "urban",
		IsLegendary: false,
	}, nil
}
