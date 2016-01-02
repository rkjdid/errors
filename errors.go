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
type Errf func(...interface{}) *Error

// Add returns a list of errors that contains both parameters, no matter their error type.
func Add(e interface{}, ee interface{}) *Errors {
	if errs, ok := e.(*Errors); ok {
		return errs.Add(ee)
	} else if err, ok := e.(*Error); ok {
		errs := &Errors{errs: make([]*Error, 0)}
		errs = errs.Add(err)
		return errs.Add(ee)
	} else {
		errs := New(ee)
		return errs.Add(ee)
	}
}

// Add returns a list of errors with the parameter added to the receiver,
// it will behave correctly with a simple error, as well as with an errors.Error and an errors.Errors as parameters.
// It will also log the error using glog if a verbosity of 3 or more is specified.
func (e *Errors) Add(ee interface{}) *Errors {
	if ee != nil {
		var err error

		switch ee := ee.(type) {
		case error:
			err = ee
		default:
			err = fmt.Errorf("%v", ee)
		}
		if e == nil {
			e = &Errors{errs: make([]*Error, 0)}
		}

		if ne, ok := err.(*Error); ok {
			e.errs = append(e.errs, ne)
		} else if errs, ok := err.(*Errors); ok {
			for _, err := range errs.errs {
				e = e.Add(err)
			}
		} else {
			ne = NewError(e)
			e.errs = append(e.errs, ne)
		}
		if glog.V(3) {
			glog.Errorln(err)
		}
	}
	return e
}

// Addf is a wrapper around Add to simply add a descriptive error to the list.
func (e *Errors) Addf(fmts string, args ...interface{}) *Errors {
	return e.Add(fmt.Errorf(fmts, args...))
}

// Error displays all the stack traces and error messages of the included errors.
func (e *Errors) Error() string {
	if e == nil {
		return ""
	}
	ret := make([]string, 0)
	for _, err := range e.errs {
		ret = append(ret, err.ErrorStack())
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
func New(err interface{}) *Errors {
	if err != nil {
		e := &Errors{errs: make([]*Error, 0)}
		return e.Add(err)
	} else {
		return nil
	}
}
