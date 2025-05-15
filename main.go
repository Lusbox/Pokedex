package main

import (
	"fmt"
	"bufio"
	"os"
	)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var commands map[string]cliCommand
		

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays the help message",
			callback: commandHelp,
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
			err := c.callback()
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, value := range commands {
		msg := fmt.Sprintf("%s: %s", value.name, value.description)
		fmt.Println(msg)
	}
	return nil
}