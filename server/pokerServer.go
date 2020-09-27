package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	guuid "github.com/google/uuid"
	mux "github.com/gorilla/mux"
)

// A Card has a symbol and a number
type Card struct {
	Value  int
	Symbol string
}

// A Suit is a list of cards
type Suit [13]Card

// A Deck consists of a max of 52 Cards
type Deck []Card

var cardSet = [13]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

var nameMap = map[int]string{
	1:  "Ace",
	2:  "Two",
	3:  "Three",
	4:  "Four",
	5:  "Five",
	6:  "Six",
	7:  "Seven",
	8:  "Eight",
	9:  "Nine",
	10: "Ten",
	11: "Jack",
	12: "Queen",
	13: "King",
}

func main() {
	// dealAndShuffle()
	games := make(ActiveGames)
	connectToDealer(games)
}

// Player defines what cards are being held and their ids
type Player struct {
	playerID string
	cards    []Card
}

// Game defines an id, a group of players and the deck state
type Game struct {
	gameID  string
	players []*Player
	deck    Deck
}

// ActiveGames is a hash of currently active games
type ActiveGames map[string]*Game

func getNewUUIDString() string {
	newUUID := guuid.New()
	return newUUID.String()
}

func getNewPlayers() []*Player {
	randomNumPlayers := 1
	if randNum := rand.Intn(5); randNum > 0 {
		randomNumPlayers = randNum
	}

	players := []*Player{}

	for randomNumPlayers > 0 {
		playerID := getNewUUIDString()
		cards := []Card{}
		newPlayer := &Player{
			playerID: playerID,
			cards:    cards,
		}
		players = append(players, newPlayer)
		randomNumPlayers--
	}
	return players
}

func newGame(games ActiveGames, res http.ResponseWriter, req *http.Request) {
	gameID := getNewUUIDString()
	unshuffledDeck := buildDeck()
	deck := shuffleDeck(unshuffledDeck)
	var players = getNewPlayers()

	games[gameID] = &Game{
		gameID:  gameID,
		players: players,
		deck:    deck,
	}

	fmt.Println("GameID:", gameID)
	fmt.Println("Active Players:", games[gameID].players)
	outputStr := fmt.Sprintf("Game Started!\n Playing Game with an ID of:\n %v", gameID)
	io.WriteString(res, outputStr)
}

// ShuffState is a JSON blob getting returned from shuffling regarding deck
type ShuffState struct {
	Prev    Deck `json:"Prev"`
	Current Deck `json:"Current"`
}

// ShuffleJSONRes is the shape of the JSON blob getting returned from shuffling.
type ShuffleJSONRes struct {
	ActivePlayerIDs []string   `json:"ActivePlayerIDs"`
	ShuffledState   ShuffState `json:"ShuffledState"`
}

func shuffleGameDeck(games ActiveGames, res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux.Vars(req)
	gameID := vars["gameID"]
	fmt.Println("Passed GameID:", gameID)

	if gameID == "" {
		http.Error(res, "Missing Game ID", http.StatusBadRequest)
		return
	}
	// res.WriteHeader(http.StatusOK)

	unshuffledDeck := games[gameID].deck
	shuffledDeck := shuffleDeck(unshuffledDeck)
	games[gameID].deck = shuffledDeck

	activePlayerIDs := []string{}
	for i := 0; i < len(games[gameID].players); i++ {
		playerID := (*games[gameID].players[i]).playerID

		activePlayerIDs = append(activePlayerIDs, playerID)
	}

	jsonResponse := ShuffleJSONRes{
		ActivePlayerIDs: activePlayerIDs,
		ShuffledState: ShuffState{
			Prev:    unshuffledDeck,
			Current: shuffledDeck,
		},
	}
	byteArray, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(res, "Error converting to JSON", http.StatusInternalServerError)
	}
	io.WriteString(res, string(byteArray))
}

// DealJSONRes is the structure of the JSON for the deal call
type DealJSONRes struct {
	Hand []Card `json:"Hand"`
}

func dealCardsToPlayer(games ActiveGames, res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux.Vars(req)
	gameID := vars["gameID"]
	playerID := vars["playerID"]
	numCardsToDeal, err := strconv.Atoi(req.FormValue("numToDeal"))
	fmt.Println("Passed GameID/PlayerID:", gameID, playerID)

	if gameID == "" {
		http.Error(res, "Missing Game ID", http.StatusBadRequest)
		return
	}

	if playerID == "" {
		http.Error(res, "Missing Player ID", http.StatusBadRequest)
		return
	}

	gameDeck := games[gameID].deck
	playerHand := []Card{}

	for i := 0; i <= len(games[gameID].players); i++ {
		player := (*games[gameID].players[i])
		if pID := player.playerID; pID == playerID {
			log.Println("Player's Current Hand:", player.cards)

			dealtCards := dealCards(&gameDeck, numCardsToDeal)
			games[gameID].deck = gameDeck
			showCards(dealtCards)

			for _, dealtCard := range dealtCards {
				playerHand = append(playerHand, dealtCard.card)
			}
			games[gameID].players[i].cards = playerHand
			break
		}
	}

	fmt.Println("Player's New Hand", playerHand)
	jsonResponse := DealJSONRes{
		Hand: playerHand,
	}
	byteArray, err := json.Marshal(jsonResponse)
	if err != nil {
		http.Error(res, "Error converting to JSON", http.StatusInternalServerError)
	}
	io.WriteString(res, string(byteArray))
}

func connectToDealer(games ActiveGames) {
	fmt.Println("Connecting to dealer...")
	rand.Seed(time.Now().UTC().UnixNano())

	newGameHandler := func(res http.ResponseWriter, req *http.Request) {
		newGame(games, res, req)
	}
	shuffleGameDeckHandler := func(res http.ResponseWriter, req *http.Request) {
		shuffleGameDeck(games, res, req)
	}
	dealCardsToPlayerHandler := func(res http.ResponseWriter, req *http.Request) {
		dealCardsToPlayer(games, res, req)
	}

	r := mux.NewRouter()
	r.HandleFunc("/newgame/", newGameHandler).Methods("GET")
	r.HandleFunc("/shuffle/{gameID}/", shuffleGameDeckHandler).Methods("POST")
	r.HandleFunc("/deal/{gameID}/{playerID}/", dealCardsToPlayerHandler).Methods("POST")
	http.Handle("/", r)
	fmt.Println("Booting new server on localhost:9999")

	log.Fatal(http.ListenAndServe("localhost:9999", nil))

}

func dealAndShuffle() {

	// Create the unshuffled deck
	unshuffledDeck := buildDeck()
	fmt.Println("Unshuffled Deck:")
	fmt.Println(unshuffledDeck)

	// Initially shuffle the deck
	shuffledDeck := shuffleDeck(unshuffledDeck)
	fmt.Println("Inintially Shuffled Deck:")
	fmt.Println(shuffledDeck)

	var cardsToReshuffleAfter []Card

	// Deal the first card
	dealtCards := dealCards(&shuffledDeck, 1)

	addCardsToReshuffle(&cardsToReshuffleAfter, dealtCards)
	showCards(dealtCards)

	// Deal the second card
	dealtCards = dealCards(&shuffledDeck, 1)

	addCardsToReshuffle(&cardsToReshuffleAfter, dealtCards)
	showCards(dealtCards)

	// // Take the cards that were pulled out and reshuffle them back into the deck a new
	var deckToReshuffle Deck = append(shuffledDeck, cardsToReshuffleAfter...)
	reshuffledDeck := shuffleDeck(deckToReshuffle)

	// Show the reshuffled deck
	fmt.Println("Newly Reshuffled Deck")
	fmt.Println(reshuffledDeck)
}

func addCardsToReshuffle(cardsToReshuffleAfter *[]Card, dealtCards []DealtCard) {
	for _, dealt := range dealtCards {
		*cardsToReshuffleAfter = append((*cardsToReshuffleAfter)[:], dealt.card)
	}
}

func showCards(dealtCards []DealtCard) {
	for _, card := range dealtCards {
		fmt.Println(card.phrase, "Remaining Cards:", len(card.remainingDeck))
	}
}

func buildSuit(suitSymbol string) Suit {
	var cardSuit [13]Card
	for i := 0; len(cardSet)-1 >= i; i++ {
		cardValue := cardSet[i]
		cardItem := Card{Value: cardValue, Symbol: suitSymbol}
		cardSuit[i] = cardItem
	}
	return cardSuit
}

var (
	hearts   Suit = buildSuit("♡")
	diamonds Suit = buildSuit("♢")
	clubs    Suit = buildSuit("♣")
	spades   Suit = buildSuit("♠")
)

func buildDeck() (deck Deck) {
	// Build a fresh deck of ordered cards/suits
	var suits = [4]Suit{hearts, diamonds, clubs, spades}
	for _, suit := range suits {
		for _, card := range suit {
			deck = append(deck, card)
		}
	}
	return
}

func shuffleDeck(deck Deck) Deck {
	// Reset the seed time
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

// DealtCard For the whole dealt card
type DealtCard struct {
	phrase        string
	card          Card
	remainingDeck Deck
}

func dealCards(shuffledDeck *Deck, count int) []DealtCard {

	var dealtCards []DealtCard

	// Always default to at least 1 card
	if count == 0 {
		count = 1
	}

	for i := 0; count > i; i++ {

		cardToPullIndex := (len(*shuffledDeck) - 1) - i
		dealtCard := (*shuffledDeck)[:][cardToPullIndex]
		*shuffledDeck = (*shuffledDeck)[:][:cardToPullIndex]

		formattedCardNum := nameMap[dealtCard.Value]
		suitSymbol := dealtCard.Symbol

		cardStr := fmt.Sprintf("Dealt Card: %v of %s", formattedCardNum, suitSymbol)
		cardJustDealt := DealtCard{phrase: cardStr, card: dealtCard, remainingDeck: *shuffledDeck}
		dealtCards = append(dealtCards, cardJustDealt)
	}
	return dealtCards
}
