// tmaze
// Ansi pixels version of the C64 classic maze

package main

import (
	"flag"
	"math/rand/v2"
	"os"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func main() {
	os.Exit(Main())
}

var runes = []rune{'╱', '╲'}

func (w Walls) Rune() rune {
	if w == Left {
		return runes[0]
	}
	return runes[1]
}

type State struct {
	ap              *ansipixels.AnsiPixels
	mono            bool
	newlines        bool
	showPath        bool
	width           int
	height          int
	maze            [][]Walls
	solverPosition  [2]int
	solverDirection [2]int
	start           [2]int
	end             [2]int
}

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fMono := flag.Bool("mono", false, "Use monochrome mode")
	fFPS := flag.Float64("fps", 120, "Frames per second (ansipixels rendering)")
	fNewLines := flag.Bool("nl", false, "Add newlines at end of each line (helps with copy/paste)")
	fWidth := flag.Int("width", 0, "Width of the maze (0 for full terminal width)")
	fHeight := flag.Int("height", 0, "Height of the maze (0 for full terminal height)")
	cli.Main()
	ap := ansipixels.NewAnsiPixels(*fFPS)
	ap.TrueColor = *fTrueColor
	if err := ap.Open(); err != nil {
		return 1 // error already logged
	}
	ap.HideCursor()
	defer ap.Restore()
	st := &State{
		ap:       ap,
		mono:     *fMono,
		newlines: *fNewLines,
		width:    *fWidth,
		height:   *fHeight,
	}
	if st.width > 0 {
		// need newlines when width isn't full terminal width
		st.newlines = true
	}
	ap.OnResize = func() error {
		st.GenerateMaze()
		ap.ClearScreen()
		ap.StartSyncMode()
		st.RepaintAll()
		st.ResetSolver()
		ap.EndSyncMode()
		return nil
	}
	_ = ap.OnResize() // initial draw.
	ap.MoveCursor(0, ap.H-1)
	ap.SaveCursorPos() // Ticks save cursor to prepare for where we want it on exit.
	err := ap.FPSTicks(st.Tick)
	if err != nil {
		log.Infof("Exiting on %v", err)
		return 1
	}
	return 0
}

func (st *State) RepaintAll() {
	st.ap.MoveCursor(0, 0)
	if st.mono {
		st.ap.WriteString(tcolor.Reset)
	}
	width, height := st.GetSize()
	for l := range height {
		if st.newlines && l > 0 {
			st.ap.WriteString("\r\n") // not technically needed but helps copy paste
		}
		for c := range width {
			st.EmitColor()
			st.ap.WriteRune(st.maze[l][c].Rune())
		}
	}
	st.ap.WriteString(tcolor.BrightGreen.Foreground())
}

func (st *State) ResetSolver() {
	st.ResetSolverState()
}

func (st *State) drawPath() {
	if st.solverPosition == st.end {
		// Reached the end, stop
		st.showPath = false
		return
	}
	if st.solverPosition == st.start {
		st.ap.MoveCursor(0, 0)
		st.ap.WriteRune(runes[1]) // start always right
	}
	// Move to new position
	path := st.NewPos()
	cur := st.maze[path[0]][path[1]]
	st.ap.MoveCursor(path[1], path[0])
	st.ap.WriteRune(cur.Rune())
}

func (st *State) EmitColor() {
	if st.mono {
		return
	}
	color := tcolor.Oklchf(.75, .5, rand.Float64()) //nolint:gosec // just for visual effect
	st.ap.WriteString(st.ap.ColorOutput.Foreground(color))
}

func (st *State) Tick() bool {
	if st.showPath {
		st.drawPath()
	}
	if len(st.ap.Data) == 0 {
		return true
	}
	c := st.ap.Data[0]
	switch c {
	case 'q', 'Q', 3: // Ctrl-C
		return false
	case 'c', 'C':
		st.mono = !st.mono
		st.RepaintAll()
		st.ResetSolver()
	case 'P', 'p', 'S', 's':
		st.showPath = !st.showPath
	case 'r', 'R':
		st.showPath = false
		st.RepaintAll()
		st.ResetSolver()

	default:
		// Regen a new maze on any other key
		_ = st.ap.OnResize()
	}
	return true
}
