package result

import (
	"errors"
	"strconv"
	"testing"
)

func TestOk(t *testing.T) {
	r := Ok()

	success, status, value, err := r.Deconstruct()
	if !success {
		t.Errorf("expected success to be true")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err != nil {
		t.Errorf("expected err to be nil, but got: %v", err)
	}
}

func TestOkWithStatus(t *testing.T) {
	r := Ok().WithStatusCode(200)

	success, status, value, err := r.Deconstruct()
	if !success {
		t.Errorf("expected success to be true")
	}

	if status != 200 {
		t.Errorf("expected status to be 200, but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err != nil {
		t.Errorf("expected err to be nil, but got: %v", err)
	}
}

func TestOkWithValue(t *testing.T) {
	r := Ok().WithValue("hello")

	success, status, value, err := r.Deconstruct()
	if !success {
		t.Errorf("expected success to be true")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value == nil {
		t.Errorf("expected value to have a value, but got nil")
	} else {
		switch value.(type) {
		case string:
			break
		default:
			t.Errorf("expected value to be 'hello', but got: %v", value)
		}
	}

	if err != nil {
		t.Errorf("expected err to be nil, but got: %v", err)
	}
}

func TestFailureWithError(t *testing.T) {
	testErr := errors.New("oops, look like an error occured")
	r := Failure(testErr)

	success, status, value, err := r.Deconstruct()
	if success {
		t.Errorf("expected success to be false")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err == nil {
		t.Errorf("expected err to be '%v', but got nil", err)
	}

	if err.Error() != testErr.Error() {
		t.Errorf("expected error to have a value of '%s' but got '%s'", testErr.Error(), err.Error())
	}
}

func TestFailureWithErrorString(t *testing.T) {
	testErr := "oops, look like an error occured"
	r := Failure(testErr)

	success, status, value, err := r.Deconstruct()
	if success {
		t.Errorf("expected success to be false")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err == nil {
		t.Errorf("expected err to be '%v', but got nil", err)
	}

	if e := err.Error(); e != testErr {
		t.Errorf("expected err to have a value of '%s', but got '%s'", testErr, e)
	}
}

func TestFailureWithErrorInterface(t *testing.T) {
	testErrValue := 500
	var testErr interface{} = testErrValue

	r := Failure(testErr)

	success, status, value, err := r.Deconstruct()
	if success {
		t.Errorf("expected success to be false")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err == nil {
		t.Errorf("expected err to be '%v', but got nil", err)
	}

	if e := err.Error(); e != strconv.Itoa(testErrValue) {
		t.Errorf("expected err to have a value of '%d', but got '%s'", testErrValue, e)
	}
}

func TestFailureWithStatus(t *testing.T) {
	testErr := "oops, an error occured"
	r := Failure(testErr).WithStatusCode(200)

	success, status, value, err := r.Deconstruct()
	if success {
		t.Errorf("expected success to be false")
	}

	if status != 200 {
		t.Errorf("expected status to be 200, but got: %d", status)
	}

	if value != nil {
		t.Errorf("expected value to be nil, but got: %v", value)
	}

	if err == nil {
		t.Errorf("expected err to be '%v', but got nil", err)
	}

	if e := err.Error(); e != testErr {
		t.Errorf("expected err to have a value of '%s', but got '%s'", testErr, e)
	}
}

func TestFailureWithValue(t *testing.T) {
	testErr := "oops, an error occured"
	r := Failure(testErr).WithValue("hello")

	success, status, value, err := r.Deconstruct()
	if success {
		t.Errorf("expected success to be false")
	}

	if status != 0 {
		t.Errorf("expected status to be default int value (0), but got: %d", status)
	}

	if value == nil {
		t.Errorf("expected value to have a value, but got nil")
	} else {
		switch value.(type) {
		case string:
			break
		default:
			t.Errorf("expected value to be 'hello', but got: %v", value)
		}
	}

	if err == nil {
		t.Errorf("expected err to be '%v', but got nil", err)
	}

	if e := err.Error(); e != testErr {
		t.Errorf("expected err to have a value of '%s', but got '%s'", testErr, e)
	}
}
