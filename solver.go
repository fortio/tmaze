package main

// NewPos progress along the path from the top left.
// Principle: imagine that we're always going u,d,l, or r,
// and each time we hit a character we bounce off of it.
func (st *State) NewPos() [2]int {
	st.solverPosition[0] += st.solverDirection[0]
	st.solverPosition[1] += st.solverDirection[1]
	sign := int(st.maze[st.solverPosition[0]][st.solverPosition[1]])
	st.solverDirection = [2]int{sign * st.solverDirection[1], sign * st.solverDirection[0]}
	return st.solverPosition
}
