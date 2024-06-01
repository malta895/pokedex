package funtranslations

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestPokemonByName(t *testing.T) {
// 	tests := map[string]struct {
// 		pokemonName         string
// 		mockPokeAPIResponse string
// 		nonOKStatusCode     int

// 		expectedPokemon *types.Pokemon
// 		expectedError   error
// 		expectApiCalled bool
// 	}{}

// 	for name, tt := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			var apiCalled bool
// 			var statusCode int = http.StatusOK
// 			if tt.nonOKStatusCode != 0 {
// 				statusCode = tt.nonOKStatusCode
// 			}
// 			server := mockPokeAPIServer(
// 				t,
// 				tt.pokemonName,
// 				tt.mockPokeAPIResponse,
// 				statusCode,
// 				func() {
// 					apiCalled = true
// 				},
// 			)
// 			defer server.Close()
// 			pokemonClient := &client{server.URL}

// 			foundResp, err := pokemonClient.FunTranslate(tt.pokemonName)
// 			if err != tt.expectedError {
// 				t.Errorf(
// 					"received error %v; want %v",
// 					err,
// 					tt.expectedError,
// 				)
// 			}
// 			if !reflect.DeepEqual(foundResp, tt.expectedPokemon) {
// 				t.Errorf(
// 					"PokemonByName(%s) = %#v; want %#v",
// 					tt.pokemonName,
// 					foundResp,
// 					tt.expectedPokemon,
// 				)
// 			}
// 			if tt.expectApiCalled != apiCalled {
// 				t.Errorf("found apiCalled=%v; want %v", apiCalled, tt.expectApiCalled)
// 			}
// 		})
// 	}
// }

func mockFunTranslationsServer(
	t *testing.T,
	translatorPath string,
	mockResp string,
	statusCode int,
	assertCalled func(),
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCalled()
		expectedPath := funtranslationsBaseURL + "/" + translatorPath
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
	t.Run("new should return a client with the default funtranslations url", func(t *testing.T) {
		found := NewClient().(*client).baseURL

		if found != funtranslationsBaseURL {
			t.Errorf("unexpected baseUrl %s; want %s", found, funtranslationsBaseURL)
		}
	})
}
