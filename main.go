package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Lusbox/Pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}

	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if input.Scan() {
			word := input.Text()
			c, ok := commands[word]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}
			err := c.callback(cfg)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, value := range commands {
		msg := fmt.Sprintf("%s: %s", value.name, value.description)
		fmt.Println(msg)
	}
	return nil
}

func commandMap(c *config) error {
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

func commandMapb(c *config) error {
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
