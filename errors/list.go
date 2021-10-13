// © Copyright 2021 HP Development Company, L.P.
// SPDX-License Identifier: BSD-2-Clause

package errors

import (
	"fmt"
	"strings"
)

// List is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type List struct {
	Errors []error
}

func (e *List) Len() int {
	if e == nil {
		return 0
	}
	return len(e.Errors)
}

func (e *List) Less(i, j int) bool {
	if e == nil {
		return false
	}
	return e.Errors[i].Error() < e.Errors[j].Error()
}

func (e *List) Swap(i, j int) {
	if e == nil {
		return
	}
	e.Errors[i], e.Errors[j] = e.Errors[j], e.Errors[i]
}

func (e *List) Error() string {
	if e == nil {
		return ""
	}
	return listFormatFunc(e.Errors)
}

// Unwrap returns the first error of the chain if not empty.
func (e *List) Unwrap() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e.Errors[0]
}

// Append is a helper function that will append more errors
// onto an List in order to create a larger multi-error.
//
// If err is not a errors.List, then it will be turned into
// one. If any of the errs are errors.List, they will be flattened
// one level into err.
func Append(err error, errs ...error) error {
	if len(errs) == 0 {
		return err
	}
	switch err := err.(type) {
	case *List:
		// Go through each error and flatten
		for _, e := range errs {
			if e == nil {
				continue
			}
			if err == nil {
				err = new(List)
			}
			switch e := e.(type) {
			case *List:
				if e != nil {
					err.Errors = append(err.Errors, e.Errors...)
				}
			default:
				if e != nil {
					err.Errors = append(err.Errors, e)
				}
			}
		}
		return err
	default:
		var newErrs []error
		if err != nil {
			newErrs = append(newErrs, err)
		}
		for _, e := range errs {
			if e == nil {
				continue
			}
			newErrs = append(newErrs, e)
		}
		if len(newErrs) != 0 {
			return Append(&List{}, newErrs...)
		}
		return nil
	}
}

func listFormatFunc(es []error) string {
	if len(es) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n", es[0])
	}

	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n",
		len(es), strings.Join(points, "\n\t"))
}
