# Game state

`game.State` is the authoritative state transition API for one Wordle game.
It has exactly six turns and owns the hidden solution, accepted-guess set,
remaining candidate solutions, and accepted guess/feedback history.

Create it with the same stable word ordering used by the caller:

```go
state, err := game.NewState(solution, validGuesses, candidateSolutions)
```

`NewState` requires a non-empty, duplicate-free valid-guess list and candidate
list. Every candidate, including the hidden solution, must be a valid guess.
It copies the supplied lists.

Call `ApplyGuess` for each move. It rejects guesses outside the valid-guess
list, repeated guesses, and every guess after the game is solved or six guesses
have been accepted. For an accepted guess it returns authoritative feedback and
filters the candidate list to words with exactly that feedback. The solution is
therefore retained in the shortlist.

`CandidateSolutions` and `History` return copies in stable turn/list order.
`TurnCount`, `Solved`, and `Finished` expose the terminal state without
revealing the hidden solution. This lets a caller convert the shortlist into
model inputs, choose a valid action, and send that action back through the
engine without duplicating Wordle transition logic.
