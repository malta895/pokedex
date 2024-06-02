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
		expectAPICalled        bool
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
			expectAPICalled:        true,
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
			expectAPICalled:        true,
		},
		"should return correct error if translator type is not recognized": {
			translatorType: "unknownTranslatorType",
			inputText:      "You are Mr. Luca",

			expectedTranslatorPath: "",
			expectedError:          ErrUnrecognizedTranslator,
			expectAPICalled:        false,
		},
		"should respond with error if the api responds with 429 error": {
			translatorType: TranslatorShakespeare,
			inputText:      "You are Mr. Luca",
			mockAPIResponse: `{
				"error": {
					"code": 429,
					"message": "Too Many Requests: Rate limit of 10 requests per hour exceeded. Please wait for 44 minutes and 58 seconds."
				  }
			  }`,
			nonOKStatusCode: http.StatusTooManyRequests,

			expectedTranslatorPath: shakespearePath,
			expectedTranslation:    "",
			expectedError:          ErrAPIStatusCode,
			expectAPICalled:        true,
		},
		"should respond with error if the api responds with 500 error": {
			translatorType: TranslatorShakespeare,
			inputText:      "You are Mr. Luca",
			mockAPIResponse: `{
				"error": {
					"code": 500,
					"message": "Internal Server Error"
				  }
			  }`,
			nonOKStatusCode: http.StatusInternalServerError,

			expectedTranslatorPath: shakespearePath,
			expectedTranslation:    "",
			expectedError:          ErrAPIStatusCode,
			expectAPICalled:        true,
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
				func() {
					apiCalled = true
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
			if tt.expectAPICalled != apiCalled {
				t.Errorf("found apiCalled=%v; want %v", apiCalled, tt.expectAPICalled)
			}
		})
	}
}

func mockFunTranslationsServer(
	t *testing.T,
	translatorPath, mockResp, expectedInputText string,
	statusCode int,
	assertCalled func(),
) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertCalled()
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
			t.Errorf("unexpected baseURL %s; want %s", found, funtranslationsBaseURL)
		}
	})
}
