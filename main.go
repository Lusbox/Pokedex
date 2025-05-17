package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
			return fmt.Errorf("error reading body")
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
