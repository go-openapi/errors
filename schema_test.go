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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx
func TestSchemaErrors(t *testing.T) {
	t.Run("with InvalidType", func(t *testing.T) {
		err := InvalidType("confirmed", "query", "boolean", nil)
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed in query must be of type boolean", err.Error())

		err = InvalidType("confirmed", "", "boolean", nil)
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed must be of type boolean", err.Error())

		err = InvalidType("confirmed", "query", "boolean", "hello")
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed in query must be of type boolean: \"hello\"", err.Error())

		err = InvalidType("confirmed", "query", "boolean", errors.New("hello"))
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed in query must be of type boolean, because: hello", err.Error())

		err = InvalidType("confirmed", "", "boolean", "hello")
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed must be of type boolean: \"hello\"", err.Error())

		err = InvalidType("confirmed", "", "boolean", errors.New("hello"))
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "confirmed must be of type boolean, because: hello", err.Error())
	})

	t.Run("with DuplicateItems", func(t *testing.T) {
		err := DuplicateItems("uniques", "query")
		require.Error(t, err)
		assert.EqualValues(t, UniqueFailCode, err.Code())
		assert.Equal(t, "uniques in query shouldn't contain duplicates", err.Error())

		err = DuplicateItems("uniques", "")
		require.Error(t, err)
		assert.EqualValues(t, UniqueFailCode, err.Code())
		assert.Equal(t, "uniques shouldn't contain duplicates", err.Error())
	})

	t.Run("with TooMany/TooFew Items", func(t *testing.T) {
		err := TooManyItems("something", "query", 5, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxItemsFailCode, err.Code())
		assert.Equal(t, "something in query should have at most 5 items", err.Error())
		assert.Equal(t, 6, err.Value)

		err = TooManyItems("something", "", 5, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxItemsFailCode, err.Code())
		assert.Equal(t, "something should have at most 5 items", err.Error())
		assert.Equal(t, 6, err.Value)

		err = TooFewItems("something", "", 5, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinItemsFailCode, err.Code())
		assert.Equal(t, "something should have at least 5 items", err.Error())
		assert.Equal(t, 4, err.Value)
	})

	t.Run("with ExceedsMaximum", func(t *testing.T) {
		err := ExceedsMaximumInt("something", "query", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumInt("something", "", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumInt("something", "query", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumInt("something", "", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumUint("something", "query", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumUint("something", "", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumUint("something", "query", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximumUint("something", "", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximum("something", "query", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximum("something", "", 5, false, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than or equal to 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximum("something", "query", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something in query should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)

		err = ExceedsMaximum("something", "", 5, true, 6)
		require.Error(t, err)
		assert.EqualValues(t, MaxFailCode, err.Code())
		assert.Equal(t, "something should be less than 5", err.Error())
		assert.Equal(t, 6, err.Value)
	})

	t.Run("with ExceedsMinimum", func(t *testing.T) {
		err := ExceedsMinimumInt("something", "query", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumInt("something", "", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumInt("something", "query", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumInt("something", "", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumUint("something", "query", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumUint("something", "", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumUint("something", "query", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimumUint("something", "", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimum("something", "query", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimum("something", "", 5, false, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than or equal to 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimum("something", "query", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something in query should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = ExceedsMinimum("something", "", 5, true, 4)
		require.Error(t, err)
		assert.EqualValues(t, MinFailCode, err.Code())
		assert.Equal(t, "something should be greater than 5", err.Error())
		assert.Equal(t, 4, err.Value)

		err = NotMultipleOf("something", "query", 5, 1)
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfFailCode, err.Code())
		assert.Equal(t, "something in query should be a multiple of 5", err.Error())
		assert.Equal(t, 1, err.Value)
	})

	t.Run("with MultipleOf", func(t *testing.T) {
		err := NotMultipleOf("something", "query", float64(5), float64(1))
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfFailCode, err.Code())
		assert.Equal(t, "something in query should be a multiple of 5", err.Error())
		assert.InDelta(t, float64(1), err.Value, 1e-6)

		err = NotMultipleOf("something", "query", uint64(5), uint64(1))
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfFailCode, err.Code())
		assert.Equal(t, "something in query should be a multiple of 5", err.Error())
		assert.Equal(t, uint64(1), err.Value)

		err = NotMultipleOf("something", "", 5, 1)
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfFailCode, err.Code())
		assert.Equal(t, "something should be a multiple of 5", err.Error())
		assert.Equal(t, 1, err.Value)

		err = MultipleOfMustBePositive("path", "body", float64(-10))
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfMustBePositiveCode, err.Code())
		assert.Equal(t, `factor MultipleOf declared for path must be positive: -10`, err.Error())
		assert.InDelta(t, float64(-10), err.Value, 1e-6)

		err = MultipleOfMustBePositive("path", "body", int64(-10))
		require.Error(t, err)
		assert.EqualValues(t, MultipleOfMustBePositiveCode, err.Code())
		assert.Equal(t, `factor MultipleOf declared for path must be positive: -10`, err.Error())
		assert.Equal(t, int64(-10), err.Value)
	})

	t.Run("with EnumFail", func(t *testing.T) {
		err := EnumFail("something", "query", "yada", []interface{}{"hello", "world"})
		require.Error(t, err)
		assert.EqualValues(t, EnumFailCode, err.Code())
		assert.Equal(t, "something in query should be one of [hello world]", err.Error())
		assert.Equal(t, "yada", err.Value)

		err = EnumFail("something", "", "yada", []interface{}{"hello", "world"})
		require.Error(t, err)
		assert.EqualValues(t, EnumFailCode, err.Code())
		assert.Equal(t, "something should be one of [hello world]", err.Error())
		assert.Equal(t, "yada", err.Value)
	})

	t.Run("with Required", func(t *testing.T) {
		err := Required("something", "query", nil)
		require.Error(t, err)
		assert.EqualValues(t, RequiredFailCode, err.Code())
		assert.Equal(t, "something in query is required", err.Error())
		assert.Nil(t, err.Value)

		err = Required("something", "", nil)
		require.Error(t, err)
		assert.EqualValues(t, RequiredFailCode, err.Code())
		assert.Equal(t, "something is required", err.Error())
		assert.Nil(t, err.Value)
	})

	t.Run("with ReadOnly", func(t *testing.T) {
		err := ReadOnly("something", "query", nil)
		require.Error(t, err)
		assert.EqualValues(t, ReadOnlyFailCode, err.Code())
		assert.Equal(t, "something in query is readOnly", err.Error())
		assert.Nil(t, err.Value)

		err = ReadOnly("something", "", nil)
		require.Error(t, err)
		assert.EqualValues(t, ReadOnlyFailCode, err.Code())
		assert.Equal(t, "something is readOnly", err.Error())
		assert.Nil(t, err.Value)
	})

	t.Run("with TooLong/TooShort", func(t *testing.T) {
		err := TooLong("something", "query", 5, "abcdef")
		require.Error(t, err)
		assert.EqualValues(t, TooLongFailCode, err.Code())
		assert.Equal(t, "something in query should be at most 5 chars long", err.Error())
		assert.Equal(t, "abcdef", err.Value)

		err = TooLong("something", "", 5, "abcdef")
		require.Error(t, err)
		assert.EqualValues(t, TooLongFailCode, err.Code())
		assert.Equal(t, "something should be at most 5 chars long", err.Error())
		assert.Equal(t, "abcdef", err.Value)

		err = TooShort("something", "query", 5, "a")
		require.Error(t, err)
		assert.EqualValues(t, TooShortFailCode, err.Code())
		assert.Equal(t, "something in query should be at least 5 chars long", err.Error())
		assert.Equal(t, "a", err.Value)

		err = TooShort("something", "", 5, "a")
		require.Error(t, err)
		assert.EqualValues(t, TooShortFailCode, err.Code())
		assert.Equal(t, "something should be at least 5 chars long", err.Error())
		assert.Equal(t, "a", err.Value)
	})

	t.Run("with FailedPattern", func(t *testing.T) {
		err := FailedPattern("something", "query", "\\d+", "a")
		require.Error(t, err)
		assert.EqualValues(t, PatternFailCode, err.Code())
		assert.Equal(t, "something in query should match '\\d+'", err.Error())
		assert.Equal(t, "a", err.Value)

		err = FailedPattern("something", "", "\\d+", "a")
		require.Error(t, err)
		assert.EqualValues(t, PatternFailCode, err.Code())
		assert.Equal(t, "something should match '\\d+'", err.Error())
		assert.Equal(t, "a", err.Value)
	})

	t.Run("with InvalidType", func(t *testing.T) {
		err := InvalidTypeName("something")
		require.Error(t, err)
		assert.EqualValues(t, InvalidTypeCode, err.Code())
		assert.Equal(t, "something is an invalid type name", err.Error())
	})

	t.Run("with AdditionalItemsNotAllowed", func(t *testing.T) {
		err := AdditionalItemsNotAllowed("something", "query")
		require.Error(t, err)
		assert.EqualValues(t, NoAdditionalItemsCode, err.Code())
		assert.Equal(t, "something in query can't have additional items", err.Error())

		err = AdditionalItemsNotAllowed("something", "")
		require.Error(t, err)
		assert.EqualValues(t, NoAdditionalItemsCode, err.Code())
		assert.Equal(t, "something can't have additional items", err.Error())
	})

	err := InvalidCollectionFormat("something", "query", "yada")
	require.Error(t, err)
	assert.EqualValues(t, InvalidTypeCode, err.Code())
	assert.Equal(t, "the collection format \"yada\" is not supported for the query param \"something\"", err.Error())

	t.Run("with CompositeValidationError", func(t *testing.T) {
		err := CompositeValidationError()
		require.Error(t, err)
		assert.EqualValues(t, CompositeErrorCode, err.Code())
		assert.Equal(t, "validation failure list", err.Error())

		testErr1 := errors.New("first error")
		testErr2 := errors.New("second error")
		err = CompositeValidationError(testErr1, testErr2)
		require.Error(t, err)
		assert.EqualValues(t, CompositeErrorCode, err.Code())
		assert.Equal(t, "validation failure list:\nfirst error\nsecond error", err.Error())

		require.ErrorIs(t, err, testErr1)
		require.ErrorIs(t, err, testErr2)
	})

	t.Run("should set validation name in CompositeValidation error", func(t *testing.T) {
		err := CompositeValidationError(
			InvalidContentType("text/html", []string{"application/json"}),
			CompositeValidationError(
				InvalidTypeName("y"),
			),
		)
		_ = err.ValidateName("new-name")
		const expectedMessage = `validation failure list:
new-name.unsupported media type "text/html", only [application/json] are allowed
validation failure list:
new-namey is an invalid type name`
		assert.Equal(t, expectedMessage, err.Error())
	})

	t.Run("with PropertyNotAllowed", func(t *testing.T) {
		err = PropertyNotAllowed("path", "body", "key")
		require.Error(t, err)
		assert.EqualValues(t, UnallowedPropertyCode, err.Code())
		// unallowedProperty         = "%s.%s in %s is a forbidden property"
		assert.Equal(t, "path.key in body is a forbidden property", err.Error())

		err = PropertyNotAllowed("path", "", "key")
		require.Error(t, err)
		assert.EqualValues(t, UnallowedPropertyCode, err.Code())
		// unallowedPropertyNoIn     = "%s.%s is a forbidden property"
		assert.Equal(t, "path.key is a forbidden property", err.Error())
	})

	t.Run("with TooMany/TooFew properties", func(t *testing.T) {
		err := TooManyProperties("path", "body", 10)
		require.Error(t, err)
		assert.EqualValues(t, TooManyPropertiesCode, err.Code())
		// tooManyProperties         = "%s in %s should have at most %d properties"
		assert.Equal(t, "path in body should have at most 10 properties", err.Error())

		err = TooManyProperties("path", "", 10)
		require.Error(t, err)
		assert.EqualValues(t, TooManyPropertiesCode, err.Code())
		// tooManyPropertiesNoIn     = "%s should have at most %d properties"
		assert.Equal(t, "path should have at most 10 properties", err.Error())

		err = TooFewProperties("path", "body", 10)
		require.Error(t, err)
		assert.EqualValues(t, TooFewPropertiesCode, err.Code())
		// tooFewProperties          = "%s in %s should have at least %d properties"
		assert.Equal(t, "path in body should have at least 10 properties", err.Error())

		err = TooFewProperties("path", "", 10)
		require.Error(t, err)
		assert.EqualValues(t, TooFewPropertiesCode, err.Code())
		// tooFewPropertiesNoIn      = "%s should have at least %d properties"
		assert.Equal(t, "path should have at least 10 properties", err.Error())
	})

	t.Run("with PatternProperties", func(t *testing.T) {
		err := FailedAllPatternProperties("path", "body", "key")
		require.Error(t, err)
		assert.EqualValues(t, FailedAllPatternPropsCode, err.Code())
		// failedAllPatternProps     = "%s.%s in %s failed all pattern properties"
		assert.Equal(t, "path.key in body failed all pattern properties", err.Error())

		err = FailedAllPatternProperties("path", "", "key")
		require.Error(t, err)
		assert.EqualValues(t, FailedAllPatternPropsCode, err.Code())
		// failedAllPatternPropsNoIn = "%s.%s failed all pattern properties"
		assert.Equal(t, "path.key failed all pattern properties", err.Error())
	})
}
