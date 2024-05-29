package pokemon

type Pokemon struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Habitat     string `json:"habitat"`
	IsLegendary bool   `json:"isLegendary"`
}

func PokemonByName(name string) (Pokemon, error) {
	return Pokemon{
		Name:        "ditto",
		Description: "Capable of copying\nan enemy's genetic\ncode to instantly\ftransform itself\ninto a duplicate\nof the enemy.",
		Habitat:     "urban",
		IsLegendary: false,
	}, nil
}
