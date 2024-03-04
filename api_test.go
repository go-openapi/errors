// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customError struct {
	apiError
}

func TestServeError(t *testing.T) {
	// method not allowed wins
	// err abides by the Error interface
	err := MethodNotAllowed("GET", []string{"POST", "PUT"})
	recorder := httptest.NewRecorder()
	ServeError(recorder, nil, err)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	assert.Equal(t, "POST,PUT", recorder.Header().Get("Allow"))
	// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
	assert.Equal(t, `{"code":405,"message":"method GET is not allowed, but [POST,PUT] are"}`, recorder.Body.String())

	// renders status code from error when present
	err = NotFound("")
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, err)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
	assert.Equal(t, `{"code":404,"message":"Not found"}`, recorder.Body.String())

	// renders mapped status code from error when present
	err = InvalidTypeName("someType")
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, err)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
	assert.Equal(t, `{"code":601,"message":"someType is an invalid type name"}`, recorder.Body.String())

	// same, but override DefaultHTTPCode
	func() {
		oldDefaultHTTPCode := DefaultHTTPCode
		defer func() { DefaultHTTPCode = oldDefaultHTTPCode }()
		DefaultHTTPCode = http.StatusBadRequest

		err = InvalidTypeName("someType")
		recorder = httptest.NewRecorder()
		ServeError(recorder, nil, err)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
		assert.Equal(t, `{"code":601,"message":"someType is an invalid type name"}`, recorder.Body.String())
	}()

	// defaults to internal server error
	simpleErr := errors.New("some error")
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, simpleErr)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	// assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
	assert.Equal(t, `{"code":500,"message":"some error"}`, recorder.Body.String())

	// composite errors

	// unrecognized: return internal error with first error only - the second error is ignored
	compositeErr := &CompositeError{
		Errors: []error{
			errors.New("firstError"),
			errors.New("anotherError"),
		},
	}
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, compositeErr)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, `{"code":500,"message":"firstError"}`, recorder.Body.String())

	// recognized: return internal error with first error only - the second error is ignored
	compositeErr = &CompositeError{
		Errors: []error{
			New(600, "myApiError"),
			New(601, "myOtherApiError"),
		},
	}
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, compositeErr)
	assert.Equal(t, CompositeErrorCode, recorder.Code)
	assert.Equal(t, `{"code":600,"message":"myApiError"}`, recorder.Body.String())

	// recognized API Error, flattened
	compositeErr = &CompositeError{
		Errors: []error{
			&CompositeError{
				Errors: []error{
					New(600, "myApiError"),
					New(601, "myOtherApiError"),
				},
			},
		},
	}
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, compositeErr)
	assert.Equal(t, CompositeErrorCode, recorder.Code)
	assert.Equal(t, `{"code":600,"message":"myApiError"}`, recorder.Body.String())

	// check guard against empty CompositeError (e.g. nil Error interface)
	compositeErr = &CompositeError{
		Errors: []error{
			&CompositeError{
				Errors: []error{},
			},
		},
	}
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, compositeErr)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, `{"code":500,"message":"Unknown error"}`, recorder.Body.String())

	// check guard against nil type
	recorder = httptest.NewRecorder()
	ServeError(recorder, nil, nil)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, `{"code":500,"message":"Unknown error"}`, recorder.Body.String())

	recorder = httptest.NewRecorder()
	var z *customError
	ServeError(recorder, nil, z)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, `{"code":500,"message":"Unknown error"}`, recorder.Body.String())
}

func TestAPIErrors(t *testing.T) {
	err := New(402, "this failed %s", "yada")
	require.Error(t, err)
	assert.EqualValues(t, 402, err.Code())
	assert.EqualValues(t, "this failed yada", err.Error())

	err = NotFound("this failed %d", 1)
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code())
	assert.EqualValues(t, "this failed 1", err.Error())

	err = NotFound("")
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code())
	assert.EqualValues(t, "Not found", err.Error())

	err = NotImplemented("not implemented")
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotImplemented, err.Code())
	assert.EqualValues(t, "not implemented", err.Error())

	err = MethodNotAllowed("GET", []string{"POST", "PUT"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusMethodNotAllowed, err.Code())
	assert.EqualValues(t, "method GET is not allowed, but [POST,PUT] are", err.Error())

	err = InvalidContentType("application/saml", []string{"application/json", "application/x-yaml"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusUnsupportedMediaType, err.Code())
	assert.EqualValues(t, "unsupported media type \"application/saml\", only [application/json application/x-yaml] are allowed", err.Error())

	err = InvalidResponseFormat("application/saml", []string{"application/json", "application/x-yaml"})
	require.Error(t, err)
	assert.EqualValues(t, http.StatusNotAcceptable, err.Code())
	assert.EqualValues(t, "unsupported media type requested, only [application/json application/x-yaml] are available", err.Error())
}

func TestValidateName(t *testing.T) {
	v := &Validation{Name: "myValidation", message: "myMessage"}

	// unchanged
	vv := v.ValidateName("")
	assert.EqualValues(t, "myValidation", vv.Name)
	assert.EqualValues(t, "myMessage", vv.message)

	// forced
	vv = v.ValidateName("myNewName")
	assert.EqualValues(t, "myNewName.myValidation", vv.Name)
	assert.EqualValues(t, "myNewName.myMessage", vv.message)

	v.Name = ""
	v.message = "myMessage"

	// unchanged
	vv = v.ValidateName("")
	assert.EqualValues(t, "", vv.Name)
	assert.EqualValues(t, "myMessage", vv.message)

	// forced
	vv = v.ValidateName("myNewName")
	assert.EqualValues(t, "myNewName", vv.Name)
	assert.EqualValues(t, "myNewNamemyMessage", vv.message)
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
