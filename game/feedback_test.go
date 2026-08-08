package game

import (
	"testing"

	"github.com/sam-bee/wordle-ml_game-engine/words"
)

func TestGetFeedback(t *testing.T) {
	tests := []struct {
		name     string
		solution words.Word
		guess    words.Word
		expected string
	}{
		{
			name:     "existing case",
			solution: words.Word("SPEED"),
			guess:    words.Word("SPARE"),
			expected: "GG--Y",
		},
		{
			name:     "duplicate letters in solution and guess",
			solution: words.Word("GEESE"),
			guess:    words.Word("EERIE"),
			expected: "YG--G",
		},
		{
			name:     "excess repeated letters in guess",
			solution: words.Word("CIGAR"),
			guess:    words.Word("ARRAY"),
			expected: "-Y-G-",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			feedback := GetFeedback(test.solution, test.guess)
			got := feedback.String()
			if got != test.expected {
				t.Errorf("expected %s, got %s", test.expected, got)
			}
		})
	}
}
