// Copyright 2025- The sacloud/http-client-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client_test

import (
	"testing"

	. "github.com/sacloud/http-client-go"
	"github.com/stretchr/testify/suite"
)

type ErrorTestSuite struct{ suite.Suite }

func TestErrorTestSuite(t *testing.T) { suite.Run(t, new(ErrorTestSuite)) }

func (s *ErrorTestSuite) TestAPIError_zero() {
	var zero Error

	actual := zero.Error()
	expected := "API Error"

	s.Equal(expected, actual)
}

func (s *ErrorTestSuite) TestAPIError_code() {
	e := NewError(404, "", nil)
	actual := e.Error()
	expected := "API Error 404"

	s.Equal(expected, actual)
}

func (s *ErrorTestSuite) TestAPIError_msg() {
	e := NewError(0, "message", nil)
	actual := e.Error()
	expected := "API Error - message"

	s.Equal(expected, actual)
}

type customErr struct{}

func (e *customErr) Error() string { return "custom error" }

func (s *ErrorTestSuite) TestAPIError_err() {
	e := NewError(0, "", &customErr{})
	actual := e.Error()
	expected := "API Error: custom error"

	s.Equal(expected, actual)
}

func (s *ErrorTestSuite) TestAPIError_combined() {
	e := NewError(500, "internal server error", &customErr{})
	actual := e.Error()
	expected := "API Error 500 - internal server error: custom error"

	s.Equal(expected, actual)
}

func (s *ErrorTestSuite) TestIsNotFoundError_null() {
	s.False(IsNotFoundError(nil))
}

func (s *ErrorTestSuite) TestIsNotFoundError_anotherType() {
	s.False(IsNotFoundError(&customErr{}))
}

func (s *ErrorTestSuite) TestIsNotFoundError_not404() {
	e := NewError(500, "internal server error", nil)
	s.False(IsNotFoundError(e))
}

func (s *ErrorTestSuite) TestIsNotFoundError_404() {
	e := NewError(404, "not found", nil)
	s.True(IsNotFoundError(e))
}
