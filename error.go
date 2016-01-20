// Package errors provides a list of errors that have stack-traces, and formattable errors that can be compared.
//
// This is particularly useful when you want to understand the
// state of execution when an error was returned unexpectedly.
//
// It provides the type *Errors which implements the standard
// golang error interface, so you can use this library interchangably
// with code that is expecting a normal error return.
//
// It also provides two extra types: *Errf and *Errors. I wanted to be able to have comparable, more helpful errors
// so *Errf is a closure for Errorf that stores an error type, so each time an error made with Newf is called, it will
// return an error that compares as true to the original Errf, while adding extra information to the error message.
// I also wanted to be able to chain errors together for the case when an error is not fatal for a function,
// but I still want to store it which is what the *Errors type does.
//
// The Error() function of the *Errors type displays the StackTrace by default.
//
// For example:
//
//  package crashy
//
//  import "github.com/soul9/errors"
//
//  var Crashed = errors.Newf("oh %s")
//
//  func Crash(s string) error {
//      return Crashed(s)
//  }
//
// This can be called as follows:
//
//  package main
//
//  import (
//      "crashy"
//      "fmt"
//      "github.com/soul9/errors"
//  )
//
//  func main() {
//      err := crashy.Crash("dear")
//      err = errors.Add(err, crashy.Crash("my"))
//      if err != nil {
//          if errors.Is(err, crashy.Crashed()) {
//              fmt.Println(err)
//          } else {
//              panic(err)
//          }
//      }
//  }
//
// This package was original written to allow reporting to Bugsnag,
// but after I found similar packages by Facebook and Dropbox, it
// was moved to one canonical location so everyone can benefit.
package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
)

var (
	// Unique identifier to make formatted errors comparable.
	errid = 1
	// The maximum number of stackframes on any error.
	MaxStackDepth = 50
)

func newid() int {
	r := errid
	errid++
	return r
}

// Error is an error with an attached stacktrace. It can be used
// wherever the builtin error interface is expected.
type Error struct {
	Err    error
	stack  []uintptr
	frames []StackFrame
	prefix string
	id     int
}

// NewError makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The stacktrace will point to the line of code that
// called New.
func NewError(e interface{}) *Error {
	var err error

	switch e := e.(type) {
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
		id:    newid(),
	}
}

// Wrap makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The skip parameter indicates how far up the stack
// to start the stacktrace. 0 is from the current call, 1 from its caller, etc.
func Wrap(e interface{}, skip int) *Error {
	var err error

	switch e := e.(type) {
	case *Error:
		return e
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
		id:    newid(),
	}
}

// WrapPrefix makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The prefix parameter is used to add a prefix to the
// error message when calling Error(). The skip parameter indicates how far
// up the stack to start the stacktrace. 0 is from the current call,
// 1 from its caller, etc.
func WrapPrefix(e interface{}, prefix string, skip int) *Error {

	err := Wrap(e, skip)

	if err.prefix != "" {
		err.prefix = fmt.Sprintf("%s: %s", prefix, err.prefix)
	} else {
		err.prefix = prefix
	}

	return err

}

// Is detects whether the error is equal to a given error. Errors
// are considered equal by this function if they are the same object,
// or if they both contain the same error inside an errors.Error,
// or if they have the same id inside an errors.Error,
// or if one of the errors in an errors.Errors Is the same as the error.
func Is(e error, original error) bool {
	if e == original {
		return true
	}
	if te, ok := e.(Errf); ok {
		return Is(te(), original)
	}
	if te, ok := original.(Errf); ok {
		return Is(e, te())
	}
	ee, eok := e.(*Error)
	ooriginal, ook := original.(*Error)

	if eok {
		if Is(ee.Err, original) {
			return true
		}
	} else if ook {
		if Is(e, ooriginal.Err) {
			return true
		}
	}
	if ook && eok && ee.id != 0 {
		return ee.id == ooriginal.id
	} else if errs, ok := e.(*Errors); ok {
		return errs.Is(original)
	} else if errs, ok := original.(*Errors); ok {
		return errs.Is(e)
	}

	return false
}

// Errorf creates a new error with the given message. You can use it
// as a drop-in replacement for fmt.Errorf() to provide descriptive
// errors in return values.
func Errorf(format string, a ...interface{}) *Error {
	return Wrap(fmt.Errorf(format, a...), 1)
}

// Error returns the underlying error's message.
func (err *Error) Error() string {

	msg := err.Err.Error()
	if err.prefix != "" {
		msg = fmt.Sprintf("%s: %s", err.prefix, msg)
	}

	return msg
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *Error) StackFrames() []StackFrame {
	if err.frames == nil {
		err.frames = make([]StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

// Newf with a format string returns a closure on top of Errorf that can be called with parameters to
// provide helpful error messages, while retaining the ability to be compared.
func Newf(s string) Errf {
	id := newid()
	return func(args ...interface{}) *Error {
		e := Errorf(s, args...)
		e.id = id
		return e
	}
}

func (e Errf) Error() string {
	return e().Error()
}
