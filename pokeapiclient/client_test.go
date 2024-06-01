package pokeapiclient

import (
	"malta895/pokedex/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPokemonByName(t *testing.T) {
	tests := map[string]struct {
		pokemonName         string
		mockPokeAPIResponse string
		nonOKStatusCode     int

		expectedPokemon *types.Pokemon
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

			expectedPokemon: &types.Pokemon{
				Name:        "fakelegend",
				Description: "This is a mock legendary pokemon",
				Habitat:     "mockHabitat",
				IsLegendary: true,
			},
			expectedError:   nil,
			expectApiCalled: true,
		},
		"should respond with the first english description": {
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
				  },
				  {
					"flavor_text": "This is the second english description of a mock big pokemon",
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

			expectedPokemon: &types.Pokemon{
				Name:        "bigpokemon",
				Description: "This is a mock big pokemon",
				Habitat:     "mockHabitat",
				IsLegendary: false,
			},
			expectedError:   nil,
			expectApiCalled: true,
		},
		"should respond with the pokemon not found error if the api responds 404": {
			pokemonName:         "nonexisting",
			mockPokeAPIResponse: `Not Found`,
			nonOKStatusCode:     http.StatusNotFound,

			expectedPokemon: nil,
			expectedError:   ErrPokemonNotFound,
			expectApiCalled: true,
		},
		"should respond with a generic error if the api responds with a non-404 error": {
			pokemonName:         "nonexisting",
			mockPokeAPIResponse: `Bad Request`,
			nonOKStatusCode:     http.StatusBadRequest,

			expectedPokemon: nil,
			expectedError:   ErrUnknown,
			expectApiCalled: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var apiCalled bool
			var statusCode int = http.StatusOK
			if tt.nonOKStatusCode != 0 {
				statusCode = tt.nonOKStatusCode
			}
			server := mockPokeAPIServer(
				t,
				tt.pokemonName,
				tt.mockPokeAPIResponse,
				statusCode,
				func() {
					apiCalled = true
				},
			)
			defer server.Close()
			pokemonClient := &Client{server.URL}

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
	statusCode int,
	assertCalled func(),
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCalled()
		expectedPath := pokemonSpeciesPath + "/" + pokemonName
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request %s, got: %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockResp))
		if err != nil {
			t.Errorf("Expect nil err, got %s", err)
		}
	}))
}

func TestNew(t *testing.T) {
	t.Run("new should return a client with the default pokeapi url", func(t *testing.T) {
		found := New().baseUrl

		if found != pokeAPIBaseURL {
			t.Errorf("unexpected baseUrl %s; want %s", found, pokeAPIBaseURL)
		}
	})
}
