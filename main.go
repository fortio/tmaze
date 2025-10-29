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
)

func main() {
	os.Exit(Main())
}

type State struct {
	ap *ansipixels.AnsiPixels
}

func Main() int {
	truecolorDefault := ansipixels.DetectColorMode().TrueColor
	fTrueColor := flag.Bool("truecolor", truecolorDefault,
		"Use true color (24-bit RGB) instead of 8-bit ANSI colors (default is true if COLORTERM is set)")
	fFPS := flag.Float64("fps", 60, "Frames per second (ansipoixels rendering)")
	cli.Main()
	ap := ansipixels.NewAnsiPixels(*fFPS)
	ap.TrueColor = *fTrueColor
	if err := ap.Open(); err != nil {
		return 1 // error already logged
	}
	ap.HideCursor()
	defer ap.Restore()
	ap.OnResize = func() error {
		ap.ClearScreen()
		ap.StartSyncMode()
		runes := []rune{'╱', '╲'}
		for range ap.H * ap.W {
			idx := rand.IntN(len(runes))
			ap.WriteRune(runes[idx])
		}
		ap.EndSyncMode()
		return nil
	}
	_ = ap.OnResize() // initial draw.
	st := State{ap: ap}
	err := ap.FPSTicks(st.Tick)
	if err != nil {
		log.Infof("Exiting on %v", err)
		return 1
	}
	return 0
}

func (st *State) Tick() bool {
	if len(st.ap.Data) == 0 {
		return true
	}
	c := st.ap.Data[0]
	switch c {
	case 'q', 'Q', 3: // Ctrl-C
		log.Infof("Exiting on %q", c)
		return false
	default:
		log.Debugf("Input %q...", c)
		// Do something
	}
	return true
}
