// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package errors

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestUnauthenticated(t *testing.T) {
	err := Unauthenticated("basic")
	assert.EqualValues(t, 401, err.Code())
	assert.Equal(t, "unauthenticated for basic", err.Error())
}
