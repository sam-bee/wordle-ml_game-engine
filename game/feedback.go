package game

import (
	"github.com/sam-bee/wordle-ml_game-engine/words"
)

const (
	green int = iota
	yellow
	grey
)

type Feedback struct {
	colours []int
}

func GetFeedback(solution words.Word, guess words.Word) Feedback {
	colours := make([]int, len(solution))
	consumed := make([]bool, len(solution))

	for i := range solution {
		colours[i] = grey
		if solution[i] == guess[i] {
			colours[i] = green
			consumed[i] = true
		}
	}

	for i := range solution {
		if colours[i] == green {
			continue
		}

		for j := range solution {
			if !consumed[j] && solution[j] == guess[i] {
				colours[i] = yellow
				consumed[j] = true
				break
			}
		}
	}

	return Feedback{colours: colours}
}

func (f *Feedback) String() string {
	feedbackString := ""
	for _, colour := range f.colours {
		switch colour {
		case grey:
			feedbackString += "-"
		case yellow:
			feedbackString += "Y"
		case green:
			feedbackString += "G"
		}
	}
	return feedbackString
}

func (f *Feedback) Equals(another Feedback) bool {
	return f.String() == another.String()
}
