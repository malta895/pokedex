package pokemon

// pokeAPIPokemonSpeciesURL is the complete URL of the pokemon-species pokeapi endpoint
//
// Reference: https://pokeapi.co/docs/v2#pokemon-species
const pokeAPIPokemonSpeciesURL = "https://pokeapi.co/api/v2/pokemon-species"

// pokemonSpecies is a partial representation of the `PokemonSpecies` pokeapi type
//
// Reference: https://pokeapi.co/docs/v2#pokemonspecies
type pokemonSpecies struct {
	Name              string         `json:"name"`
	FlavorTextEntries []flavorText   `json:"flavor_text_entries"`
	Habitat           pokemonHabitat `json:"habitat"`
	IsLegendary       bool           `json:"is_legendary"`
}

// flavorText is a partial representation of the `FlavorText` pokeapi type
//
// Reference: https://pokeapi.co/docs/v2#flavortext
type flavorText struct {
	FlavorText string `json:"flavor_text"`
	Language   struct {
		Iso3166 string `json:"iso3166"`
	} `json:"language"`
}

// pokemonHabitat is a partial representation of the `PokemonHabitat` pokeapi type
//
// Reference: https://pokeapi.co/docs/v2#pokemonhabitat
type pokemonHabitat struct {
	Name string `json:"name"`
}
