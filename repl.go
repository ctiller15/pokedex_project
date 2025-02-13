package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ctiller15/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

var supportedCommands map[string]cliCommand

var api = pokeapi.NewAPI()

func listCommands() {
	for name, command := range supportedCommands {
		fmt.Printf("%s: %s\n", name, command.description)
	}
}

func commandExit(args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	listCommands()
	return nil
}

func commandMap(args ...string) error {
	return api.FetchLocation(api.NextPageUrl)

}

func commandMapB(args ...string) error {
	if api.PrevPageUrl == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	return api.FetchLocation(api.PrevPageUrl)
}

func commandExplore(args ...string) error {
	if len(args) < 1 {
		return errors.New("explore must be called with a valid location")
	}
	location := args[0]
	fmt.Printf("Exploring %s...\n", location)

	return api.ExploreLocation(location)
}

func commandCatch(args ...string) error {
	if len(args) < 1 {
		return errors.New("catch must be called with a valid pokemon name")
	}

	pokemon := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	return api.CatchPokemon(pokemon)
}

func commandInspect(args ...string) error {
	if len(args) < 1 {
		return errors.New("inspect must be called with a valid pokemon name")
	}

	pokemon := args[0]
	return api.InspectPokemon(pokemon)
}

func commandPokedex(args ...string) error {
	return api.CheckPokedex()
}

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	spaceRegex := regexp.MustCompile(`\s+`)
	trimmedNoWhiteSpace := spaceRegex.ReplaceAllString(trimmed, " ")
	lowercased := strings.ToLower(trimmedNoWhiteSpace)
	if lowercased == "" {
		return []string{}
	}

	return strings.Split(lowercased, " ")
}

func startRepl() {
	// Break initialization cycle
	supportedCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "displays 20 location areas in the pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays previous 20 location areas in the pokemon world.",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "displays all of the pokemon in a specific location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "tries to catch a pokemon and adds it to your pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "shows info about a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "lists all of the pokemon that the user has caught",
			callback:    commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		text := cleanInput(scanner.Text())

		if len(text) > 0 {
			command := text[0]
			arguments := text[1:]

			registeredCommand, ok := supportedCommands[command]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			err := registeredCommand.callback(arguments...)
			if err != nil {
				fmt.Printf("error occurred: %v\n", err)
			}
		}
	}
}
