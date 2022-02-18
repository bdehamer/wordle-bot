package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

var words []string
var inLetters []rune

var posWord = [][]rune{{}, {}, {}, {}, {}}

func main() {
	rand.Seed(time.Now().Unix())
	loadWordList()
	for i := range posWord {
		for j := 'a'; j <= 'z'; j++ {
			posWord[i] = append(posWord[i], j)
		}
	}

	for {
		p := calculatePattern()
		x := makeGuess(p)
		reader := bufio.NewReader(os.Stdin)
		key, _ := reader.ReadString('\n')
		processGuess(x, key)
	}
}

func makeGuess(pattern string) string {
	hits := []string{}

	re := regexp2.MustCompile(pattern, 0)
	for _, word := range words {
		if b, _ := re.MatchString(word); b {
			hits = append(hits, word)
		}
	}

	var guess string
	switch len(hits) {
	case 0:
		fmt.Println("Sorry -- couldn't find the word")
		os.Exit(1)
		break
	case 1:
		fmt.Println(hits[0])
		os.Exit(0)
		break
	default:
		i := rand.Intn(len(hits))
		guess = hits[i]
		fmt.Printf("%s (%d options)\n", guess, len(hits))
	}

	return guess
}

func calculatePattern() string {
	pattern := ""
	for _, letters := range posWord {
		pattern = pattern + fmt.Sprintf("[%s]", string(letters))
	}
	pattern = fmt.Sprintf("(?=%s)", pattern)

	if len(inLetters) > 0 {
		for _, c := range inLetters {
			pattern = pattern + fmt.Sprintf("(?=.*%c.*)", c)
		}
	}
	return pattern

}

func processGuess(guess, key string) {
	for i := range guess {
		ltr := rune(guess[i])
		s := key[i]
		switch s {
		case 'x':
			removePossibleLetterFromAllPositions(ltr)
			break
		case 'o':
			affirmLetter(ltr)
			break
		case '^':
			confirmLetterPosition(ltr, i)
			break
		}
	}
}

func confirmLetterPosition(c rune, i int) {
	posWord[i] = []rune{c}
}

func removePossibleLetterFromAllPositions(c rune) {
	for i := range posWord {
		removePossibleLetterFromPosition(c, i)
	}
}

func removePossibleLetterFromPosition(c rune, i int) {
	ary := posWord[i]

	idx := findIndexOfLetter(ary, c)

	if idx != -1 {
		posWord[i] = removeEntryFromList(ary, idx)
	}
}

func findIndexOfLetter(ary []rune, c rune) int {
	matchIdx := -1
	for idx, ltr := range ary {
		if ltr == c {
			matchIdx = idx
			break
		}
	}
	return matchIdx
}

func removeEntryFromList(ary []rune, i int) []rune {
	ary[i] = ary[len(ary)-1]
	ary = ary[:len(ary)-1]
	return ary
}

func affirmLetter(c rune) {
	existing := false
	for _, x := range inLetters {
		if x == c {
			existing = true
			break
		}
	}

	if !existing {
		inLetters = append(inLetters, c)
	}
}

func loadWordList() {
	f, err := os.OpenFile("words.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return
		}

		words = append(words, strings.TrimSpace(line))
	}
}
