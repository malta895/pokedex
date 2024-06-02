package testutils

import (
	"strings"
	"testing"
)

func TestJsonEq(t *testing.T) {
	tests := map[string]struct {
		foundBody string

		expectedBody        string
		expected            bool
		expectedErrorSubstr string
	}{
		"should report equal json": {
			foundBody:    `{"name": "Charizard", "isLegendary": false}`,
			expectedBody: `{"name": "Charizard", "isLegendary": false}`,
			expected:     true,
		},
		"should report equal json with different order": {
			foundBody:    `{"isLegendary": false, "name": "Charizard"}`,
			expectedBody: `{"name": "Charizard", "isLegendary": false}`,
			expected:     true,
		},
		"should report not equal json": {
			foundBody:    `{"name": "Charizard", "isLegendary": false}`,
			expectedBody: `{"name": "Charizard", "isLegendary": true}`,
			expected:     false,
		},
		"should report not equal json with different order": {
			foundBody:    `{"isLegendary": false, "name": "Charizard"}`,
			expectedBody: `{"name": "Charmeleon", "isLegendary": true}`,
			expected:     false,
		},
		"should return error if invalid json in found body": {
			foundBody:    `{"name": "Charizard", "isLegendary": false`,
			expectedBody: `{"name": "Charizard", "isLegendary": false}`,

			expectedErrorSubstr: "json syntax error in found body:",
		},

		"should return error invalid json in expected body": {
			foundBody:    `{"name": "Charizard", "isLegendary": false}`,
			expectedBody: `{"name": "Charizard", "isLegendary": false`,
			expected:     false,

			expectedErrorSubstr: "json syntax error in expected body:",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := JsonEq(tt.foundBody, tt.expectedBody)
			if err != nil && tt.expectedErrorSubstr == "" {
				t.Fatalf("unexpected error: %v", err)
			}

			if err != nil && tt.expectedErrorSubstr != "" {
				if !strings.Contains(err.Error(), tt.expectedErrorSubstr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErrorSubstr, err)
				}
			}

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
