package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/ctiller15/pokedexcli/internal/pokecache"
)

type encounterMethod struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type version struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type versionDetail struct {
	Rate    int     `json:"rate"`
	Version version `json:"version"`
}

type encounterMethodRate struct {
	EncounterMethod encounterMethod `json:"encounter_method"`
	VersionDetails  []versionDetail `json:"version_details"`
}

type language struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type alternateName struct {
	Name     string   `json:"name"`
	Language language `json:"language"`
}

type pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type encounterDetail struct {
	Chance int `json:"chance"`
	// ConditionValues
	MaxLevel int             `json:"max_level"`
	MinLevel int             `json:"min_level"`
	Method   encounterMethod `json:"method"`
}

type pokemonVersionDetail struct {
	EncounterDetails []encounterDetail `json:"encounter_details"`
	MaxChance        int               `json:"max_chance"`
	Version          version           `json:"version"`
}

type pokemonEncounter struct {
	Pokemon        pokemon                `json:"pokemon"`
	VersionDetails []pokemonVersionDetail `json:"version_details"`
}

type exploreResponse struct {
	EncounterMethodRates []encounterMethodRate `json:"encounter_method_rates"`
	GameIndex            int                   `json:"game_index"`
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	AltNames             []alternateName       `json:"names"`
	PokemonEncounters    []pokemonEncounter    `json:"pokemon_encounters"`
}

type mapResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

type location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type StatDetail struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonStat struct {
	BaseStat int        `json:"base_stat"`
	Stat     StatDetail `json:"stat"`
}

type PokemonTypeInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonType struct {
	Slot string            `json:"slot"`
	Type []PokemonTypeInfo `json:"type"`
}

type pokemonResponseType struct {
	Slot int             `json:"slot"`
	Type PokemonTypeInfo `json:"type"`
}

type pokemonResponse struct {
	Name           string                `json:"name"`
	BaseExperience int                   `json:"base_experience"`
	Stats          []PokemonStat         `json:"stats"`
	Height         int                   `json:"height"`
	Weight         int                   `json:"weight"`
	Types          []pokemonResponseType `json:"types"`
}

type API struct {
	PrevPageUrl string
	NextPageUrl string
	c           pokecache.Cache
}

// type PokemonStats struct {
// 	hp             int
// 	attack         int
// 	defense        int
// 	specialAttack  int
// 	specialDefense int
// 	speed          int
// }

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  map[string]int
	Types  []string
}

var pokedex map[string]Pokemon = make(map[string]Pokemon)

func NewAPI() *API {
	api := API{
		PrevPageUrl: "",
		NextPageUrl: "https://pokeapi.co/api/v2/location-area",
		c:           *pokecache.NewCache(time.Hour * 6),
	}

	return &api
}

func (a *API) fetchWithCache(url string) ([]byte, error) {
	var resBody []byte
	data, found := a.c.Get(url)

	if found {
		resBody = data
	} else {
		res, err := http.Get(url)

		if err != nil {
			return nil, err
		}
		resBody, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("response failed with status code: %d and \nbody: %s", res.StatusCode, resBody)
		}
		if err != nil {
			return nil, err
		}

		a.c.Add(url, resBody)
	}

	return resBody, nil
}

func (a *API) addToPokedex(pokemon pokemonResponse) {
	pokemonStats := make(map[string]int)
	types := make([]string, 0)

	for _, stat := range pokemon.Stats {
		pokemonStats[stat.Stat.Name] = stat.BaseStat
	}

	for _, pokeType := range pokemon.Types {
		types = append(types, pokeType.Type.Name)
	}

	pokedex[pokemon.Name] = Pokemon{
		Name:   pokemon.Name,
		Height: pokemon.Height,
		Weight: pokemon.Weight,
		Stats:  pokemonStats,
		Types:  types,
	}
}

func (a *API) FetchLocation(url string) error {
	body, err := a.fetchWithCache(url)
	if err != nil {
		return err
	}

	var responseStruct mapResponse
	err = json.Unmarshal(body, &responseStruct)

	if err != nil {
		return err
	}

	for _, location := range responseStruct.Results {
		fmt.Println(location.Name)
	}

	a.NextPageUrl = responseStruct.Next
	a.PrevPageUrl = responseStruct.Previous
	return nil
}

func (a *API) ExploreLocation(location string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + location

	body, err := a.fetchWithCache(url)

	if err != nil {
		return err
	}

	var exploreResponse exploreResponse
	err = json.Unmarshal(body, &exploreResponse)

	if err != nil {
		return err
	}

	if len(exploreResponse.PokemonEncounters) > 0 {
		fmt.Println("Found Pokemon:")
	}

	for _, pokemonEncounter := range exploreResponse.PokemonEncounters {
		fmt.Println("- " + pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func (a *API) CatchPokemon(pokemon string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	body, err := a.fetchWithCache(url)

	if err != nil {
		return err
	}

	var pokemonResponse pokemonResponse
	err = json.Unmarshal(body, &pokemonResponse)

	if err != nil {
		return err
	}

	baseExp := pokemonResponse.BaseExperience

	catchChance := 36.0 / float64(baseExp)

	caughtRand := rand.Float64()

	if catchChance > caughtRand {
		fmt.Printf("Caught %s!\n", pokemon)
		_, ok := pokedex[pokemon]
		if !ok {
			fmt.Println("Adding to Pokedex...")
			a.addToPokedex(pokemonResponse)
			fmt.Printf("Complete! You may now inspect %s with the inspect command.\n", pokemon)
		}

	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return nil
}

func (a *API) InspectPokemon(pokemon string) error {
	pokemonData, ok := pokedex[pokemon]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemonData.Name)
	fmt.Printf("Height: %d\n", pokemonData.Height)
	fmt.Printf("Weight: %d\n", pokemonData.Weight)

	// Yeah I know this is messy. I didn't want to do a slice and
	// I wanted a deterministic outcome.
	fmt.Println("Stats:")
	fmt.Printf("  -hp: %d\n", pokemonData.Stats["hp"])
	fmt.Printf("  -attack: %d\n", pokemonData.Stats["attack"])
	fmt.Printf("  -defense: %d\n", pokemonData.Stats["defense"])
	fmt.Printf("  -special-attack: %d\n", pokemonData.Stats["special-attack"])
	fmt.Printf("  -special-defense: %d\n", pokemonData.Stats["special-defense"])
	fmt.Printf("  -speed: %d\n", pokemonData.Stats["speed"])

	fmt.Println("Types:")
	for _, pokeType := range pokemonData.Types {
		fmt.Printf("  - %s\n", pokeType)
	}

	return nil
}

func (a *API) CheckPokedex() error {
	if len(pokedex) == 0 {
		return errors.New("you haven't caught any pokemon yet")
	}

	fmt.Println("Your Pokedex:")
	for key := range pokedex {
		fmt.Println(key)
	}

	return nil
}
