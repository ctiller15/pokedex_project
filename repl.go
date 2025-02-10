package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands map[string]cliCommand

func listCommands() {
	for name, command := range supportedCommands {
		fmt.Printf("%s: %s\n", name, command.description)
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	listCommands()
	return nil
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
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		text := cleanInput(scanner.Text())

		if len(text) > 0 {
			command := text[0]

			registeredCommand, ok := supportedCommands[command]
			if !ok {
				fmt.Println("Unknown command")
				continue
			}

			registeredCommand.callback()
		}
	}
}
