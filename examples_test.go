// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package errors_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/errors"
)

func ExampleNew() {
	// Create a generic API error with custom code
	err := errors.New(400, "invalid input: %s", "email")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("code: %d\n", err.Code())

	// Create common HTTP errors
	notFound := errors.NotFound("user %s not found", "john-doe")
	fmt.Printf("not found: %v\n", notFound)
	fmt.Printf("not found code: %d\n", notFound.Code())

	notImpl := errors.NotImplemented("feature: dark mode")
	fmt.Printf("not implemented: %v\n", notImpl)

	// Output:
	// error: invalid input: email
	// code: 400
	// not found: user john-doe not found
	// not found code: 404
	// not implemented: feature: dark mode
}

func ExampleServeError() {
	// Create a simple validation error
	err := errors.Required("email", "body", nil)

	// Simulate HTTP response
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/users", nil)

	// Serve the error as JSON
	errors.ServeError(recorder, request, err)

	fmt.Printf("status: %d\n", recorder.Code)
	fmt.Printf("content-type: %s\n", recorder.Header().Get("Content-Type"))

	// Parse and display the JSON response
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err == nil {
		fmt.Printf("error code: %.0f\n", response["code"])
		fmt.Printf("error message: %s\n", response["message"])
	}

	// Output:
	// status: 422
	// content-type: application/json
	// error code: 602
	// error message: email in body is required
}

func ExampleCompositeValidationError() {
	var errs []error

	// Collect multiple validation errors
	errs = append(errs, errors.Required("name", "body", nil))
	errs = append(errs, errors.TooShort("description", "body", 10, "short"))
	errs = append(errs, errors.InvalidType("age", "body", "integer", "abc"))

	// Combine them into a composite error
	compositeErr := errors.CompositeValidationError(errs...)

	fmt.Printf("error count: %d\n", len(errs))
	fmt.Printf("composite error: %v\n", compositeErr)
	fmt.Printf("code: %d\n", compositeErr.Code())

	// Can unwrap to access individual errors
	if unwrapped := compositeErr.Unwrap(); unwrapped != nil {
		fmt.Printf("unwrapped count: %d\n", len(unwrapped))
	}

	// Output:
	// error count: 3
	// composite error: validation failure list:
	// name in body is required
	// description in body should be at least 10 chars long
	// age in body must be of type integer: "abc"
	// code: 422
	// unwrapped count: 3
}
