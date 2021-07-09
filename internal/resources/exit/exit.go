// Package exit handles exit to OS
package exit

import (
	"fmt"
	"os"
)

const (
	// MNone set the exit mode at default. When Exit.Exit is called it exits
	// with the as normal with the supplied return code.
	MNone = iota
	// MPanic set the exit mode to panic. When Exit.Exit is called then a panic
	// will ensue instead of exiting.
	MPanic
)

var (
	// ExitMode is the current exit mode. Should be one of MNone or MPanic
	ExitMode int
)

// Handler exit handling
type Handler struct {
	mode int
}

// Options options to New
type Options struct {
	Mode int
}

// New create a Handler object
func New(opts Options) *Handler {
	return &Handler{
		mode: opts.Mode,
	}
}

// Exit exit function with exit code. 0 is success
func (e *Handler) Exit(code int) {
	switch e.mode {
	case MNone:
		os.Exit(code)
	case MPanic:
		panic(fmt.Sprintf("Exit %d\n", code))
	default:
		panic(fmt.Sprintf("Unknown exit mode: %d", e.mode))
	}
}

// Exit exit function with exit code. 0 is success
func Exit(code int) {
	h := Handler{
		mode: ExitMode,
	}
	h.Exit(code)
}

// Mode sets the exit mode for the exit.Exit call. Current modes are,
// exit.None (default - call exit mode directly) & exit.Panic (panic instead of exiting)
// Returns true if exit mode was was successfully set.
func Mode(mode int) (ok bool) {
	switch mode {
	case MNone:
		ExitMode = MNone
		ok = true
	case MPanic:
		ExitMode = MPanic
		ok = true
	default:
		panic(fmt.Sprintf("Unknown exit mode: %d", mode))
	}
	return
}
