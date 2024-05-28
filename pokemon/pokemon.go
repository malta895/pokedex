package pokemon

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"isLegendary"`
}

// PokeAPIBody describes the GraphQL body sent to the PokéAPI
// `pokemon-species` API
type PokeAPIBody struct {
}

// PokeAPIResponse describes the GraphQL response sent by the PokéAPI
// `pokemon-species` API.
type PokeAPIResponse struct {
}

func PokemonByName(name string) (*Pokemon, error) {
	return nil, nil
}
