package words

import (
	"testing"
)

func TestValidGuessesReadCorrectly(t *testing.T) {
	wl, err := GetValidGuesses()
	expectedLength := 12947
	expectedFirst, _ := NewWord("AAHED")
	wordListShouldBeAsExpected(t, wl, expectedLength, expectedFirst, err)
}

func TestValidSolutionsReadCorrectly(t *testing.T) {
	wl, err := GetValidSolutions()
	expectedLength := 2309
	expectedFirst, _ := NewWord("ABACK")
	wordListShouldBeAsExpected(t, wl, expectedLength, expectedFirst, err)
}

func wordListShouldBeAsExpected(t *testing.T, wl []Word, expectedLength int, expectedFirst Word, err error) {
	gotLength := len(wl)
	gotFirst := wl[0]

	if err != nil {
		t.Errorf("Error reading word list: %s", err)
	}

	if gotLength != expectedLength {
		t.Errorf("Expected %d words, got %d", expectedLength, gotLength)
	}

	if gotFirst != expectedFirst {
		t.Errorf("Expected %q, got %q", expectedFirst, gotFirst)
	}
}
