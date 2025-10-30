package main

// imagine that we're always going u,d,l, or r, and each time we hit a character we bounce off of it
func path(maze [][]rune) [][2]int {
	runes := []rune{'╱', '╲'}
	cur := [2]int{0, 0}
	curDirection := [2]int{1, 0}
	path := [][2]int{}

	for cur != [2]int{len(maze) - 1, len(maze[0]) - 1} {
		path = append(path, cur)
		cur[0] += curDirection[0]
		cur[1] += curDirection[1]
		char := maze[cur[0]][cur[1]]
		// this could be cleaned up
		switch curDirection {
		case [2]int{0, 1}:
			curDirection[1] = 0
			curDirection[0] = 1
			if char == runes[0] {
				curDirection[0] = -1
			}
		case [2]int{0, -1}:
			curDirection[1] = 0
			curDirection[0] = -1
			if char == runes[0] {
				curDirection[0] = 1
			}
		case [2]int{1, 0}:
			curDirection[0] = 0
			curDirection[1] = 1
			if char == runes[0] {
				curDirection[1] = -1
			}
		case [2]int{-1, 0}:
			curDirection[0] = 0
			curDirection[1] = -1
			if char == runes[0] {
				curDirection[1] = 1
			}
		}
	}
	return path
}
