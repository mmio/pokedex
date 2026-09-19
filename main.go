package main

import (
	"bufio"
	"fmt"
	"os"
)

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(config *Config) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:`)

	for _, command := range config.commands {
		_, err := fmt.Printf("%v: %v\n", command.name, command.description)
		if err != nil {
			return err
		}
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	commands map[string]cliCommand
}

func main() {
	config := Config{
		commands: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exits the Pokedex",
				callback:    commandExit,
			}, "help": {
				name:        "help",
				description: "Shows help for the Pokedex",
				callback:    commandHelp,
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	var input string
	var safeInput []string
	var commandName string
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			fmt.Printf("Couldn't scan input: %v", scanner.Err())
			break
		}

		input = scanner.Text()
		safeInput = cleanInput(input)

		if len(safeInput) == 0 {
			continue
		}

		commandName = safeInput[0]
		command, ok := config.commands[commandName]
		if !ok {
			fmt.Printf("Unknown command '%v'\n", safeInput[0])
			continue
		}

		if err := command.callback(&config); err != nil {
			fmt.Println("Couldn't execute command")
			continue
		}
	}
}
