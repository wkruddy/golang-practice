package main

import (
	"fmt"
	"math/rand"
	"time"
)

// A Card has a symbol and a number
type Card struct {
	value  int
	symbol string
}

// A Suit is a list of cards
type Suit []Card

var cardSet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

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

func buildSuit(suitSymbol string) Suit {
	var cardSuit []Card
	for i := 0; len(cardSet)-1 >= i; i++ {
		cardValue := cardSet[i]
		cardItem := Card{value: cardValue, symbol: suitSymbol}
		cardSuit = append(cardSuit, cardItem)
	}
	return cardSuit
}

var (
	hearts   Suit = buildSuit("♡")
	diamonds Suit = buildSuit("♢")
	clubs    Suit = buildSuit("♣")
	spades   Suit = buildSuit("♠")
)

func buildDeck() []Card {
	// Build a fresn deck of ordered cards/suits
	deck := hearts
	deck = append(deck, diamonds...)
	deck = append(deck, clubs...)
	deck = append(deck, spades...)
	return deck
}

func shuffleDeck(deck []Card) []Card {
	// Instead of doing nested loops, I opted for something different to verify things were randomly pulled/shuffled into a new deck:
	// For every card, get a random index of the 52, pull that card. If the card exists in the "already shuffled cards" map
	// Recursively try again for another random card until it finds one that hasnt been used, then it adds it to the shuffled deck
	var (
		cardsRemaining  int = len(deck)
		shuffledDeck    []Card
		addShuffledCard func()
	)
	usedCards := make(map[string]Card)

	addShuffledCard = func() {
		cardIndexToPull := getRandomizedCardNum(len(deck))
		cardToAdd := deck[cardIndexToPull]

		key := fmt.Sprintf("%v-%s", cardToAdd.value, cardToAdd.symbol)

		_, cardBeenShuffled := usedCards[key]

		if cardBeenShuffled && len(usedCards) < len(deck) {
			addShuffledCard()
		} else {
			usedCards[key] = cardToAdd

			shuffledDeck = append(shuffledDeck, cardToAdd)
			cardsRemaining--
			return
		}
	}

	for cardsRemaining > 0 {
		addShuffledCard()
	}
	return shuffledDeck
}

func getRandomizedCardNum(maxVal int) (num int) {
	num = 0 + rand.Intn(maxVal)
	return
}

func dealCard(shuffledDeck []Card) (string, Card, []Card) {
	card, remainingDeck := shuffledDeck[len(shuffledDeck)-1], shuffledDeck[:len(shuffledDeck)-1]

	formattedCardNum := nameMap[card.value]
	suitSymbol := card.symbol

	cardStr := fmt.Sprintf("Dealt Card: %v of %s", formattedCardNum, suitSymbol)

	return cardStr, card, remainingDeck
}

func main() {
	// Initialize the program time seed
	rand.Seed(time.Now().UTC().UnixNano())

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
	cardStr, dealtCard, remainingDeck := dealCard(shuffledDeck)
	cardsToReshuffleAfter = append(cardsToReshuffleAfter, dealtCard)
	fmt.Println(cardStr, "Remaining Cards:", len(remainingDeck))

	// Deal the second card
	cardStr, dealtCard, remainingDeck = dealCard(remainingDeck)
	cardsToReshuffleAfter = append(cardsToReshuffleAfter, dealtCard)
	fmt.Println(cardStr, "Remaining Cards:", len(remainingDeck))

	// Take the cards that were pulled out and reshuffle them back into the deck a new
	var deckToReshuffle []Card = append(remainingDeck, cardsToReshuffleAfter...)
	reshuffledDeck := shuffleDeck(deckToReshuffle)

	// Show the reshuffled deck
	fmt.Println("Newly Reshuffled Deck")
	fmt.Println(reshuffledDeck)
}
