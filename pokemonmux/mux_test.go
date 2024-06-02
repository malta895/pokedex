package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
	"malta895/pokedex/apiclients/funtranslations"
	"malta895/pokedex/apiclients/pokeapi"
	"malta895/pokedex/testutils"
	"malta895/pokedex/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockPokeAPIClient struct {
	mockResp *types.Pokemon
	mockErr  error
}

func (mpc *mockPokeAPIClient) PokemonByName(name string) (*types.Pokemon, error) {
	return mpc.mockResp, mpc.mockErr
}

func TestBasicPokemonInfo(t *testing.T) {
	testCases := map[string]struct {
		mockPokeAPIClient *mockPokeAPIClient
		pokemonName       string

		expectedResp       string
		expectedStatusCode int
	}{
		"should respond with 200 OK and correct information if client gives it": {
			mockPokeAPIClient: &mockPokeAPIClient{
				mockResp: &types.Pokemon{
					Name:        "mewtwo",
					Description: "some description",
					Habitat:     "rare",
					IsLegendary: true,
				},
				mockErr: nil,
			},
			pokemonName: "mewtwo",

			expectedResp: `{
				"name": "mewtwo",
				"description": "some description",
				"habitat": "rare",
				"isLegendary": true
			}`,
			expectedStatusCode: http.StatusOK,
		},
		"should respond with 404 Not Found if pokemon is not found": {
			mockPokeAPIClient: &mockPokeAPIClient{
				mockResp: nil,
				mockErr:  pokeapi.ErrPokemonNotFound,
			},
			pokemonName: "mewtwo",

			expectedResp:       "Not Found",
			expectedStatusCode: http.StatusNotFound,
		},
		"should respond with 500 Internal Server Error with unknown error": {
			mockPokeAPIClient: &mockPokeAPIClient{
				mockResp: nil,
				mockErr:  pokeapi.ErrUnknown,
			},
			pokemonName: "mewtwo",

			expectedResp:       "Internal Server Error",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			handler := New(log.Default(), tt.mockPokeAPIClient, nil)
			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf("/pokemon/%s", tt.pokemonName),
				nil,
			)
			if err != nil {
				t.Errorf("found err=%s; want nil", err)
			}

			respRecorder := httptest.NewRecorder()
			handler.ServeHTTP(respRecorder, req)

			if tt.expectedStatusCode != respRecorder.Code {
				t.Errorf("found statusCode=%d; want %d", respRecorder.Code, tt.expectedStatusCode)
			}

			foundResp := respRecorder.Body.String()
			if json.Valid([]byte(tt.expectedResp)) {
				contentType := respRecorder.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("found response contentType=%s; want application/json", contentType)
				}
				bodyOK, err := testutils.JsonEq(foundResp, tt.expectedResp)
				if err != nil {
					t.Error(err)
				}
				if !bodyOK {
					t.Errorf("found respBody=%s; want %s", foundResp, tt.expectedResp)
				}
			} else {
				if strings.TrimSpace(foundResp) != tt.expectedResp {
					t.Errorf("found respBody=%s, want %s", foundResp, tt.expectedResp)
				}
			}

		})
	}
}

type mockFunTranslationsClient struct {
	mockResp            string
	mockErr             error
	foundTranslatorType string
	foundText           string
}

func (mft *mockFunTranslationsClient) FunTranslate(translatorType, text string) (string, error) {
	mft.foundTranslatorType = translatorType
	mft.foundText = text
	return mft.mockResp, mft.mockErr
}

func TestTranslatedPokemonInfo(t *testing.T) {
	testCases := map[string]struct {
		mockPokeAPIClient         *mockPokeAPIClient
		mockFunTranslationsClient *mockFunTranslationsClient
		pokemonName               string

		expectedTranslatorType string
		expectedText           string
		expectedResp           string
		expectedStatusCode     int
	}{
		"should respond with a yoda translation for a cave pokemon": {
			mockPokeAPIClient: &mockPokeAPIClient{
				mockResp: &types.Pokemon{
					Name:        "somecavepokemon",
					Habitat:     "cave",
					Description: "this is some cave pokemon",
					IsLegendary: false,
				},
				mockErr: nil,
			},
			mockFunTranslationsClient: &mockFunTranslationsClient{
				mockResp: "some cave pokemon, this is",
				mockErr:  nil,
			},
			pokemonName: "somecavepokemon",

			expectedTranslatorType: funtranslations.TranslatorYoda,
			expectedText:           "this is some cave pokemon",
			expectedResp: `{
				"name": "somecavepokemon",
				"description": "some cave pokemon, this is",
				"habitat": "cave",
				"isLegendary": false
			}`,
			expectedStatusCode: http.StatusOK,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			handler := New(log.Default(), tt.mockPokeAPIClient, tt.mockFunTranslationsClient)
			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf("/pokemon/translated/%s", tt.pokemonName),
				nil,
			)
			if err != nil {
				t.Errorf("found err=%s; want nil", err)
			}

			respRecorder := httptest.NewRecorder()
			handler.ServeHTTP(respRecorder, req)

			foundTranslator := tt.mockFunTranslationsClient.foundTranslatorType
			if foundTranslator != tt.expectedTranslatorType {
				t.Errorf("found translatorType=%s; want %s", foundTranslator, tt.expectedTranslatorType)
			}
			foundText := tt.mockFunTranslationsClient.foundText
			if foundText != tt.expectedText {
				t.Errorf("found text=%s; want %s", foundText, tt.expectedText)
			}

			if tt.expectedStatusCode != respRecorder.Code {
				t.Errorf("found statusCode=%d; want %d", respRecorder.Code, tt.expectedStatusCode)
			}

			foundResp := respRecorder.Body.String()
			if json.Valid([]byte(tt.expectedResp)) {
				contentType := respRecorder.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("found response contentType=%s; want application/json", contentType)
				}
				bodyOK, err := testutils.JsonEq(foundResp, tt.expectedResp)
				if err != nil {
					t.Error(err)
				}
				if !bodyOK {
					t.Errorf("found respBody=%s; want %s", foundResp, tt.expectedResp)
				}
			} else {
				if strings.TrimSpace(foundResp) != tt.expectedResp {
					t.Errorf("found respBody=%s, want %s", foundResp, tt.expectedResp)
				}
			}
		})
	}
}
