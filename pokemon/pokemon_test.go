package pokemon

import "testing"

func TestPokemonByName(t *testing.T) {
	tests := map[string]struct {
		pokemonName         string
		mockPokeAPIResponse PokeAPIResponse

		expectedPokeAPIBody PokeAPIBody
		expectedPokemon     Pokemon
		expectedError       error
	}{
		"": {},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockPokeAPIServer(
				t,
				tt.mockPokeAPIResponse,
				func(foundBody PokeAPIBody) {
					if tt.expectedPokeAPIBody != foundBody {
						t.Errorf(
							"pokeAPI found body was %#v; want %#v",
							foundBody,
							tt.expectedPokeAPIBody,
						)
					}
				},
			)

			foundResp, err := PokemonByName(tt.pokemonName)
			if err != tt.expectedError {
				t.Errorf(
					"received error %s; want %s",
					err,
					tt.expectedError,
				)
			}
			if *foundResp != tt.expectedPokemon {
				t.Errorf(
					"PokemonByName(%s) = %#v; want %#v",
					tt.pokemonName,
					foundResp,
					tt.expectedPokeAPIBody,
				)
			}
		})
	}
}

func mockPokeAPIServer(t *testing.T, mockResp PokeAPIResponse, assertBody func(b PokeAPIBody)) {
	t.Helper()

}
