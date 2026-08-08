package game

import (
	"errors"
	"fmt"

	"github.com/sam-bee/wordle-ml_game-engine/words"
)

// MaxTurns is the fixed number of guesses in one Wordle game.
const MaxTurns = 6

var (
	// ErrInvalidGuess indicates that a guess is not in the game's valid-guess list.
	ErrInvalidGuess = errors.New("guess is not valid for this game")
	// ErrRepeatedGuess indicates that a guess has already been played.
	ErrRepeatedGuess = errors.New("guess has already been played")
	// ErrFinished indicates that no more guesses can be played.
	ErrFinished = errors.New("game is already finished")
)

// Turn records one accepted guess and the feedback it received.
type Turn struct {
	Guess    words.Word
	Feedback Feedback
}

// State is one fixed six-turn Wordle game. It owns the hidden solution, the
// guesses accepted by the game, the remaining candidate solutions, and the
// completed turns.
type State struct {
	solution     words.Word
	validGuesses map[words.Word]struct{}
	candidates   []words.Word
	guessed      map[words.Word]struct{}
	history      []Turn
	solved       bool
	finished     bool
}

// NewState creates a game with a hidden solution, its accepted guesses, and
// its initial candidate solutions. The solution must be both a valid guess and
// an initial candidate. Lists are copied, so later caller mutations do not
// alter the game.
func NewState(solution words.Word, validGuesses, candidateSolutions []words.Word) (*State, error) {
	if err := validateStateWord("solution", solution); err != nil {
		return nil, err
	}
	if len(validGuesses) == 0 {
		return nil, errors.New("valid guesses must not be empty")
	}
	if len(candidateSolutions) == 0 {
		return nil, errors.New("candidate solutions must not be empty")
	}

	valid := make(map[words.Word]struct{}, len(validGuesses))
	for _, guess := range validGuesses {
		if err := validateStateWord("valid guess", guess); err != nil {
			return nil, err
		}
		if _, duplicate := valid[guess]; duplicate {
			return nil, fmt.Errorf("valid guesses contains duplicate %q", guess)
		}
		valid[guess] = struct{}{}
	}
	if _, ok := valid[solution]; !ok {
		return nil, fmt.Errorf("solution %q is not a valid guess", solution)
	}

	candidates := make([]words.Word, len(candidateSolutions))
	seenCandidates := make(map[words.Word]struct{}, len(candidateSolutions))
	solutionPresent := false
	for index, candidate := range candidateSolutions {
		if err := validateStateWord("candidate solution", candidate); err != nil {
			return nil, err
		}
		if _, ok := valid[candidate]; !ok {
			return nil, fmt.Errorf("candidate solution %q is not a valid guess", candidate)
		}
		if _, duplicate := seenCandidates[candidate]; duplicate {
			return nil, fmt.Errorf("candidate solutions contains duplicate %q", candidate)
		}
		seenCandidates[candidate] = struct{}{}
		candidates[index] = candidate
		solutionPresent = solutionPresent || candidate == solution
	}
	if !solutionPresent {
		return nil, fmt.Errorf("candidate solutions does not contain solution %q", solution)
	}

	return &State{
		solution:     solution,
		validGuesses: valid,
		candidates:   candidates,
		guessed:      make(map[words.Word]struct{}, MaxTurns),
		history:      make([]Turn, 0, MaxTurns),
	}, nil
}

// ApplyGuess records one legal, previously unused guess, returns its feedback,
// and updates the candidate shortlist. It rejects guesses after a solved game
// or after the sixth guess.
func (s *State) ApplyGuess(guess words.Word) (Feedback, error) {
	if s.finished {
		return Feedback{}, ErrFinished
	}
	if _, ok := s.validGuesses[guess]; !ok {
		return Feedback{}, fmt.Errorf("%w: %q", ErrInvalidGuess, guess)
	}
	if _, repeated := s.guessed[guess]; repeated {
		return Feedback{}, fmt.Errorf("%w: %q", ErrRepeatedGuess, guess)
	}

	feedback := GetFeedback(s.solution, guess)
	updatedCandidates := make([]words.Word, 0, len(s.candidates))
	solutionRetained := false
	for _, candidate := range s.candidates {
		candidateFeedback := GetFeedback(candidate, guess)
		if !candidateFeedback.Equals(feedback) {
			continue
		}
		updatedCandidates = append(updatedCandidates, candidate)
		solutionRetained = solutionRetained || candidate == s.solution
	}
	if !solutionRetained {
		return Feedback{}, errors.New("internal error: feedback removed the solution from candidates")
	}

	s.candidates = updatedCandidates
	s.guessed[guess] = struct{}{}
	s.history = append(s.history, Turn{Guess: guess, Feedback: feedback})
	s.solved = guess == s.solution
	s.finished = s.solved || len(s.history) == MaxTurns
	return feedback, nil
}

// TurnCount returns the number of accepted guesses made so far.
func (s *State) TurnCount() int {
	return len(s.history)
}

// CandidateSolutions returns the current shortlist in its initial stable order.
func (s *State) CandidateSolutions() []words.Word {
	return append([]words.Word(nil), s.candidates...)
}

// History returns a copy of the completed game history in turn order.
func (s *State) History() []Turn {
	return append([]Turn(nil), s.history...)
}

// Solved reports whether the hidden solution has been guessed.
func (s *State) Solved() bool {
	return s.solved
}

// Finished reports whether the game has been solved or reached six guesses.
func (s *State) Finished() bool {
	return s.finished
}

func validateStateWord(name string, word words.Word) error {
	if _, err := words.NewWord(string(word)); err != nil {
		return fmt.Errorf("%s %q: %w", name, word, err)
	}
	return nil
}
