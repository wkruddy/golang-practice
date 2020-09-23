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
		cardItem := Card{value: cardValue, symbol: suitSymbol}
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

	for i := 0; count > i; i++ {

		cardToPullIndex := (len(*shuffledDeck) - 1) - i
		dealtCard := (*shuffledDeck)[:][cardToPullIndex]
		*shuffledDeck = (*shuffledDeck)[:][:cardToPullIndex]

		formattedCardNum := nameMap[dealtCard.value]
		suitSymbol := dealtCard.symbol

		cardStr := fmt.Sprintf("Dealt Card: %v of %s", formattedCardNum, suitSymbol)
		cardJustDealt := DealtCard{phrase: cardStr, card: dealtCard, remainingDeck: *shuffledDeck}
		dealtCards = append(dealtCards, cardJustDealt)
	}
	return dealtCards
}
