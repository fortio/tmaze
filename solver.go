package main

// imagine that we're always going u,d,l, or r, and each time we hit a character we bounce off of it.
func (st *State) path() [2]int {
	st.solver[0] += st.solverDirection[0]
	st.solver[1] += st.solverDirection[1]
	// Rotate 90° perpendicular to movement: swap coords and negate one
	// char determines which one to negate
	sign := 1
	if st.maze[st.solver[0]][st.solver[1]] == runes[0] {
		sign = -1
	}
	st.solverDirection = [2]int{sign * st.solverDirection[1], sign * st.solverDirection[0]}
	return st.solver
}
