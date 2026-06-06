package words

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidGuessesReadCorrectly(t *testing.T) {
	withWordListServer(t, map[string]string{
		validGuessesFilename: "AAHED\nAALII",
	})

	wl, err := GetValidGuesses()
	expectedLength := 2
	expectedFirst, _ := NewWord("AAHED")
	wordListShouldBeAsExpected(t, wl, expectedLength, expectedFirst, err)
}

func TestValidSolutionsReadCorrectly(t *testing.T) {
	withWordListServer(t, map[string]string{
		validSolutionsFilename: "ABACK\nABASE",
	})

	wl, err := GetValidSolutions()
	expectedLength := 2
	expectedFirst, _ := NewWord("ABACK")
	wordListShouldBeAsExpected(t, wl, expectedLength, expectedFirst, err)
}

func TestWordListsUseTaggedSource(t *testing.T) {
	if !strings.Contains(wordListBaseURL, "/"+wordListVersion+"/") {
		t.Errorf("Expected word list URL %q to include version %q", wordListBaseURL, wordListVersion)
	}
}

func TestWordListFetchFailureReturnsError(t *testing.T) {
	withWordListServer(t, map[string]string{})

	wl, err := GetValidGuesses()

	if err == nil {
		t.Errorf("Expected error fetching missing word list")
	}
	if len(wl) != 0 {
		t.Errorf("Expected no words, got %d", len(wl))
	}
}

func wordListShouldBeAsExpected(t *testing.T, wl []Word, expectedLength int, expectedFirst Word, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Error reading word list: %s", err)
	}

	gotLength := len(wl)
	gotFirst := wl[0]

	if gotLength != expectedLength {
		t.Errorf("Expected %d words, got %d", expectedLength, gotLength)
	}

	if gotFirst != expectedFirst {
		t.Errorf("Expected %q, got %q", expectedFirst, gotFirst)
	}
}

func withWordListServer(t *testing.T, files map[string]string) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/")
		contents, ok := files[filename]
		if !ok {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, contents)
	}))

	originalBaseURL := wordListBaseURL
	wordListBaseURL = server.URL

	t.Cleanup(func() {
		wordListBaseURL = originalBaseURL
		server.Close()
	})
}
