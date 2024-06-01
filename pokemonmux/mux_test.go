package pokemonmux

import (
	"encoding/json"
	"fmt"
	"log"
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
					Name: "somepokemon",
				},
				mockErr: nil,
			},
			pokemonName: "somepokemon",

			expectedResp: `{
				"name": "mewtwo",
				"description": "some description",
				"habitat": "rare",
				"isLegendary": true
			}`,
			expectedStatusCode: http.StatusOK,
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
				t.Errorf("found statusCode=%d; want %d", tt.expectedStatusCode, respRecorder.Code)
			}

			foundResp := respRecorder.Body.String()
			bodyOK, err := jsonEq(foundResp, tt.expectedResp)
			if err != nil {
				t.Error(err)
			}
			if !bodyOK {
				t.Errorf("found respBody=%s; want %s", foundResp, tt.expectedResp)
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
