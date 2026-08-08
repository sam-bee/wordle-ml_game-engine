package game

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sam-bee/wordle-ml_game-engine/words"
)

func TestStateFiltersCandidatesUsingRepeatedLetterFeedback(t *testing.T) {
	state := newTestState(t,
		words.Word("GEESE"),
		[]words.Word{"GEESE", "EERIE", "CIGAR", "SPEED"},
		[]words.Word{"GEESE", "CIGAR", "SPEED"},
	)

	feedback, err := state.ApplyGuess(words.Word("EERIE"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := feedback.String(), "YG--G"; got != want {
		t.Fatalf("feedback = %q, want %q", got, want)
	}
	if got, want := state.CandidateSolutions(), []words.Word{"GEESE"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("candidates = %v, want %v", got, want)
	}
	if state.Solved() || state.Finished() {
		t.Fatalf("state solved=%t finished=%t after non-solution guess", state.Solved(), state.Finished())
	}
}

func TestStateShortlistMatchesBruteForceFeedbackFiltering(t *testing.T) {
	solution := words.Word("CIGAR")
	candidates := []words.Word{"CIGAR", "REBUT", "SISSY", "HUMPH", "AWAKE", "BLUSH", "FOCAL", "EVADE", "NAVAL", "SERVE"}
	guesses := []words.Word{"ARRAY", "REBUT", "SISSY", "HUMPH"}
	validGuesses := append(append([]words.Word(nil), candidates...), guesses...)
	validGuesses = uniqueWords(validGuesses)
	state := newTestState(t, solution, validGuesses, candidates)
	want := append([]words.Word(nil), candidates...)

	for _, guess := range guesses {
		feedback, err := state.ApplyGuess(guess)
		if err != nil {
			t.Fatalf("ApplyGuess(%q): %v", guess, err)
		}
		want = bruteForceShortlist(want, guess, feedback)
		if got := state.CandidateSolutions(); !reflect.DeepEqual(got, want) {
			t.Fatalf("after %q candidates = %v, want %v", guess, got, want)
		}
		if !containsWord(state.CandidateSolutions(), solution) {
			t.Fatalf("after %q candidates lost solution %q", guess, solution)
		}
	}
}

func TestStateRejectsInvalidRepeatedAndFinishedGuesses(t *testing.T) {
	solution := words.Word("CIGAR")
	validGuesses := []words.Word{"CIGAR", "REBUT", "SISSY", "HUMPH", "AWAKE", "BLUSH", "FOCAL"}
	state := newTestState(t, solution, validGuesses, []words.Word{"CIGAR", "REBUT", "SISSY"})

	if _, err := state.ApplyGuess(words.Word("XXXXX")); !errors.Is(err, ErrInvalidGuess) {
		t.Fatalf("invalid guess error = %v, want ErrInvalidGuess", err)
	}
	if got := state.TurnCount(); got != 0 {
		t.Fatalf("turn count after invalid guess = %d, want 0", got)
	}

	if _, err := state.ApplyGuess(words.Word("REBUT")); err != nil {
		t.Fatal(err)
	}
	if _, err := state.ApplyGuess(words.Word("REBUT")); !errors.Is(err, ErrRepeatedGuess) {
		t.Fatalf("repeated guess error = %v, want ErrRepeatedGuess", err)
	}
	if got := state.TurnCount(); got != 1 {
		t.Fatalf("turn count after repeated guess = %d, want 1", got)
	}

	for _, guess := range []words.Word{"SISSY", "HUMPH", "AWAKE", "BLUSH", "FOCAL"} {
		if _, err := state.ApplyGuess(guess); err != nil {
			t.Fatalf("ApplyGuess(%q): %v", guess, err)
		}
	}
	if !state.Finished() || state.Solved() {
		t.Fatalf("after six misses solved=%t finished=%t, want false/true", state.Solved(), state.Finished())
	}
	if got := state.TurnCount(); got != MaxTurns {
		t.Fatalf("turn count = %d, want %d", got, MaxTurns)
	}
	if _, err := state.ApplyGuess(solution); !errors.Is(err, ErrFinished) {
		t.Fatalf("post-finish error = %v, want ErrFinished", err)
	}
}

func TestStateFinishesWhenSolvedAndAccessorsAreDefensive(t *testing.T) {
	solution := words.Word("CIGAR")
	state := newTestState(t,
		solution,
		[]words.Word{"CIGAR", "REBUT"},
		[]words.Word{"CIGAR", "REBUT"},
	)

	candidates := state.CandidateSolutions()
	candidates[0] = words.Word("XXXXX")
	if got := state.CandidateSolutions()[0]; got != solution {
		t.Fatalf("candidate accessor leaked state: got %q, want %q", got, solution)
	}

	if _, err := state.ApplyGuess(solution); err != nil {
		t.Fatal(err)
	}
	if !state.Solved() || !state.Finished() {
		t.Fatalf("after solution solved=%t finished=%t, want true/true", state.Solved(), state.Finished())
	}
	history := state.History()
	history[0].Guess = words.Word("REBUT")
	if got := state.History()[0].Guess; got != solution {
		t.Fatalf("history accessor leaked state: got %q, want %q", got, solution)
	}
	if _, err := state.ApplyGuess(words.Word("REBUT")); !errors.Is(err, ErrFinished) {
		t.Fatalf("post-solve error = %v, want ErrFinished", err)
	}
}

func TestNewStateRejectsInvalidInitialState(t *testing.T) {
	valid := []words.Word{"CIGAR", "REBUT"}
	for name, newState := range map[string]func() (*State, error){
		"empty valid guesses": func() (*State, error) {
			return NewState(words.Word("CIGAR"), nil, []words.Word{"CIGAR"})
		},
		"empty candidates": func() (*State, error) {
			return NewState(words.Word("CIGAR"), valid, nil)
		},
		"solution not legal": func() (*State, error) {
			return NewState(words.Word("CIGAR"), []words.Word{"REBUT"}, []words.Word{"CIGAR"})
		},
		"solution not candidate": func() (*State, error) {
			return NewState(words.Word("CIGAR"), valid, []words.Word{"REBUT"})
		},
		"candidate not legal": func() (*State, error) {
			return NewState(words.Word("CIGAR"), valid, []words.Word{"CIGAR", "SISSY"})
		},
		"duplicate valid guess": func() (*State, error) {
			return NewState(words.Word("CIGAR"), []words.Word{"CIGAR", "CIGAR"}, []words.Word{"CIGAR"})
		},
		"duplicate candidate": func() (*State, error) {
			return NewState(words.Word("CIGAR"), valid, []words.Word{"CIGAR", "CIGAR"})
		},
		"malformed word": func() (*State, error) {
			return NewState(words.Word("NOPE"), valid, []words.Word{"NOPE"})
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := newState(); err == nil {
				t.Fatal("NewState unexpectedly succeeded")
			}
		})
	}
}

func newTestState(t *testing.T, solution words.Word, validGuesses, candidates []words.Word) *State {
	t.Helper()
	state, err := NewState(solution, validGuesses, candidates)
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func bruteForceShortlist(candidates []words.Word, guess words.Word, feedback Feedback) []words.Word {
	filtered := make([]words.Word, 0, len(candidates))
	for _, candidate := range candidates {
		candidateFeedback := GetFeedback(candidate, guess)
		if candidateFeedback.Equals(feedback) {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}

func uniqueWords(wordsIn []words.Word) []words.Word {
	unique := make([]words.Word, 0, len(wordsIn))
	seen := make(map[words.Word]struct{}, len(wordsIn))
	for _, word := range wordsIn {
		if _, found := seen[word]; found {
			continue
		}
		seen[word] = struct{}{}
		unique = append(unique, word)
	}
	return unique
}

func containsWord(wordsIn []words.Word, want words.Word) bool {
	for _, word := range wordsIn {
		if word == want {
			return true
		}
	}
	return false
}
