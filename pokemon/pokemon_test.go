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
		mockPokeAPIResponse pokemonSpecies

		expectedPokemon Pokemon
		expectedError   error
	}{
		"should respond with correct fields when response is 200 OK": {
			pokemonName: "ditto",
			expectedPokemon: Pokemon{
				Name:        "ditto",
				Description: "Capable of copying\nan enemy's genetic\ncode to instantly\ftransform itself\ninto a duplicate\nof the enemy.",
				Habitat:     "urban",
				IsLegendary: false,
			},
			expectedError: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			mockPokeAPIServer(
				t,
				tt.mockPokeAPIResponse,
			)

			foundResp, err := PokemonByName(tt.pokemonName)
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
		})
	}
}

const mockOKResp = `{
	"base_happiness": 50,
	"capture_rate": 35,
	"color": {
	  "name": "purple",
	  "url": "https://pokeapi.co/api/v2/pokemon-color/7/"
	},
	"egg_groups": [
	  {
		"name": "ditto",
		"url": "https://pokeapi.co/api/v2/egg-group/13/"
	  }
	],
	"evolution_chain": {
	  "url": "https://pokeapi.co/api/v2/evolution-chain/66/"
	},
	"evolves_from_species": null,
	"flavor_text_entries": [
	  {
		"flavor_text": "Capable of copying\nan enemy's genetic\ncode to instantly\ftransform itself\ninto a duplicate\nof the enemy.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "red",
		  "url": "https://pokeapi.co/api/v2/version/1/"
		}
	  },
	  {
		"flavor_text": "Capable of copying\nan enemy's genetic\ncode to instantly\ftransform itself\ninto a duplicate\nof the enemy.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "blue",
		  "url": "https://pokeapi.co/api/v2/version/2/"
		}
	  },
	  {
		"flavor_text": "Its transformation\nability is per­\nfect. However, if\fmade to laugh, it\ncan't maintain its\ndisguise.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "silver",
		  "url": "https://pokeapi.co/api/v2/version/5/"
		}
	  },
	  {
		"flavor_text": "DITTO rearranges its cell structure to\ntransform itself into other shapes.\nHowever, if it tries to transform itself\finto something by relying on its memory,\nthis POKéMON manages to get details\nwrong.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "ruby",
		  "url": "https://pokeapi.co/api/v2/version/7/"
		}
	  },
	  {
		"flavor_text": "DITTO rearranges its cell structure to\ntransform itself into other shapes.\nHowever, if it tries to transform itself\finto something by relying on its memory,\nthis POKéMON manages to get details\nwrong.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "sapphire",
		  "url": "https://pokeapi.co/api/v2/version/8/"
		}
	  },
	  {
		"flavor_text": "A DITTO rearranges its cell structure to\ntransform itself. However, if it tries to\nchange based on its memory, it will get\ndetails wrong.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "emerald",
		  "url": "https://pokeapi.co/api/v2/version/9/"
		}
	  },
	  {
		"flavor_text": "It can freely recombine its own cellular\nstructure to transform into other life-\nforms.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "firered",
		  "url": "https://pokeapi.co/api/v2/version/10/"
		}
	  },
	  {
		"flavor_text": "It has the ability to reconstitute\nits entire cellular structure to\ntransform into whatever it sees.",
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"version": {
		  "name": "platinum",
		  "url": "https://pokeapi.co/api/v2/version/14/"
		}
	  }
	  {
		"flavor_text": "Puede transformarse en cualquier cosa que vea,\npero, si intenta hacerlo de memoria, habrá\ndetalles que se le escapen.",
		"language": {
		  "name": "es",
		  "url": "https://pokeapi.co/api/v2/language/7/"
		},
		"version": {
		  "name": "ultra-moon",
		  "url": "https://pokeapi.co/api/v2/version/30/"
		}
	  }
	],
	"generation": {
	  "name": "generation-i",
	  "url": "https://pokeapi.co/api/v2/generation/1/"
	},
	"growth_rate": {
	  "name": "medium",
	  "url": "https://pokeapi.co/api/v2/growth-rate/2/"
	},
	"habitat": {
	  "name": "urban",
	  "url": "https://pokeapi.co/api/v2/pokemon-habitat/8/"
	},
	"has_gender_differences": false,
	"hatch_counter": 20,
	"id": 132,
	"is_baby": false,
	"is_legendary": false,
	"is_mythical": false,
	"name": "ditto",
	"names": [
	  {
		"language": {
		  "name": "it",
		  "url": "https://pokeapi.co/api/v2/language/8/"
		},
		"name": "Ditto"
	  },
	  {
		"language": {
		  "name": "en",
		  "url": "https://pokeapi.co/api/v2/language/9/"
		},
		"name": "Ditto"
	  }
	],
	"order": 156,
	"pal_park_encounters": [
	  {
		"area": {
		  "name": "field",
		  "url": "https://pokeapi.co/api/v2/pal-park-area/2/"
		},
		"base_score": 70,
		"rate": 20
	  }
	],
	"pokedex_numbers": [
	  {
		"entry_number": 132,
		"pokedex": {
		  "name": "national",
		  "url": "https://pokeapi.co/api/v2/pokedex/1/"
		}
	  },
	  {
		"entry_number": 132,
		"pokedex": {
		  "name": "kanto",
		  "url": "https://pokeapi.co/api/v2/pokedex/2/"
		}
	  }
	],
	"shape": {
	  "name": "ball",
	  "url": "https://pokeapi.co/api/v2/pokemon-shape/1/"
	},
	"varieties": [
	  {
		"is_default": true,
		"pokemon": {
		  "name": "ditto",
		  "url": "https://pokeapi.co/api/v2/pokemon/132/"
		}
	  }
	]
  }
`

func mockPokeAPIServer(t *testing.T, mockResp pokemonSpecies) {
	t.Helper()

	var apiCalled bool = false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pokeAPIPokemonSpeciesURL {
			t.Errorf("Expected to request %s, got: %s", pokeAPIPokemonSpeciesURL, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockOKResp))
		if err != nil {
			t.Errorf("Expect nil err, got %s", err)
		}
		apiCalled = true
	}))
	defer server.Close()

	if !apiCalled {
		t.Error("expect pokeAPI to be called, got no call.")
	}
}
