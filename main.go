package main
import ("fmt"; "strings"; "bufio"; "os"; "net/http"; "encoding/json"; "math/rand"; 
"errors"; "time"; "github.com/gabeportillo51/PokeDex/internal/pokecache"; "io"; "bytes")

func cleanInput(text string) []string{
	// clean up user input for processing
	text = strings.TrimSpace(text) 
	text = strings.ToLower(text) 
	clean_text := strings.Split(text, " ")
	return clean_text
}

func commandExit(loc *config, param string) error {
	// close the program
	fmt.Println()
	fmt.Println("Closing the Pokedex... Goodbye!")
	fmt.Println()
	os.Exit(0)
	return nil
}

func commandHelp(loc *config, param string) error {
	// provide info on how to use the program
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for key, value := range commandRegistry {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMap(loc *config, param string) error {
	// navigate pages in forward order
	if loc.Next == nil {
		return errors.New("That page is nil")
	}
	loc.Current = loc.Next
	var reader *bytes.Reader
	var byteData []byte
	byteData, ok := poke_cache.Get(*loc.Next)
	if !ok {
    	response, err := http.Get(*loc.Next)
    	if err != nil {
        	return err
    	}
    	defer response.Body.Close()
    	byteData, err = io.ReadAll(response.Body)
    	if err != nil {
        	return err
    	}
    	poke_cache.Add(*loc.Next, byteData)
	}
	reader = bytes.NewReader(byteData)
	loc.Page_number += 1
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&area); err != nil {
    	return err
	}
	fmt.Printf("\nYou are on page number: %d\n\n", loc.Page_number)
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	fmt.Println()
	loc.Next = area.Next
	loc.Previous = area.Previous
	return nil
}

func commandMapB(loc *config, param string) error {
	// navigate pages in reverse order
	if loc.Previous == nil {
		return errors.New("That page is nil")
	}
	loc.Current = loc.Previous
	var reader *bytes.Reader
	var byteData []byte
	byteData, ok := poke_cache.Get(*loc.Previous)
	if !ok {
    	response, err := http.Get(*loc.Previous)
    	if err != nil {
        	return err
    	}
    	defer response.Body.Close()
    	byteData, err = io.ReadAll(response.Body)
    	if err != nil {
        	return err
    	}
    	poke_cache.Add(*loc.Previous, byteData)
	}
	reader = bytes.NewReader(byteData)

	loc.Page_number -= 1
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&area); err != nil {
    	return err
	}
	fmt.Printf("\nYou are on page number: %d\n\n", loc.Page_number)
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	fmt.Println()
	loc.Next = area.Next
	loc.Previous = area.Previous
	return nil
}

func commandExplore(loc *config, area_name string) error {
	if area_name == "" {
		return errors.New("You didn't provide an area to explore. Try Again. Type: 'explore <area_name>'")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", area_name)
	var reader *bytes.Reader
	var byteData []byte
	byteData, ok := poke_cache.Get(url)
	if ok {
		fmt.Println("")
	}
	if !ok {
		fmt.Println("")
    	response, err := http.Get(url)
    	if err != nil {
        	return err
    	}
    	defer response.Body.Close()
    	byteData, err = io.ReadAll(response.Body)
    	if err != nil {
        	return err
    	}
    	poke_cache.Add(url, byteData)
	}
	fmt.Println("")
	fmt.Printf("Exploring %s...\n", area_name)
	reader = bytes.NewReader(byteData)
	decoder := json.NewDecoder(reader)
	var area_pokemon areaPokemon
	if err := decoder.Decode(&area_pokemon); err != nil {
    	return err
	}
	fmt.Println("")
	for _, p := range area_pokemon.PokemonEncounters {
		fmt.Println(p.Pokemon.Name)
	}
	fmt.Println()
	return nil
}

func commandCatch(loc *config, pokemon_name string) error {
	if pokemon_name == "" {
		return errors.New("You didn't specify which Pokemon you want to catch. Try Again.")
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon_name)
	var reader *bytes.Reader
	var byteData []byte
    response, err := http.Get(url)
    if err != nil {
        return err
    }
    defer response.Body.Close()
    byteData, err = io.ReadAll(response.Body)
    if err != nil {
        return err
    }
	fmt.Println("")
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon_name)
	reader = bytes.NewReader(byteData)
	decoder := json.NewDecoder(reader)
	var pokemon Pokemon
	if err := decoder.Decode(&pokemon); err != nil {
    	return err
	}
	catchProbability := 100 - (pokemon.BaseExperience / 3)
	randomNumber := rand.Intn(100)
	if randomNumber < catchProbability {
		fmt.Printf("%s was caught!\n", pokemon_name)
		UserPokeDex[pokemon_name] = pokemon
		fmt.Println("You may now inspect it with the 'inspect' command.")
	} else {
		fmt.Printf("%s escaped!", pokemon_name)
	}
	fmt.Println()
	return nil
}

func commandInspect(loc *config, pokemon_name string) error {
	if pokemon_name == "" {
		return errors.New("You didn't specify a pokemon name. Try Again.")
	}
	pokemon_data, ok := UserPokeDex[pokemon_name]
	if !ok {
		return errors.New("You don't have a Pokemon by that name.")
	}
	fmt.Println()
	fmt.Printf("Name: %s\n", pokemon_data.Name)
	fmt.Printf("Height: %d\n", pokemon_data.Height)
	fmt.Printf("Weight: %d\n", pokemon_data.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon_data.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon_data.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}
	fmt.Println()
	return nil
}

func commandPokedex(loc *config, param string) error {
	if len(UserPokeDex) == 0 {
		return errors.New("Your PokeDex is empty!")
	}
	fmt.Println()
	fmt.Println("Your Pokedex:")
	for _, pokemon := range UserPokeDex {
		fmt.Printf("  -%s\n", pokemon.Name)
	}
	fmt.Println()
	return nil
}

type cliCommand struct {
	// This is a struct the is used to describe a command
	name string     
	description string     
	callback func(*config, string) error    
}

type config struct {
	// this struct stores page info to navigate the map
	Next *string
	Previous *string
	Current *string
	Page_number int
}

type locationArea struct {
	// holds info about the current page of map
    Count    int    `json:"count"`
    Next     *string `json:"next"`     
    Previous *string `json:"previous"` 
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

type areaPokemon struct {
	// holds list of all pokemon in a speicifc location area
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct{
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`

}

var commandRegistry map[string]cliCommand
var poke_cache *pokecache.Cache
var area locationArea
var UserPokeDex map[string]Pokemon

func main(){
	UserPokeDex = map[string]Pokemon{}
	poke_cache = pokecache.NewCache(5 * time.Second)
	commandRegistry = map[string]cliCommand {
		"help": {
			name: "help",
			description: "How to use the Pokedex.",
			callback: commandHelp,
		},
		"exit": {
			name: "exit",
			description: "Exit the PokeDex.",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "Display the next 20 areas.",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Display the previous 20 areas.",
			callback: commandMapB,
		},
		"explore": {
			name: "explore",
			description: "Show all Pokemon in the specified area that are available to catch.",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Attempt to catch specified Pokemon using a pokeball.",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Examine specified Pokemon to learn more about it.",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "View all of the Pokemon in your PokeDex.",
			callback: commandPokedex,
		},
	}

	beginning := "https://pokeapi.co/api/v2/location-area/"
	location := config{
	Next: &beginning,
	Previous: nil,
	Current: nil,
	Page_number: 0,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {                           
		fmt.Print("PokeDex > ")
		scanner.Scan()                           
		words := cleanInput(scanner.Text())   
		if len(words) > 2 {
			fmt.Println("Too many arguments provided")
		} else {   
			var param string
			command := words[0]
			if len(words) == 2 {
				param = words[1]
			} else {
				param = ""
			}        
			com, ok := commandRegistry[command]    
			if ok {                                   
				err := com.callback(&location, param)       
				if err != nil {                        
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				fmt.Println("Unknown command")                
			}
		}
	}
	
}