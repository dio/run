package run

import (
	"fmt"
	"runtime"

	"github.com/oklog/run"
)

// Group is a drop-in replacement for oklog/run.Group with automatic panic recovery.
// All panics in actors are caught and converted to PanicError.
type Group struct {
	group *run.Group

	// PanicHandler is called when an actor panics.
	// If nil, panics are converted to PanicError and returned.
	// If set, this function is called with the panic value and stack trace.
	PanicHandler func(panicValue any, stackTrace string) error
}

// New creates a new Group with panic recovery always enabled.
func New() *Group {
	return &Group{
		group: &run.Group{},
	}
}

// NewWithHandler creates a new Group with a custom panic handler.
func NewWithHandler(handler func(panicValue any, stackTrace string) error) *Group {
	return &Group{
		group:        &run.Group{},
		PanicHandler: handler,
	}
}

// Add an actor to the group with automatic panic recovery.
// This is a drop-in replacement for oklog/run.Group.Add().
func (g *Group) Add(execute func() error, interrupt func(error)) {
	safeExecute := g.wrapWithPanicRecovery(execute)
	g.group.Add(safeExecute, interrupt)
}

// Run executes all actors. This is identical to oklog/run.Group.Run().
func (g *Group) Run() error {
	return g.group.Run()
}

// wrapWithPanicRecovery wraps an execute function with panic recovery.
func (g *Group) wrapWithPanicRecovery(execute func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Capture stack trace.
				buf := make([]byte, 4096)
				stackSize := runtime.Stack(buf, false)
				stackTrace := string(buf[:stackSize])

				if g.PanicHandler != nil {
					err = g.PanicHandler(r, stackTrace)
				} else {
					err = &PanicError{
						Value:      r,
						StackTrace: stackTrace,
					}
				}
			}
		}()
		return execute()
	}
}

// PanicError represents an error that occurred during the execution of a program.
type PanicError struct {
	Value      any
	StackTrace string
}

// Error implements the error interface for PanicError.
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic: %v\n%s", e.Value, e.StackTrace)
}
