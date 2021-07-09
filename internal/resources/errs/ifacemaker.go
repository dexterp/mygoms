package errs

// HandlerIface handle errors
type HandlerIface interface {
	FatalIface
	LogIface
}

// LogIface a set of functions that will log if an error is contained in v.
type LogIface interface {
	// Logf will log an error if any argument passed is an error and return true
	// to reflect an error has been found.
	Logf(format string, v ...interface{}) bool
}

// FatalIface a set of functions that will exit execution.
type FatalIface interface {
	// Panic will log an error and panic if err is not nil.
	Panic(err error)
	// Panicf will log an error and panic if any argument passed to format is an
	// error.
	Panicf(format string, v ...interface{})
	// Fatal will log an error and exit if err is not nil.
	Fatal(err error)
	// Fatalf will log an error and exit if any argument passed to fatal is an
	// error.
	Fatalf(format string, v ...interface{})
}
