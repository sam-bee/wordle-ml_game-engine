package game

import (
	"testing"

	"github.com/sam-bee/wordle-ml_game-engine/words"
)

func TestGetFeedback(t *testing.T) {
	solution := words.Word("SPEED")
	guess := words.Word("SPARE")

	feedback := GetFeedback(solution, guess)

	expected := "GG--Y"
	got := feedback.String()

	if got != expected {
		t.Errorf("Feedback handling failed. Expected %s, got '%s'", expected, got)
	}
}
