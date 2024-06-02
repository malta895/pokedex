//go:build integration
// +build integration

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"malta895/pokedex/testutils"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	t.Setenv("HTTP_PORT", "5000")

	// start server in a goroutine
	go func() {
		main()
	}()
	// wait for the server to start
	for {
		if _, err := http.Get("http://localhost:5000"); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	t.Run("GET /pokemon/{pokemonName}", func(t *testing.T) {
		testCases := map[string]struct {
			pokemonName string

			expectedResp       string
			expectedStatusCode int
		}{
			"should respond with 200 OK with mewtwo": {
				pokemonName: "mewtwo",

				expectedResp: `{
				"name": "mewtwo",
				"description": "It was created by\na scientist after\nyears of horrific\fgene splicing and\nDNA engineering\nexperiments.",
				"habitat": "rare",
				"isLegendary": true
			}`,
				expectedStatusCode: http.StatusOK,
			},
			"should respond 404 Not Found with nonExistingPokemon": {
				pokemonName:        "nonExistingPokemon",
				expectedResp:       "Not Found",
				expectedStatusCode: http.StatusNotFound,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("http://localhost:5000/pokemon/%s", tt.pokemonName))
				if err != nil {
					t.Errorf("error making request to server: %s", err)
					return
				}

				if tt.expectedStatusCode != resp.StatusCode {
					t.Errorf("found statusCode=%d; want %d", resp.StatusCode, tt.expectedStatusCode)
				}

				respBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %s", err)
				}
				foundResp := string(respBytes)
				if json.Valid([]byte(tt.expectedResp)) {
					contentType := resp.Header.Get("Content-Type")
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
	})

	t.Run("GET /pokemon/translated/{pokemonName}", func(t *testing.T) {
		testCases := map[string]struct {
			pokemonName string

			expectedResp       string
			expectedStatusCode int
		}{
			"should respond with 200 OK with mewtwo": {
				pokemonName: "mewtwo",

				expectedResp: `{
				"name": "mewtwo",
				"description": "Created by a scientist after years of horrific gene splicing and dna engineering experiments,  it was.",
				"habitat": "rare",
				"isLegendary": true
			}`,
				expectedStatusCode: http.StatusOK,
			},
			"should respond 404 Not Found with nonExistingPokemon": {
				pokemonName:        "nonExistingPokemon",
				expectedResp:       "Not Found",
				expectedStatusCode: http.StatusNotFound,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				resp, err := http.Get(fmt.Sprintf("http://localhost:5000/pokemon/translated/%s", tt.pokemonName))
				if err != nil {
					t.Errorf("error making request to server: %s", err)
					return
				}

				if tt.expectedStatusCode != resp.StatusCode {
					t.Errorf("found statusCode=%d; want %d", resp.StatusCode, tt.expectedStatusCode)
				}

				respBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("error reading response body: %s", err)
				}
				foundResp := string(respBytes)
				if json.Valid([]byte(tt.expectedResp)) {
					contentType := resp.Header.Get("Content-Type")
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
	})

}
