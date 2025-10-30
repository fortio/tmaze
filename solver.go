package main

import "time"

// imagine that we're always going u,d,l, or r, and each time we hit a character we bounce off of it.
func path(maze [][]rune) <-chan [2]int {
	pathChan := make(chan [2]int)
	go func() {
		defer close(pathChan)
		runes := []rune{'╱', '╲'}
		cur := [2]int{0, 0}
		curDirection := [2]int{1, 0}

		for cur != [2]int{len(maze) - 1, len(maze[0]) - 1} {
			pathChan <- cur
			time.Sleep(1 * time.Nanosecond)
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
	}()
	return pathChan
}
