# Wordle ML Game Engine

A small Go library for Wordle game rules and word handling. `game.State`
provides the authoritative fixed six-turn game state, legal-guess checks,
feedback, candidate-shortlist updates, and turn history. See
[`docs/game-state.md`](docs/game-state.md) for its API contract.

This repository is part of the Wordle ML project. It is being stripped down from an older Wordle player into a reusable game engine library.

Wordlist data is provided by `github.com/sam-bee/wordle-ml_wordlists`.

## Go Version

Go 1.26.

## Module

```text
github.com/sam-bee/wordle-ml_game-engine
```
