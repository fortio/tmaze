package main

import "time"

// imagine that we're always going u,d,l, or r, and each time we hit a character we bounce off of it.
func (st *State) path() <-chan [2]int {
	pathChan := make(chan [2]int)
	go func() {
		defer close(pathChan)
		cur := [2]int{0, 0}
		curDirection := [2]int{1, 0}
		ticker := time.NewTicker(time.Millisecond)
		defer ticker.Stop()
		for cur != [2]int{len(st.maze) - 1, len(st.maze[0]) - 1} {
			pathChan <- cur
			<-ticker.C
			oldDirection := curDirection
			cur[0] += curDirection[0]
			cur[1] += curDirection[1]
			// Rotate 90° perpendicular to movement: swap coords and negate one
			// char determines which one to negate
			sign := 1
			if st.maze[cur[0]][cur[1]] == runes[0] {
				sign = -1
			}
			curDirection = [2]int{sign * oldDirection[1], sign * oldDirection[0]}
		}
	}()
	return pathChan
}
