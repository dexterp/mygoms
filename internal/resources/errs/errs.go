// kick:render
package errs

import (
	"fmt"
	"os"

	"log"

	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/exit"
	"${GOSERVER}/${GOGROUP}/${PROJECT_NAME}/internal/resources/logger"
)

// Handler error handling
type Handler struct {
	ex  exit.HandlerIface  `validate:"required"` // Exit handler
	log logger.OutputIface `validate:"required"` // Default logger
}

// Options to New
type Options struct {
	ExitHandler exit.HandlerIface  // Exit handler
	Logger      logger.OutputIface // Logger
}

// New create new handler
func New(opts Options) *Handler {
	return &Handler{
		ex:  opts.ExitHandler,
		log: opts.Logger,
	}
}

// Panic will log an error and panic if err is not nil.
func (e *Handler) Panic(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// Panicf will log an error and panic if any argument passed to format is an error
func (e *Handler) Panicf(format string, v ...interface{}) {
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// Logf will log an error if any argument passed to format is an error
func (e *Handler) Logf(format string, v ...interface{}) bool { // nolint
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func (e *Handler) Fatal(err error) {
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// Fatalf will log an error and exit if any argument passed to fatal is an error
func (e *Handler) Fatalf(format string, v ...interface{}) { // nolint
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func (e *Handler) hasErrPrint(err error) bool {
	if err == nil {
		return false
	}
	o := e.log.Output(3, err.Error())
	if o != nil {
		panic(o)
	}
	return true
}

func (e *Handler) hasErrPrintf(format string, v ...interface{}) bool {
	hasError := false
	for _, elm := range v {
		if _, ok := elm.(error); ok {
			hasError = true
			break
		}
	}
	if !hasError {
		return false
	}
	out := fmt.Errorf(format, v...)
	e.log.Output(3, out.Error()) // nolint
	return true
}

// Panic will log an error and panic if err is not nil.
func Panic(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	panic(err)
}

// Panicf will log an error and panic if any argument passed to format is an error
func Panicf(format string, v ...interface{}) {
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	panic(fmt.Errorf(format, v...))
}

// Logf will log an error if any argument passed to format is an error
func Logf(format string, v ...interface{}) bool { // nolint
	e := makeErrors()
	return e.hasErrPrintf(format, v...)
}

// Fatal will log an error and exit if err is not nil.
func Fatal(err error) {
	e := makeErrors()
	has := e.hasErrPrint(err)
	if !has {
		return
	}
	e.ex.Exit(255)
}

// Fatalf will log an error and exit if any argument passed to fatal is an error
func Fatalf(format string, v ...interface{}) { // nolint
	e := makeErrors()
	hasErr := e.hasErrPrintf(format, v...)
	if !hasErr {
		return
	}
	e.ex.Exit(255)
}

func makeErrors() *Handler {
	eh := exit.New(exit.Options{
		Mode: exit.ExitMode,
	})
	e := &Handler{
		ex:  eh,
		log: logger.New(os.Stderr, "", log.LstdFlags, logger.ErrorLevel, eh),
	}
	return e
}
