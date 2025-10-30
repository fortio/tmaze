package main

import (
	"testing"

	"fortio.org/terminal/ansipixels"
)

// mockAnsiPixels creates a minimal AnsiPixels for testing
func mockAnsiPixels(w, h int) *ansipixels.AnsiPixels {
	return &ansipixels.AnsiPixels{
		W: w,
		H: h,
	}
}

func TestGenerateMaze(t *testing.T) {
	st := &State{
		ap: mockAnsiPixels(5, 4),
	}
	st.GenerateMaze()

	// Check that maze has correct dimensions
	if len(st.maze) != st.ap.H {
		t.Errorf("Expected maze height %d, got %d", st.ap.H, len(st.maze))
	}

	for i, row := range st.maze {
		if len(row) != st.ap.W {
			t.Errorf("Expected maze width %d at row %d, got %d", st.ap.W, i, len(row))
		}
	}

	// Check that all values are either Left (-1) or Right (1)
	for y, row := range st.maze {
		for x, wall := range row {
			if wall != Left && wall != Right {
				t.Errorf("Invalid wall value at [%d][%d]: %d", y, x, wall)
			}
		}
	}

	// Verify borders follow the expected pattern (as defined in GenerateMaze)
	// The switch in GenerateMaze has priority: top/right first, then bottom/left
	for c := range st.ap.W {
		// Top border (l=0) - uses case "l == 0 || c+1 == st.ap.W"
		idx := (0 + c + 1) % 2
		expected := Walls(2*idx - 1)
		if st.maze[0][c] != expected {
			t.Errorf("Top border at column %d: expected %d, got %d", c, expected, st.maze[0][c])
		}
		// Bottom border at non-rightmost columns uses "l+1 == st.ap.H || c == 0"
		bottomRow := st.ap.H - 1
		if c < st.ap.W-1 { // not rightmost column
			idx = (bottomRow + c) % 2
			expected = Walls(2*idx - 1)
			if st.maze[bottomRow][c] != expected {
				t.Errorf("Bottom border at column %d: expected %d, got %d", c, expected, st.maze[bottomRow][c])
			}
		} else { // rightmost column uses the top/right formula
			idx = (bottomRow + c + 1) % 2
			expected = Walls(2*idx - 1)
			if st.maze[bottomRow][c] != expected {
				t.Errorf("Bottom-right corner at column %d: expected %d, got %d", c, expected, st.maze[bottomRow][c])
			}
		}
	}

	for l := range st.ap.H {
		// Left border at non-top rows uses "l+1 == st.ap.H || c == 0"
		if l > 0 { // not top row
			idx := (l + 0) % 2
			expected := Walls(2*idx - 1)
			if st.maze[l][0] != expected {
				t.Errorf("Left border at row %d: expected %d, got %d", l, expected, st.maze[l][0])
			}
		}
		// Right border uses "l == 0 || c+1 == st.ap.W"
		rightCol := st.ap.W - 1
		idx := (l + rightCol + 1) % 2
		expected := Walls(2*idx - 1)
		if st.maze[l][rightCol] != expected {
			t.Errorf("Right border at row %d: expected %d, got %d", l, expected, st.maze[l][rightCol])
		}
	}
}

func TestNewPos(t *testing.T) {
	tests := []struct {
		name         string
		initialPos   [2]int
		initialDir   [2]int
		wallAtNewPos Walls
		expectedPos  [2]int
		expectedDir  [2]int
	}{
		{
			name:         "Move right and bounce off right wall",
			initialPos:   [2]int{0, 0},
			initialDir:   [2]int{0, 1}, // moving right
			wallAtNewPos: Right,
			expectedPos:  [2]int{0, 1},
			expectedDir:  [2]int{1, 0}, // now moving down
		},
		{
			name:         "Move right and bounce off left wall",
			initialPos:   [2]int{0, 0},
			initialDir:   [2]int{0, 1}, // moving right
			wallAtNewPos: Left,
			expectedPos:  [2]int{0, 1},
			expectedDir:  [2]int{-1, 0}, // now moving up
		},
		{
			name:         "Move down and bounce off right wall",
			initialPos:   [2]int{1, 1},
			initialDir:   [2]int{1, 0}, // moving down
			wallAtNewPos: Right,
			expectedPos:  [2]int{2, 1},
			expectedDir:  [2]int{0, 1}, // now moving right
		},
		{
			name:         "Move down and bounce off left wall",
			initialPos:   [2]int{1, 1},
			initialDir:   [2]int{1, 0}, // moving down
			wallAtNewPos: Left,
			expectedPos:  [2]int{2, 1},
			expectedDir:  [2]int{0, -1}, // now moving left
		},
		{
			name:         "Move left and bounce off right wall",
			initialPos:   [2]int{2, 2},
			initialDir:   [2]int{0, -1}, // moving left
			wallAtNewPos: Right,
			expectedPos:  [2]int{2, 1},
			expectedDir:  [2]int{-1, 0}, // now moving up
		},
		{
			name:         "Move left and bounce off left wall",
			initialPos:   [2]int{2, 2},
			initialDir:   [2]int{0, -1}, // moving left
			wallAtNewPos: Left,
			expectedPos:  [2]int{2, 1},
			expectedDir:  [2]int{1, 0}, // now moving down
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple 5x5 maze
			st := &State{
				ap:              mockAnsiPixels(5, 5),
				solverPosition:  tt.initialPos,
				solverDirection: tt.initialDir,
			}

			// Initialize maze with all Right walls, then set specific wall
			st.maze = make([][]Walls, st.ap.H)
			for i := range st.ap.H {
				st.maze[i] = make([]Walls, st.ap.W)
				for j := range st.ap.W {
					st.maze[i][j] = Right
				}
			}

			// Set the wall at the position we'll move to
			newY := tt.initialPos[0] + tt.initialDir[0]
			newX := tt.initialPos[1] + tt.initialDir[1]
			st.maze[newY][newX] = tt.wallAtNewPos

			// Call NewPos
			result := st.NewPos()

			// Check position
			if result != tt.expectedPos {
				t.Errorf("Expected position %v, got %v", tt.expectedPos, result)
			}
			if st.solverPosition != tt.expectedPos {
				t.Errorf("Expected solver position %v, got %v", tt.expectedPos, st.solverPosition)
			}

			// Check direction
			if st.solverDirection != tt.expectedDir {
				t.Errorf("Expected direction %v, got %v", tt.expectedDir, st.solverDirection)
			}
		})
	}
}

func TestResetSolverState(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		initialStart [2]int
	}{
		{
			name:         "Reset with default start",
			width:        10,
			height:       8,
			initialStart: [2]int{0, 0},
		},
		{
			name:         "Reset with custom start",
			width:        15,
			height:       12,
			initialStart: [2]int{3, 5},
		},
		{
			name:         "Reset with minimal maze",
			width:        2,
			height:       2,
			initialStart: [2]int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &State{
				ap:       mockAnsiPixels(tt.width, tt.height),
				start:    tt.initialStart,
				showPath: false, // Keep false to avoid WriteString calls
			}

			// Set solver to some arbitrary position/direction before reset
			st.solverPosition = [2]int{5, 7}
			st.solverDirection = [2]int{-1, -1}

			st.ResetSolverState()

			// Check that solver position is reset to start
			if st.solverPosition != tt.initialStart {
				t.Errorf("Expected solver position %v, got %v", tt.initialStart, st.solverPosition)
			}

			// Check that direction is reset to right (down in the coordinate system)
			expectedDir := [2]int{1, 0}
			if st.solverDirection != expectedDir {
				t.Errorf("Expected solver direction %v, got %v", expectedDir, st.solverDirection)
			}

			// Check that end is set to bottom-right corner
			expectedEnd := [2]int{tt.height - 1, tt.width - 1}
			if st.end != expectedEnd {
				t.Errorf("Expected end position %v, got %v", expectedEnd, st.end)
			}
		})
	}
}
