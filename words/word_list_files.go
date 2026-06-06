package words

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	wordListVersion = "v0.1.0"

	validGuessesFilename   = "wordlist-valid-guesses.csv"
	validSolutionsFilename = "wordlist-valid-solutions.csv"
)

var (
	wordListBaseURL = "https://raw.githubusercontent.com/sam-bee/wordle-ml_wordlists/" + wordListVersion + "/data"
	wordListClient  = &http.Client{Timeout: 10 * time.Second}
)

func GetValidGuesses() ([]Word, error) {
	return getWordList(validGuessesFilename)
}

func GetValidSolutions() ([]Word, error) {
	return getWordList(validSolutionsFilename)
}

func getWordList(filename string) ([]Word, error) {
	response, err := wordListClient.Get(wordListBaseURL + "/" + filename)
	if err != nil {
		return []Word{}, fmt.Errorf("fetch word list %q: %w", filename, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Word{}, fmt.Errorf("fetch word list %q: got HTTP %d", filename, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []Word{}, fmt.Errorf("read word list %q: %w", filename, err)
	}

	return makeWordList(string(body))
}

func makeWordList(s string) ([]Word, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	wl := make([]Word, 0, len(lines))
	for _, line := range lines {
		w, err := NewWord(strings.TrimSpace(line))
		if err != nil {
			return []Word{}, err
		}
		wl = append(wl, w)
	}
	return wl, nil
}
