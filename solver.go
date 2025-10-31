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

func (st *State) GetSize() (width, height int) {
	width = st.width
	if width <= 0 {
		width = st.ap.W
	}
	height = st.height
	if height <= 0 {
		height = st.ap.H
	}
	return width, height
}

// GenerateMaze creates a new maze based on the current size (ap.W, ap.H).
func (st *State) GenerateMaze() {
	var idx int
	width, height := st.GetSize()
	st.maze = make([][]Walls, height)
	for l := range height {
		line := make([]Walls, width)
		for c := range width {
			switch {
			case l == 0 || c+1 == width:
				// top line or rightmost column
				idx = (l + c + 1) % 2
			case l+1 == height || c == 0:
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
	width, height := st.GetSize()
	st.end = [2]int{height - 1, width - 1}
}
