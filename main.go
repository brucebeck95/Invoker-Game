package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

type Spell struct {
	name           string
	keyCombination string
}

var spellList = []Spell{
	{name: "Sunstrike", keyCombination: "EEE"},
	{name: "Ghostwalk", keyCombination: "QQW"},
	{name: "Chaos Meteor", keyCombination: "EEW"},
	{name: "Forge Spirit", keyCombination: "QEE"},
	{name: "Ice Wall", keyCombination: "QQE"},
	{name: "Alacrity", keyCombination: "WWE"},
	{name: "EMP", keyCombination: "WWW"},
	{name: "Tornado", keyCombination: "WWQ"},
	{name: "Cold Snap", keyCombination: "QQQ"},
	{name: "Deafening Blast", keyCombination: "QWE"},
}


func getRandomSpell(numSpells int) []Spell {

	n := len(spellList)
	if n == 0 {
		return []Spell{}
	}

	shuffledSpells := make([]Spell, n)
	copy(shuffledSpells, spellList)
	rand.Shuffle(len(shuffledSpells), func(i, j int) {
		shuffledSpells[i], shuffledSpells[j] = shuffledSpells[j], shuffledSpells[i]
	})

	if numSpells >= n {
		return shuffledSpells
	}
	return shuffledSpells[:numSpells]
}

func areAnagrams(firstString, secondString string) bool {
	if len(firstString) != len(secondString) {
		return false
	}
	s1Chars := strings.Split(firstString, "")
	s2Chars := strings.Split(secondString, "")
	sort.Strings(s1Chars)
	sort.Strings(s2Chars)
	return strings.Join(s1Chars, "") == strings.Join(s2Chars, "")
}

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making terminal raw: %v\n", err)
		os.Exit(1)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Print("\r\nGame interrupted. Exiting.\r\n")
		os.Exit(1)
	}()

	fmt.Print("Welcome to the Invoker game!\r\n")
	fmt.Print("Get ready. The game is about to begin in 5 seconds...\r\n")

	fmt.Print("\r\n--- Instructions ---\r\n")
	fmt.Print("- Press 'Q', 'W', or 'E' to choose orbs.\r\n")
	fmt.Print("- Your last 3 orbs will be used for invocation.\r\n")
	fmt.Print("- Press 'R' to invoke the spell.\r\n")
	fmt.Print("- Press 'C' to clear your current orb selection.\r\n")
	fmt.Print("- Press 'X' to quit the game.\r\n\r\n")

	for i := 5; i > 0; i-- {
		fmt.Printf("%d...\r\n", i)
		time.Sleep(1 * time.Second)
	}

	numSpellsToPractice := 10
	spells := getRandomSpell(numSpellsToPractice)

	if len(spells) == 0 {
		fmt.Print("No spells available to practice. Exiting.\r\n")
		return
	}

	startTime := time.Now()
	score := 0

OuterLoop:
	for i, spell := range spells {
		var currentKeyQueue []byte

		fmt.Printf("\r\n\033[K--- Spell %d of %d: %s ---\r\n", i+1, len(spells), spell.name)
		fmt.Print("\033[KEnter orbs: ")

	InnerLoop:
		for {
			var b = make([]byte, 1)
			_, err := os.Stdin.Read(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\r\nError reading input: %v\r\n", err)
				break OuterLoop
			}
			key := b[0]

			switch key {
			case 'r', 'R':
				fmt.Print("\r\n\033[K[R] Invoke attempt: ")

				if len(currentKeyQueue) < len(spell.keyCombination) {
					fmt.Printf("Not enough orbs (%d/%d). Current: %s\r\n", len(currentKeyQueue), len(spell.keyCombination), string(currentKeyQueue))
					fmt.Print("\033[KEnter orbs: ", string(currentKeyQueue))
				} else {
					invokedCombination := string(currentKeyQueue)

					if areAnagrams(invokedCombination, spell.keyCombination) {
						fmt.Printf("SUCCESS! You invoked %s.\r\n", spell.name)
						score++
						currentKeyQueue = []byte{}
						break InnerLoop
					} else {
						fmt.Printf("FAILED. Expected anagram of %s, got %s.\r\n", spell.keyCombination, invokedCombination)
						currentKeyQueue = []byte{}
						break InnerLoop
					}
				}
			case 'x', 'X':
				fmt.Print("\r\n\033[KQuitting the game.\r\n")

				return
			case 'c', 'C':
				currentKeyQueue = []byte{}
				fmt.Print("\r\033[KOrbs cleared. Enter orbs: ")
			case 'q', 'Q', 'w', 'W', 'e', 'E':
				orb := strings.ToUpper(string(key))

				currentKeyQueue = append(currentKeyQueue, orb[0])

				expectedKeyLength := len(spell.keyCombination)
				if len(currentKeyQueue) > expectedKeyLength {
					currentKeyQueue = currentKeyQueue[len(currentKeyQueue)-expectedKeyLength:]
				}

				fmt.Printf("\r\033[KEnter orbs: %s", string(currentKeyQueue))

			default:
			
			}
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()

	fmt.Printf("\r\n\033[K--- Game Over! ---\r\n")
	fmt.Printf("\033[KYou attempted %d spells.\r\n", len(spells))
	fmt.Printf("\033[KYour score: %d out of %d.\r\n", score, len(spells))
	fmt.Printf("\033[KTotal time: %.2f seconds.\r\n", duration)
}
