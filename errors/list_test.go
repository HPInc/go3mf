// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package errors

import (
	"errors"
	"testing"
)

var _ error = new(List)

func TestList_Error_nil(t *testing.T) {
	want := ""
	var multi *List
	if got := multi.Error(); got != want {
		t.Errorf("List.Error() got %v, want %v", got, want)
	}
}

func TestList_Error_single(t *testing.T) {
	want := `1 error occurred:
	* foo
`

	errors := []error{
		errors.New("foo"),
	}

	multi := &List{Errors: errors}
	if got := multi.Error(); got != want {
		t.Errorf("List.Error() got %v, want %v", got, want)
	}
}

func TestList_Error_multiple(t *testing.T) {
	want := `2 errors occurred:
	* foo
	* bar
`

	errors := []error{
		errors.New("foo"),
		errors.New("bar"),
	}

	multi := &List{Errors: errors}
	if got := multi.Error(); got != want {
		t.Errorf("List.Error() got %v, want %v", got, want)
	}
}

func TestAppend_Error(t *testing.T) {
	original := &List{
		Errors: []error{errors.New("foo")},
	}

	result := Append(original, errors.New("bar")).(*List)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	original = &List{}
	result = Append(original, errors.New("bar")).(*List)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test when a typed nil is passed
	var e *List
	result = Append(e, errors.New("baz")).(*List)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}

	// Test flattening
	original = &List{
		Errors: []error{errors.New("foo")},
	}

	result = Append(original, Append(nil, errors.New("foo"), errors.New("bar"))).(*List)
	if len(result.Errors) != 3 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilError(t *testing.T) {
	var err error
	result := Append(err, errors.New("bar")).(*List)
	if len(result.Errors) != 1 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_ErrsEmpty(t *testing.T) {
	var err error
	result := Append(err, []error{}...)
	if result != nil {
		t.Fatalf("wrong err: %v", result)
	}
}

func TestAppend_ErrsWithNil(t *testing.T) {
	var err *List
	result := Append(err, []error{nil}...)
	if result.(*List) != nil {
		t.Fatalf("wrong err: %v", result)
	}
}

func TestAppend_NilErrorArg(t *testing.T) {
	var err error
	var nilErr *List
	result := Append(err, nilErr).(*List)
	if len(result.Errors) != 0 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NilErrorIfaceArg(t *testing.T) {
	var err error
	var nilErr error
	result := Append(err, nilErr)
	if result != nil {
		t.Fatalf("wrong err: %v", result)
	}
}

func TestAppend_NonError(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, errors.New("bar")).(*List)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}

func TestAppend_NonError_Error(t *testing.T) {
	original := errors.New("foo")
	result := Append(original, Append(nil, errors.New("bar"))).(*List)
	if len(result.Errors) != 2 {
		t.Fatalf("wrong len: %d", len(result.Errors))
	}
}
