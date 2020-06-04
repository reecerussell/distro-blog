package result

import (
	"errors"
	"fmt"
)

// Result is a return type to be used in place of an error. Result can
// either be "Ok" or "Failure". A result can contain a value, error and
// a status code.
type Result interface {
	IsOk() bool
	Deconstruct() (success bool, status int, value interface{}, err error)

	// Chaining methods
	WithValue(v interface{}) Result
	WithStatusCode(code int) Result
}

type basicResult struct {
	ok     bool
	status int
	err    error
	value  interface{}
}

// IsOk returns a flag determining if the result is ok.
func (r *basicResult) IsOk() bool {
	return r.ok
}

// Deconstruct returns the key internal properties of the result.
func (r *basicResult) Deconstruct() (ok bool, status int, value interface{}, err error) {
	ok = r.IsOk()
	status = r.status
	value = r.value
	err = r.err

	return
}

// WithValue sets the value of the result, then returns the result
// enabling chaining.
func (r *basicResult) WithValue(v interface{}) Result {
	r.value = v

	return r
}

// WithStatusCode sets the status of the result, then returns the
// result enabling chaining.
func (r *basicResult) WithStatusCode(code int) Result {
	r.status = code

	return r
}

// Ok returns a new Ok Result.
func Ok() Result {
	return &basicResult{
		ok: true,
	}
}

// Failure returns a new non-Ok result with an error.
func Failure(err interface{}) Result {
	var e error

	switch err.(type) {
	case error:
		e = err.(error)
		break
	case string:
		e = errors.New(err.(string))
		break
	default:
		e = fmt.Errorf("%v", err)
		break
	}

	return &basicResult{
		ok:  false,
		err: e,
	}
}
