// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

//nolint:err113
package errors

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type customError struct {
	apiError
}

func TestServeError(t *testing.T) {
	t.Run("method not allowed wins", func(t *testing.T) {
		// err abides by the Error interface
		err := MethodNotAllowed("GET", []string{"POST", "PUT"})
		require.Error(t, err)

		recorder := httptest.NewRecorder()
		ServeError(recorder, nil, err)
		assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
		assert.Equal(t, "POST,PUT", recorder.Header().Get("Allow"))
		// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
		assert.JSONEq(t,
			`{"code":405,"message":"method GET is not allowed, but [POST,PUT] are"}`,
			recorder.Body.String(),
		)
	})

	t.Run("renders status code from error", func(t *testing.T) {
		err := NotFound("")
		require.Error(t, err)

		recorder := httptest.NewRecorder()
		ServeError(recorder, nil, err)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
		assert.JSONEq(t,
			`{"code":404,"message":"Not found"}`,
			recorder.Body.String(),
		)
	})

	t.Run("renders mapped status code from error", func(t *testing.T) {
		// renders mapped status code from error when present
		err := InvalidTypeName("someType")
		require.Error(t, err)

		recorder := httptest.NewRecorder()
		ServeError(recorder, nil, err)
		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
		assert.JSONEq(t,
			`{"code":601,"message":"someType is an invalid type name"}`,
			recorder.Body.String(),
		)
	})

	t.Run("overrides DefaultHTTPCode", func(t *testing.T) {
		func() {
			oldDefaultHTTPCode := DefaultHTTPCode
			defer func() { DefaultHTTPCode = oldDefaultHTTPCode }()
			DefaultHTTPCode = http.StatusBadRequest

			err := InvalidTypeName("someType")
			require.Error(t, err)

			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, err)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
			assert.JSONEq(t,
				`{"code":601,"message":"someType is an invalid type name"}`,
				recorder.Body.String(),
			)
		}()
	})

	t.Run("defaults to internal server error", func(t *testing.T) {
		simpleErr := errors.New("some error")
		recorder := httptest.NewRecorder()
		ServeError(recorder, nil, simpleErr)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
		assert.JSONEq(t,
			`{"code":500,"message":"some error"}`,
			recorder.Body.String(),
		)
	})

	t.Run("with composite erors", func(t *testing.T) {
		t.Run("unrecognized - return internal error with first error only - the second error is ignored", func(t *testing.T) {
			compositeErr := &CompositeError{
				Errors: []error{
					errors.New("firstError"),
					errors.New("anotherError"),
				},
			}
			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, compositeErr)
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.JSONEq(t,
				`{"code":500,"message":"firstError"}`,
				recorder.Body.String(),
			)
		})

		t.Run("recognized - return internal error with first error only - the second error is ignored", func(t *testing.T) {
			compositeErr := &CompositeError{
				Errors: []error{
					New(600, "myApiError"),
					New(601, "myOtherApiError"),
				},
			}
			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, compositeErr)
			assert.Equal(t, CompositeErrorCode, recorder.Code)
			assert.JSONEq(t,
				`{"code":600,"message":"myApiError"}`,
				recorder.Body.String(),
			)
		})

		t.Run("recognized API Error, flattened", func(t *testing.T) {
			compositeErr := &CompositeError{
				Errors: []error{
					&CompositeError{
						Errors: []error{
							New(600, "myApiError"),
							New(601, "myOtherApiError"),
						},
					},
				},
			}
			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, compositeErr)
			assert.Equal(t, CompositeErrorCode, recorder.Code)
			assert.JSONEq(t,
				`{"code":600,"message":"myApiError"}`,
				recorder.Body.String(),
			)
		})

		// (e.g. nil Error interface)
		t.Run("check guard against empty CompositeError", func(t *testing.T) {
			compositeErr := &CompositeError{
				Errors: []error{
					&CompositeError{
						Errors: []error{},
					},
				},
			}
			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, compositeErr)
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.JSONEq(t,
				`{"code":500,"message":"Unknown error"}`,
				recorder.Body.String(),
			)
		})

		t.Run("check guard against nil type", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ServeError(recorder, nil, nil)
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.JSONEq(t,
				`{"code":500,"message":"Unknown error"}`,
				recorder.Body.String(),
			)
		})

		t.Run("check guard against nil value", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			var z *customError
			ServeError(recorder, nil, z)
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			assert.JSONEq(t,
				`{"code":500,"message":"Unknown error"}`,
				recorder.Body.String(),
			)
		})
	})
}

func TestAPIErrors(t *testing.T) {
	err := New(402, "this failed %s", "yada")
	require.Error(t, err)
	assert.EqualValues(t, 402, err.Code())
	assert.Equal(t, "this failed yada", err.Error())

	err = NotFound("this failed %d", 1)
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code())
	assert.Equal(t, "this failed 1", err.Error())

	err = NotFound("")
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code())
	assert.Equal(t, "Not found", err.Error())

	err = NotImplemented("not implemented")
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotImplemented, err.Code())
	assert.Equal(t, "not implemented", err.Error())

	err = MethodNotAllowed("GET", []string{"POST", "PUT"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, err.Code())
	assert.Equal(t, "method GET is not allowed, but [POST,PUT] are", err.Error())

	err = InvalidContentType("application/saml", []string{"application/json", "application/x-yaml"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusUnsupportedMediaType, err.Code())
	assert.Equal(t, "unsupported media type \"application/saml\", only [application/json application/x-yaml] are allowed", err.Error())

	err = InvalidResponseFormat("application/saml", []string{"application/json", "application/x-yaml"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotAcceptable, err.Code())
	assert.Equal(t, "unsupported media type requested, only [application/json application/x-yaml] are available", err.Error())
}

func TestValidateName(t *testing.T) {
	v := &Validation{Name: "myValidation", message: "myMessage"}

	// unchanged
	vv := v.ValidateName("")
	assert.Equal(t, "myValidation", vv.Name)
	assert.Equal(t, "myMessage", vv.message)

	// forced
	vv = v.ValidateName("myNewName")
	assert.Equal(t, "myNewName.myValidation", vv.Name)
	assert.Equal(t, "myNewName.myMessage", vv.message)

	v.Name = ""
	v.message = "myMessage"

	// unchanged
	vv = v.ValidateName("")
	assert.Empty(t, vv.Name)
	assert.Equal(t, "myMessage", vv.message)

	// forced
	vv = v.ValidateName("myNewName")
	assert.Equal(t, "myNewName", vv.Name)
	assert.Equal(t, "myNewNamemyMessage", vv.message)
}

func TestMarshalJSON(t *testing.T) {
	const (
		expectedCode = http.StatusUnsupportedMediaType
		value        = "myValue"
	)
	list := []string{"a", "b"}

	e := InvalidContentType(value, list)

	jazon, err := e.MarshalJSON()
	require.NoError(t, err)

	expectedMessage := strings.ReplaceAll(fmt.Sprintf(contentTypeFail, value, list), `"`, `\"`)

	expectedJSON := fmt.Sprintf(
		`{"code":%d,"message":"%s","name":"Content-Type","in":"header","value":"%s","values":["a","b"]}`,
		expectedCode, expectedMessage, value,
	)
	assert.JSONEq(t, expectedJSON, string(jazon))

	a := apiError{code: 1, message: "a"}
	jazon, err = a.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":1,"message":"a"}`, string(jazon))

	m := MethodNotAllowedError{code: 1, message: "a", Allowed: []string{"POST"}}
	jazon, err = m.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":1,"message":"a","allowed":["POST"]}`, string(jazon))

	c := CompositeError{Errors: []error{e}, code: 1, message: "a"}
	jazon, err = c.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, fmt.Sprintf(`{"code":1,"message":"a","errors":[%s]}`, expectedJSON), string(jazon))

	p := ParseError{code: 1, message: "x", Name: "a", In: "b", Value: "c", Reason: errors.New("d")}
	jazon, err = p.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":1,"message":"x","name":"a","in":"b","value":"c","reason":"d"}`, string(jazon))
}
