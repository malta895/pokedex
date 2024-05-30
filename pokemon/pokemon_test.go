package pokemon

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPokemonByName(t *testing.T) {
	tests := map[string]struct {
		pokemonName         string
		mockPokeAPIResponse string

		expectedPokemon Pokemon
		expectedError   error
		expectApiCalled bool
	}{
		"should respond with expected body with a mock legendary pokemon": {
			pokemonName: "fakelegend",
			mockPokeAPIResponse: `{
				"flavor_text_entries": [
				  {
					"flavor_text": "This is a mock legendary pokemon",
					"language": {
					  "name": "en"
					}
				  }
				],
				"habitat": {
				  "name": "mockHabitat"
				},
				"is_legendary": true,
				"name": "fakelegend"
			  }
			`,

			expectedPokemon: Pokemon{
				Name:        "fakelegend",
				Description: "This is a mock legendary pokemon",
				Habitat:     "mockHabitat",
				IsLegendary: true,
			},
			expectedError:   nil,
			expectApiCalled: true,
		},
		"should respond with the english description": {
			pokemonName: "bigpokemon",
			mockPokeAPIResponse: `{
				"flavor_text_entries": [
					{
						"flavor_text": "Questo è un pokemon grande di test",
						"language": {
						  "name": "it"
						}
					  },
				  {
					"flavor_text": "This is a mock big pokemon",
					"language": {
					  "name": "en"
					}
				  }
				],
				"habitat": {
				  "name": "mockHabitat"
				},
				"is_legendary": false,
				"name": "bigpokemon"
			  }
			`,

			expectedPokemon: Pokemon{
				Name:        "bigpokemon",
				Description: "This is a mock big pokemon",
				Habitat:     "mockHabitat",
				IsLegendary: false,
			},
			expectedError:   nil,
			expectApiCalled: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var apiCalled bool
			server := mockPokeAPIServer(
				t,
				tt.pokemonName,
				tt.mockPokeAPIResponse,
				func() {
					apiCalled = true
				},
			)
			defer server.Close()
			pokemonClient := NewClient(server.URL)

			foundResp, err := pokemonClient.PokemonByName(tt.pokemonName)
			if err != tt.expectedError {
				t.Errorf(
					"received error %v; want %v",
					err,
					tt.expectedError,
				)
			}
			if !reflect.DeepEqual(foundResp, tt.expectedPokemon) {
				t.Errorf(
					"PokemonByName(%s) = %#v; want %#v",
					tt.pokemonName,
					foundResp,
					tt.expectedPokemon,
				)
			}
			if tt.expectApiCalled != apiCalled {
				t.Errorf("found apiCalled=%v; want %v", apiCalled, tt.expectApiCalled)
			}
		})
	}
}

func mockPokeAPIServer(
	t *testing.T,
	pokemonName string,
	mockResp string,
	assertCalled func(),
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCalled()
		if r.URL.Path != ("/" + pokemonName) {
			t.Errorf("Expected to request %s, got: %s", pokeAPIPokemonSpeciesURL, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(mockResp))
		if err != nil {
			t.Errorf("Expect nil err, got %s", err)
		}
	}))
}
