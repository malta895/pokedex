package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
	"malta895/pokedex/pokeapi"
	"malta895/pokedex/types"
	"net/http"
	"net/http/httptest"
	"reflect"
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

			expectedResp:       `Not Found`,
			expectedStatusCode: http.StatusNotFound,
		},
		"should respond with 500 Internal Server Error with unknown error": {
			mockPokeAPIClient: &mockPokeAPIClient{
				mockResp: nil,
				mockErr:  pokeapi.ErrUnknown,
			},
			pokemonName: "mewtwo",

			expectedResp:       `Internal Server Error`,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			handler := New(log.Default(), tt.mockPokeAPIClient)
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
				bodyOK, err := jsonEq(foundResp, tt.expectedResp)
				if err != nil {
					t.Error(err)
				}
				if !bodyOK {
					t.Errorf("found respBody=%s; want %s", foundResp, tt.expectedResp)
				}
			} else {
				if foundResp != tt.expectedResp {
					t.Errorf("found respBody=%s, want %s", foundResp, tt.expectedResp)
				}
			}

		})
	}

}

func jsonEq(foundBody, expectedBody string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(foundBody), &o1)
	if err != nil {
		return false, fmt.Errorf("json syntax error in found body: %s", err)
	}
	err = json.Unmarshal([]byte(expectedBody), &o2)
	if err != nil {
		return false, fmt.Errorf("json syntax error in expected body: %s", err)
	}

	return reflect.DeepEqual(o1, o2), nil
}
