package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Lusbox/Pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

var commands map[string]cliCommand

var pokedex map[string]pokemonDetails

type maplocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type areaPokemon struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type pokemonDetails struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy any    `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices            []any  `json:"game_indices"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        int `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      any    `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        any    `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault any `json:"front_default"`
				FrontFemale  any `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      any `json:"back_default"`
					BackGray         any `json:"back_gray"`
					BackTransparent  any `json:"back_transparent"`
					FrontDefault     any `json:"front_default"`
					FrontGray        any `json:"front_gray"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      any `json:"back_default"`
					BackGray         any `json:"back_gray"`
					BackTransparent  any `json:"back_transparent"`
					FrontDefault     any `json:"front_default"`
					FrontGray        any `json:"front_gray"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           any `json:"back_default"`
					BackShiny             any `json:"back_shiny"`
					BackShinyTransparent  any `json:"back_shiny_transparent"`
					BackTransparent       any `json:"back_transparent"`
					FrontDefault          any `json:"front_default"`
					FrontShiny            any `json:"front_shiny"`
					FrontShinyTransparent any `json:"front_shiny_transparent"`
					FrontTransparent      any `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      any `json:"back_default"`
					BackShiny        any `json:"back_shiny"`
					FrontDefault     any `json:"front_default"`
					FrontShiny       any `json:"front_shiny"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      any `json:"back_default"`
					BackShiny        any `json:"back_shiny"`
					FrontDefault     any `json:"front_default"`
					FrontShiny       any `json:"front_shiny"`
					FrontTransparent any `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  any `json:"back_default"`
					BackShiny    any `json:"back_shiny"`
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  any `json:"back_default"`
					BackShiny    any `json:"back_shiny"`
					FrontDefault any `json:"front_default"`
					FrontShiny   any `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      any `json:"back_default"`
					BackFemale       any `json:"back_female"`
					BackShiny        any `json:"back_shiny"`
					BackShinyFemale  any `json:"back_shiny_female"`
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      any `json:"back_default"`
					BackFemale       any `json:"back_female"`
					BackShiny        any `json:"back_shiny"`
					BackShinyFemale  any `json:"back_shiny_female"`
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      any `json:"back_default"`
					BackFemale       any `json:"back_female"`
					BackShiny        any `json:"back_shiny"`
					BackShinyFemale  any `json:"back_shiny_female"`
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      any `json:"back_default"`
						BackFemale       any `json:"back_female"`
						BackShiny        any `json:"back_shiny"`
						BackShinyFemale  any `json:"back_shiny_female"`
						FrontDefault     any `json:"front_default"`
						FrontFemale      any `json:"front_female"`
						FrontShiny       any `json:"front_shiny"`
						FrontShinyFemale any `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      any `json:"back_default"`
					BackFemale       any `json:"back_female"`
					BackShiny        any `json:"back_shiny"`
					BackShinyFemale  any `json:"back_shiny_female"`
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault any `json:"front_default"`
					FrontFemale  any `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     any `json:"front_default"`
					FrontFemale      any `json:"front_female"`
					FrontShiny       any `json:"front_shiny"`
					FrontShinyFemale any `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault any `json:"front_default"`
					FrontFemale  any `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
}

func main() {
	cfg := &config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    pokecache.NewCache(5 * time.Minute),
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays the help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Previous locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Add area name to show Pokemon found",
			callback:    commmandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch Pokemon",
			callback:    commandCatch,
		},
	}

	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if input.Scan() {
			words := strings.Fields(input.Text())
			if len(words) == 0 {
				continue
			}

			commandName := words[0]
			args := words[1:]

			c, ok := commands[commandName]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			err := c.callback(cfg, args...)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func commandExit(c *config, name ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, name ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, value := range commands {
		msg := fmt.Sprintf("%s: %s", value.name, value.description)
		fmt.Println(msg)
	}
	return nil
}

func commandMap(c *config, name ...string) error {
	cachedData, ok := c.cache.Get(c.Next)

	var body []byte
	if ok {
		body = cachedData
	} else {
		res, err := http.Get(c.Next)
		if err != nil {
			return fmt.Errorf("error getting response: %v", err)
		}

		defer res.Body.Close()

		var err2 error
		body, err2 = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return fmt.Errorf("error with statuscode: %v", res.StatusCode)
		}
		if err2 != nil {
			return fmt.Errorf("error reading body")
		}
		c.cache.Add(c.Next, body)
	}

	var locations maplocations
	err := json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("data not json format: %v", err)
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, item := range locations.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func commandMapb(c *config, name ...string) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	cachedData, ok := c.cache.Get(c.Previous)

	var body []byte
	if ok {
		body = cachedData
	} else {
		res, err := http.Get(c.Previous)
		if err != nil {
			return fmt.Errorf("error getting response: %v", err)
		}

		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			return fmt.Errorf("error with statuscode: %v", res.StatusCode)
		}
		if err != nil {
			return fmt.Errorf("error reading body")
		}

		c.cache.Add(c.Previous, body)
	}

	var locations maplocations
	err := json.Unmarshal(body, &locations)
	if err != nil {
		return fmt.Errorf("data not json format: %v", err)
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, item := range locations.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func commmandExplore(c *config, name ...string) error {
	if len(name) == 0 || name[0] == "" {
		return fmt.Errorf("please provide a location area name")
	}

	url := "https://pokeapi.co/api/v2/location-area/" + name[0]
	cachedData, ok := c.cache.Get(url)

	var body []byte
	if ok {
		body = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error getting response: %v", err)
		}

		defer res.Body.Close()

		var err2 error
		body, err2 = io.ReadAll(res.Body)
		if err2 != nil {
			return fmt.Errorf("error reading body: %v", err2)
		}

		c.cache.Add(url, body)
	}

	var pokemon areaPokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		if strings.Contains(string(body), "Not Found") {
			return fmt.Errorf("location area '%s' not found", name[0])
		}
		return fmt.Errorf("error parsing response: %v", err)
	}

	fmt.Printf("Exploring %s...\n", name[0])
	fmt.Println("Found Pokemon:")

	for _, item := range pokemon.PokemonEncounters {
		fmt.Printf(" - %s\n", item.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, name ...string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + name[0]

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error getting response:  %v", err)
	}

	defer res.Body.Close()

	var err2 error
	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return fmt.Errorf("error reading body: %s", err2)
	}

	var pokemon pokemonDetails
	if err := json.Unmarshal(body, &pokemon); err != nil {
		if strings.Contains(string(body), "Not Found") {
			return fmt.Errorf("pokemon '%s' not found", name[0])
		}
		return fmt.Errorf("error parsing response: %v", err)
	}

	pokedex = make(map[string]pokemonDetails)

	fmt.Printf("Throwing a Pokeball at %s...\n", name[0])

	attempt := rand.Intn(pokemon.BaseExperience)
	if attempt < 20  {
		fmt.Printf("%s was caught!\n", name[0])
		fmt.Printf("Adding %s to Pokedex\n", name[0])
		pokedex[name[0]] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name[0])
		return nil
	}
	return nil

}
