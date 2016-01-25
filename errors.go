package errors

import (
	"fmt"
	"github.com/golang/glog"
	"strings"
)

// Errors is a list of errors with stack traces. It implements the error interface
type Errors struct {
	errs []*Error
}

// Errf is a closure around Errorf to provide comparable but descriptive errors
type Errf func(...interface{}) error

// Add returns a list of errors that contains both parameters, no matter their error type.
func Add(e interface{}, ee interface{}) error {
	if errs, ok := e.(*Errors); ok {
		return errs.Add(ee)
	} else if err, ok := e.(*Error); ok {
		errs := &Errors{errs: make([]*Error, 0)}
		errs = errs.Add(err).(*Errors)
		return errs.Add(ee)
	} else {
		return New(ee)
	}
}

// Add returns a list of errors with the parameter added to the receiver,
// it will behave correctly with a simple error, as well as with an errors.Error and an errors.Errors as parameters.
// It will also log the error using glog if a verbosity of 3 or more is specified.
func (e *Errors) Add(ee interface{}) error {
	if ee != nil {
		var err error

		if e == nil {
			e = &Errors{errs: make([]*Error, 0)}
		}
		switch ee := ee.(type) {
		case *Error:
			err = ee
			e.errs = append(e.errs, ee)
		case *Errors:
			err = ee
			for _, err := range err.(*Errors).errs {
				e.errs = append(e.errs, err)
			}
		default:
			err = NewError(ee)
			e.errs = append(e.errs, err.(*Error))
		}
		if glog.V(3) {
			glog.Errorln(err)
		}
	} else if e == nil {
		return nil
	}
	return e
}

// Addf is a wrapper around Add to simply add a descriptive error to the list.
func (e *Errors) Addf(fmts string, args ...interface{}) error {
	return e.Add(fmt.Errorf(fmts, args...))
}

// ErrorStack returns all the stack traces and error messages of the included errors.
func (e *Errors) ErrorStack() string {
	if e == nil {
		return ""
	}
	ret := make([]string, 0)
	for i := range e.errs {
		ret = append(ret, e.errs[i].ErrorStack())
	}
	return strings.Join(ret, "\n")
}

// Error returns all the error messages of the included errors.
func (e *Errors) Error() string {
	if e == nil {
		return ""
	}
	ret := make([]string, 0)
	for i := range e.errs {
		ret = append(ret, e.errs[i].Error())
	}
	return strings.Join(ret, "\n")
}

// Is checks whether the parameter error is contained in the list of errors.
// If the parameter is an errors.Errors, it will check whether at least one of their errors match.
func (e *Errors) Is(ee error) bool {
	if e == nil && ee == nil {
		return true
	} else if e == nil || ee == nil {
		return false
	}
	if errs, ok := ee.(*Errors); ok {
		for _, err := range errs.errs {
			if e.Is(err) {
				return true
			}
		}
	} else {
		for _, err := range e.errs {
			if Is(err, ee) {
				return true
			}
		}
	}
	return false
}

// New returns a list of errors with the parameter added to the list.
func New(err interface{}) error {
	if err != nil {
		e := &Errors{errs: make([]*Error, 0)}
		return e.Add(err)
	} else {
		return nil
	}
}
