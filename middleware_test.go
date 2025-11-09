// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestAPIVerificationFailed(t *testing.T) {
	err := &APIVerificationFailed{
		Section:              "consumer",
		MissingSpecification: []string{"application/json", "application/x-yaml"},
		MissingRegistration:  []string{"text/html", "application/xml"},
	}

	expected := `missing [text/html, application/xml] consumer registrations
missing from spec file [application/json, application/x-yaml] consumer`
	assert.Equal(t, expected, err.Error())
}
