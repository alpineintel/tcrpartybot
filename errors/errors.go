package errors

import (
	"fmt"
	"log"
	"runtime"
)

// TCRPError is an interface for errors augmented with stacktrace information
type TCRPError interface {
	error
	ErrorWithStacktrace() string
}

type tcrpError struct {
	msg        string
	stacktrace []functionCall
}

type functionCall struct {
	File string
	Line int
}

// New augments an error message with stacktrace information
func New(msg string) TCRPError {
	return &tcrpError{msg: msg, stacktrace: constructTrace()}
}

// Wrap augments an existing error with stacktrace information
func Wrap(err error) TCRPError {
	// Don't double wrap errors
	if tcrpError, ok := err.(TCRPError); ok {
		return tcrpError
	}

	return &tcrpError{msg: err.Error(), stacktrace: constructTrace()}
}

func constructTrace() []functionCall {
	pc := make([]uintptr, 10)
	length := runtime.Callers(3, pc)

	stacktrace := make([]functionCall, length-1)
	for i := 0; i < length-1; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		stacktrace[i] = functionCall{File: file, Line: line}
	}

	return stacktrace
}

// Errorf formats a given error message and attaches stacktrace info
func Errorf(msg string, rest ...interface{}) TCRPError {
	msg = fmt.Sprintf(msg, rest...)
	return &tcrpError{msg: msg, stacktrace: constructTrace()}

}

// Error returns the error message associated with a TCRPError
func (e *tcrpError) Error() string {
	return e.msg
}

// ErrorWithStacktrace returns the error message in addition to a stacktrace in
// the form of a printable string
func (e *tcrpError) ErrorWithStacktrace() string {
	str := e.msg + "\n"
	for _, call := range e.stacktrace {
		str += fmt.Sprintf("\t%s:%d\n", call.File, call.Line)
	}
	return str
}

// LogErrors properly formats incoming errors from an errChan
func LogErrors(errChan <-chan error) {
	for err := range errChan {
		if tcrpError, ok := err.(TCRPError); ok {
			log.Printf("[error] %s", tcrpError.ErrorWithStacktrace())
		} else {
			log.Printf("[error] %s", err.Error())
		}
	}
}
