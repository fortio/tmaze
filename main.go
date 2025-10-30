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

type State struct {
	ap       *ansipixels.AnsiPixels
	mono     bool
	newlines bool
}

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fMono := flag.Bool("mono", false, "Use monochrome mode")
	fFPS := flag.Float64("fps", 60, "Frames per second (ansipixels rendering)")
	fNewLines := flag.Bool("nl", false, "Add newlines at end of each line (helps with copy/paste)")
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
	}
	ap.OnResize = func() error {
		ap.ClearScreen()
		ap.StartSyncMode()
		// Debug the palette:
		// ap.WriteString(tcolor.Inverse)
		runes := []rune{'╱', '╲'}
		for l := range ap.H {
			if st.newlines && l > 0 {
				ap.WriteString("\r\n") // not technically needed but helps copy paste
			}
			for c := range ap.W {
				st.EmitColor(l)
				idx := rand.IntN(len(runes)) //nolint:gosec // just for visual effect
				if l == 0 || c+1 == ap.W {
					idx = l + c + 1
				} else if l+1 == ap.H || c == 0 {
					idx = l + c
				}
				ap.WriteRune(runes[idx%2])
			}
		}
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

func (st *State) EmitColor(_ int) {
	if st.mono {
		return
	}
	// Debug the palette:
	// color := tcolor.Oklchf(.7, .7, float64(line)/float64(st.ap.H)) //nolint:gosec // just for visual effect
	color := tcolor.Oklchf(.75, .5, rand.Float64()) //nolint:gosec // just for visual effect
	st.ap.WriteString(st.ap.ColorOutput.Foreground(color))
}

func (st *State) Tick() bool {
	if len(st.ap.Data) == 0 {
		return true
	}
	c := st.ap.Data[0]
	switch c {
	case 'q', 'Q', 3: // Ctrl-C
		return false
	default:
		// Regen on any other key
		_ = st.ap.OnResize()
	}
	return true
}
