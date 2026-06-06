package words

import (
	_ "embed"
	"strings"
)

//go:embed data/wordlist-valid-guesses.csv
var guessesFile string

//go:embed data/wordlist-valid-solutions.csv
var solutionsFile string

func GetValidGuesses() ([]Word, error) {
	return makeWordList(guessesFile)
}

func GetValidSolutions() ([]Word, error) {
	return makeWordList(solutionsFile)
}

func makeWordList(s string) ([]Word, error) {
	lines := strings.Split(s, "\n")
	wl := make([]Word, 0, len(lines))
	for _, line := range lines {
		w, err := NewWord(line)
		if err != nil {
			return []Word{}, err
		}
		wl = append(wl, w)
	}
	return wl, nil
}
