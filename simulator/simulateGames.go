package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func simulateGames() {
	fmt.Print("Simulate some poker? Y or N\n ")
	userInput := bufio.NewScanner(os.Stdin)
	for userInput.Scan() {
		if count, _ := strconv.Atoi(userInput.Text()); userInput.Text() == "Y" || userInput.Text() == "y" {
			fmt.Printf("\n ### You said %s! Let's go! ### \n", userInput.Text())

			fmt.Print("\nHow many games to simulate?\n \n ")
		} else if count <= 10 && count != 0 {
			fmt.Printf("\n ### You asked me to simulate %v games. Let's... GO? ### \n", count)
			for i := 1; i <= count; i++ {
				fmt.Println("Starting simulation:", i)
				go startGameSim()
			}
		} else if count > 10 {
			fmt.Printf("\n ### You asked me to simulate %v games dude... No way. Pick another number, or quit with Exit. ### \n", userInput.Text())
		} else if userInput.Text() == "Exit" || userInput.Text() == "exit" || userInput.Text() == "N" || userInput.Text() == "n" {
			os.Exit(0)
		} else {
			fmt.Printf("\n ### You asked me to %s ### \n", userInput.Text())
			fmt.Printf("I don't know how to %s\n", userInput.Text())
		}
	}
}

func asyncShuffler(gameID string) <-chan []string {
	r := make(chan []string)
	go func() {
		defer close(r)
		playerIDs := shuffleCardsForGame(gameID)
		r <- playerIDs
	}()
	return r
}

func startGameSim() {
	gameID := getNewGameID()
	if gameID != "" {
		playerIDs := <-asyncShuffler(gameID)

		fmt.Println("ActivePlayerIDs:", playerIDs)
		dealCardsToPlayers(gameID, playerIDs)
	}
}

// A Card has a symbol and a number
type Card struct {
	Value  int
	Symbol string
}

// ShuffState is a JSON blob getting returned from shuffling regarding deck
type ShuffState struct {
	Prev    []Card `json:"Prev"`
	Current []Card `json:"Current"`
}

// ShuffleJSONRes is the shape of the JSON blob getting returned from shuffling.
type ShuffleJSONRes struct {
	ActivePlayerIDs []string   `json:"ActivePlayerIDs"`
	ShuffledState   ShuffState `json:"ShuffledState"`
}

func shuffleCardsForGame(gameID string) []string {
	shuffleCardsURL := fmt.Sprintf("http://localhost:9999/shuffle/%v/", gameID)
	log.Println("Attempting to shuffle @", shuffleCardsURL)
	r, err := http.PostForm(shuffleCardsURL, url.Values{"gameID": {gameID}})
	if err != nil {
		log.Panicln("Oh noez! Shuffling Error:", err)
		return []string{}
	}
	defer r.Body.Close()

	var data ShuffleJSONRes

	rData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if jsonErr := json.Unmarshal(rData, &data); jsonErr != nil {
		log.Panic(err)
	}

	// log.Println("Previous Shuffled State:", data.ShuffledState.Prev)
	// log.Println("Current Shuffled State:", data.ShuffledState.Current)
	// log.Println("Active Player IDs:", data.ActivePlayerIDs)

	return data.ActivePlayerIDs
}

// DealJSONRes is the structure for the JSON response from DealURL
type DealJSONRes struct {
	Hand []Card `json:"Hand"`
}

func dealCardsToPlayers(gameID string, playerIDs []string) {
	for _, pID := range playerIDs {
		dealCardURL := fmt.Sprintf("http://localhost:9999/deal/%v/%v/", gameID, pID)
		log.Println("Attempting to deal cards @", dealCardURL)
		numToDeal := ""
		r, err := http.PostForm(dealCardURL, url.Values{"gameID": {gameID}, "playerID": {pID}, "numToDeal": {numToDeal}})
		if err != nil {
			log.Panicln("Oh noez! Dealing Cards Error:", err)
			return
		}
		defer r.Body.Close()

		var data DealJSONRes

		rData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		if jsonErr := json.Unmarshal(rData, &data); jsonErr != nil {
			log.Panic(err)
		}
		log.Println("Player's Current Hand:", data.Hand)
	}
}

func getNewGameID() string {
	r, err := http.Get("http://localhost:9999/newgame/")
	if err != nil {
		log.Panicln("Oh noezzz! Failed to create new game error:", err)
		return ""
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln("Failed to read body:", err)
		return ""
	}
	bodyString := string(body)
	result := strings.Split(bodyString, "\n")[:]
	// Last item in the response body is the GameID
	gameID := result[len(result)-1]
	log.Print("New GameID:", gameID)
	return strings.Trim(gameID, " ")
}

func main() {
	simulateGames()
}
