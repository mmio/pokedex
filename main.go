package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/mmio/pokedex/pokeapi"
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

func commandMap(config *Config) error {
	response, err := config.pokeAPI.CallMap()
	if err != nil {
		return err
	}

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(config *Config) error {
	response, err := config.pokeAPI.CallMapBack()
	if err != nil {
		return err
	}

	for _, location := range response.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(config *Config) error {
	if len(config.arguments) == 0 {
		return errors.New("Explore needs the name of the location as argument")
	}

	locationName := config.arguments[0]

	response, err := config.pokeAPI.CallExplore(locationName)
	if err != nil {
		return err
	}

	fmt.Println("Exploring pastoria-city-area...")
	fmt.Println("Found Pokemon:")
	for _, encounter := range response.PokemonEncounters {
		fmt.Println("-", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *Config) error {
	if len(config.arguments) == 0 {
		return errors.New("Catch needs the name of the pokemon as argument")
	}

	pokemonName := config.arguments[0]

	response, err := config.pokeAPI.CallPokemonInfo(pokemonName)
	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + pokemonName + "...")

	if CanICatchIt(response.BaseExperience) {
		fmt.Println(pokemonName + " was caught!")
		config.pokemons = append(config.pokemons, pokemonName)
		return nil
	}

	fmt.Println(pokemonName + " escaped!")
	return nil
}

func commandInspect(config *Config) error {
	if len(config.arguments) == 0 {
		return errors.New("Inspect needs the name of the pokemon as argument")
	}

	pokemonName := config.arguments[0]

	if !slices.Contains(config.pokemons, pokemonName) {
		fmt.Println("You haven't yet caught this pokemon")
		return nil
	}

	response, err := config.pokeAPI.CallPokemonInfo(pokemonName)
	if err != nil {
		return err
	}

	fmt.Println(response)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	pokeAPI   pokeapi.PokeAPIState
	commands  map[string]cliCommand
	command   string
	arguments []string
	pokemons  []string
}

func main() {
	pas, err := pokeapi.NewPokeAPI()
	if err != nil {
		fmt.Println("Error initializing poke api", err)
		os.Exit(0)
	}

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
			}, "map": {
				name:        "map",
				description: "Shows next location areas",
				callback:    commandMap,
			}, "mapb": {
				name:        "mapb",
				description: "Shows previous location areas",
				callback:    commandMapBack,
			}, "explore": {
				name:        "explore",
				description: "Shows pokemons at location",
				callback:    commandExplore,
			}, "catch": {
				name:        "cache",
				description: "Catch a pokemon",
				callback:    commandCatch,
			}, "inspect": {
				name:        "inspect",
				description: "Inspect a caught pokemon",
				callback:    commandInspect,
			},
		},
		pokeAPI:   pas,
		command:   "",
		arguments: []string{},
		pokemons:  []string{},
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

		config.command = safeInput[0]
		config.arguments = safeInput[1:]
		commandName = safeInput[0]
		command, ok := config.commands[commandName]
		if !ok {
			fmt.Printf("Unknown command '%v'\n", safeInput[0])
			continue
		}

		if err := command.callback(&config); err != nil {
			fmt.Println("Couldn't execute command", err)
			continue
		}
	}
}
