package cmd

import (
	"io"
	"os"
	"os/exec"
	"runtime"
)

// Pager tries to start a pager process (less/more) for better output display.
// If it succeeds, returns a writer connected to the pager and a cleanup func.
// If it fails, it falls back to os.Stdout with a no-op cleanup.
//
// The pager is useful for displaying large amounts of text (like configuration
// documentation) in a scrollable interface rather than dumping everything to
// the terminal at once.
func Pager() (io.Writer, func()) {
	pagerCmd := determinePagerCommand()

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

// determinePagerCommand returns the appropriate pager command for the current platform.
// On Windows, defaults to "more". On Unix-like systems, uses the PAGER environment
// variable or defaults to "less".
func determinePagerCommand() string {
	if runtime.GOOS == "windows" {
		return "more"
	}

	// Check for user-specified pager
	if pagerCmd := os.Getenv("PAGER"); pagerCmd != "" {
		return pagerCmd
	}

	// Default to less on Unix-like systems
	return "less"
}
