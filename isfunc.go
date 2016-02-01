package errors

import "os"

// IsNotExist checks for os.IsNotExist on provided error,
// whether it is an *Errors, *Error, or error
func IsNotExist(err error) bool {
	return IsFunc(os.IsNotExist, err)
}

// IsExist checks for os.IsExist on provided error,
// whether it is an *Errors, *Error, or error
func IsExist(err error) bool {
	return IsFunc(os.IsExist, err)
}

// IsPermission checks for os.IsPermission on provided error,
// whether it is an *Errors, *Error, or error
func IsPermission(err error) bool {
	return IsFunc(os.IsPermission, err)
}

// IsFunc try-casts err for *Error or *Errors,
// and checks the underlying error(s) against provided fn.
// If error is not of type *Error or *Errors, IsFunc simply calls fn(err)
func IsFunc(fn func(error) bool, err error) bool {
	switch err.(type) {
	case *Error:
		return fn(err.(*Error).Err)
	case *Errors:
		for _, errn := range err.(*Errors).errs {
			if fn(errn.Err) {
				return true
			}
		}
		return false
	default:
		return fn(err)
	}
}
