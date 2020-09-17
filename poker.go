package main

import (
	"fmt"
	"math/rand"
	"time"
)

type card struct {
	value  int
	symbol string
}

type suit []card

var cardSet = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

func buildSuit(suitSymbol string) suit {
	var cardSuit []card
	for i := 0; len(cardSet)-1 >= i; i++ {
		cardValue := cardSet[i]
		cardItem := card{value: cardValue, symbol: suitSymbol}
		cardSuit = append(cardSuit, cardItem)
	}
	return cardSuit
}

var (
	hearts   suit = buildSuit("♡")
	diamonds suit = buildSuit("♢")
	clubs    suit = buildSuit("♣")
	spades   suit = buildSuit("♠")
)

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

func buildDeck() []card {
	deck := hearts
	deck = append(hearts, diamonds...)
	deck = append(deck, clubs...)
	deck = append(deck, spades...)
	return deck
}

func shuffleDeck(deck []card) []card {
	var (
		cardsRemaining  int = len(deck)
		shuffledDeck    []card
		addShuffledCard func()
	)
	usedCards := make(map[string]card)

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

func dealCard(shuffledDeck []card) (string, card, []card) {
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
	fmt.Println("Unshuffled Deck:", unshuffledDeck)

	shuffledDeck := shuffleDeck(unshuffledDeck)
	fmt.Println("Inintially Shuffled Deck:", shuffledDeck)

	var cardsToReshuffleAfter []card

	// Deal the first card
	cardStr, dealtCard, remainingDeck := dealCard(shuffledDeck)
	cardsToReshuffleAfter = append(cardsToReshuffleAfter, dealtCard)
	fmt.Println(cardStr, "Remaining Cards:", len(remainingDeck))

	// Deal the second card
	cardStr, dealtCard, remainingDeck = dealCard(remainingDeck)
	cardsToReshuffleAfter = append(cardsToReshuffleAfter, dealtCard)
	fmt.Println(cardStr, "Remaining Cards:", len(remainingDeck))

	// Take the cards that were pulled out and reshuffle them back into the deck a new
	var deckToReshuffle []card = append(remainingDeck, cardsToReshuffleAfter...)
	reshuffledDeck := shuffleDeck(deckToReshuffle)
	fmt.Println("Newly Reshuffled Deck", reshuffledDeck)
}
