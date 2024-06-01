package funtranslations

import (
	"errors"
	"fmt"
	"io"
	"malta895/pokedex/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFunTranslate(t *testing.T) {
	tests := map[string]struct {
		translatorType  string
		inputText       string
		mockAPIResponse string
		nonOKStatusCode int

		expectedTranslatorPath string
		expectedTranslation    string
		expectedError          error
		expectApiCalled        bool
	}{
		"should respond with correct yoda translation": {
			translatorType: TranslatorYoda,
			inputText:      "this is some translation",
			mockAPIResponse: `{
				"success": {
				  "total": 1
				},
				"contents": {
				  "translated": "Some translation, this is",
				  "text": "this is some translation",
				  "translation": "yoda"
				}
			  }`,

			expectedTranslatorPath: yodaPath,
			expectedTranslation:    "Some translation, this is",
			expectedError:          nil,
			expectApiCalled:        true,
		},
		"should respond with correct shakespeare translation": {
			translatorType: TranslatorShakespeare,
			inputText:      "You are Mr. Luca",
			mockAPIResponse: `{
				"success": {
				  "total": 1
				},
				"contents": {
				  "translated": "Ye art mr. Luca",
				  "text": "You are Mr. Luca",
				  "translation": "shakespeare"
				}
			  }`,

			expectedTranslatorPath: shakespearePath,
			expectedTranslation:    "Ye art mr. Luca",
			expectedError:          nil,
			expectApiCalled:        true,
		},
		"should return correct error if translator type is not recognized": {
			translatorType: "unknownTranslatorType",
			inputText:      "You are Mr. Luca",

			expectedTranslatorPath: "",
			expectedError:          ErrUnrecognizedTranslator,
			expectApiCalled:        false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var apiCalled bool
			var statusCode int = http.StatusOK
			if tt.nonOKStatusCode != 0 {
				statusCode = tt.nonOKStatusCode
			}
			server := mockFunTranslationsServer(
				t,
				tt.expectedTranslatorPath,
				tt.mockAPIResponse,
				tt.inputText,
				statusCode,
				func() bool {
					apiCalled = true
					return tt.expectApiCalled
				},
			)
			defer server.Close()
			pokemonClient := &client{server.URL}

			foundTranslation, err := pokemonClient.FunTranslate(tt.translatorType, tt.inputText)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf(
					"received error %v; want %v",
					err,
					tt.expectedError,
				)
			}
			if tt.expectedTranslation != strings.TrimSpace(foundTranslation) {
				t.Errorf("found translation %s; want %s", foundTranslation, tt.expectedTranslation)
			}
			if tt.expectApiCalled != apiCalled {
				t.Errorf("found apiCalled=%v; want %v", apiCalled, tt.expectApiCalled)
			}
		})
	}
}

func mockFunTranslationsServer(
	t *testing.T,
	translatorPath, mockResp, expectedInputText string,
	statusCode int,
	assertCalled func() bool,
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assertCalled() {
			return
		}
		expectedPath := "/" + translatorPath
		if r.URL.Path != expectedPath {
			t.Errorf("Expected to request %s, got: %s", expectedPath, r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("received req method %s; want %s", r.Method, http.MethodPost)
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("error while reading req body: %s", err)
		}
		expectedBody := fmt.Sprintf(`{"text":"%s"}`, expectedInputText)
		ok, err := testutils.JsonEq(string(bodyBytes), expectedBody)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Errorf("found req body %s; want %s", string(bodyBytes), expectedBody)
		}
		w.WriteHeader(statusCode)

		_, err = w.Write([]byte(mockResp))
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
