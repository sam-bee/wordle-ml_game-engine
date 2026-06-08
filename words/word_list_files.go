package words

import (
	"strings"

	wordlists "github.com/sam-bee/wordle-ml_wordlists"
)

func GetValidGuesses() ([]Word, error) {
	return makeWordList(wordlists.ValidGuessesCSV())
}

func GetValidSolutions() ([]Word, error) {
	return makeWordList(wordlists.ValidSolutionsCSV())
}

func GetActionSpace() ([]Word, error) {
	return makeWordList(wordlists.ActionSpaceCSV())
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
