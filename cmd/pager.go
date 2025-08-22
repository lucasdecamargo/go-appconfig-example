package cmd

import (
	"io"
	"os"
	"os/exec"
	"runtime"
)

// [Pager] tries to start a pager process (less/more).
// If it succeeds, returns a writer connected to the pager and a cleanup func.
// If it fails, it falls back to os.Stdout with a no-op cleanup.
func Pager() (io.Writer, func()) {
	var pagerCmd string
	if runtime.GOOS == "windows" {
		pagerCmd = "more"
	} else {
		pagerCmd = os.Getenv("PAGER")
		if pagerCmd == "" {
			pagerCmd = "less"
		}
	}

	cmd := exec.Command(pagerCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	w, err := cmd.StdinPipe()
	if err == nil && cmd.Start() == nil {
		// Pager successfully started
		cleanup := func() {
			_ = w.Close()  // close writer so pager sees EOF
			_ = cmd.Wait() // wait for pager to exit
		}
		return w, cleanup
	}

	// Fallback: just print to stdout
	return os.Stdout, func() {}
}
