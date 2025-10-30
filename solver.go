package main

import (
	"math/rand/v2"
)

// NewPos progresses along the path from the top left.
// Principle: imagine that we're always going u,d,l, or r,
// and each time we hit a character we bounce off of it.
func (st *State) NewPos() [2]int {
	st.solverPosition[0] += st.solverDirection[0]
	st.solverPosition[1] += st.solverDirection[1]
	sign := int(st.maze[st.solverPosition[0]][st.solverPosition[1]])
	st.solverDirection = [2]int{sign * st.solverDirection[1], sign * st.solverDirection[0]}
	return st.solverPosition
}

// GenerateMaze creates a new maze based on the current size (ap.W, ap.H).
func (st *State) GenerateMaze() {
	var idx int
	st.maze = make([][]Walls, st.ap.H)
	for l := range st.ap.H {
		line := make([]Walls, st.ap.W)
		for c := range st.ap.W {
			switch {
			case l == 0 || c+1 == st.ap.W:
				// top line or rightmost column
				idx = (l + c + 1) % 2
			case l+1 == st.ap.H || c == 0:
				// bottom line or leftmost column
				idx = (l + c) % 2
			default:
				// inside is random
				idx = rand.IntN(len(runes)) //nolint:gosec // just for visual effect
			}
			line[c] = Walls(2*idx - 1) // -1 for left, +1 for right
		}
		st.maze[l] = line
	}
}

func (st *State) ResetSolverState() {
	st.solverPosition = st.start // zero value
	st.solverDirection = [2]int{1, 0}
	st.end = [2]int{st.ap.H - 1, st.ap.W - 1}
}
