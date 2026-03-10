// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitoringsuite

import (
	"strconv"
)

// generic-ish type cast helper function
func intoOpt[T, U any, P interface {
	*T
	Reset()
	SetTo(u U)
}](v *U) (opt T) {
	if v == nil {
		P(&opt).Reset()
		return
	}
	P(&opt).SetTo(*v)
	return opt
}

// generic-ish type cast helper function
func intoNil[T, U any, P interface {
	*T
	SetTo(u U)
	SetToNull()
}](v *U) (opt T) {
	if v == nil {
		P(&opt).SetToNull()
		return
	}
	P(&opt).SetTo(*v)
	return
}

// generic-ish type cast helper function
func intoOptNil[T, U any, P interface {
	*T
	SetTo(u U)
	SetToNull()
	Reset()
}](v *U) T {
	return intoOpt[T, U, P](v)
}

// string parser
func fromStringPtr[
	T any,
	U ~int | ~int8 | ~int16 | ~int32 | ~int64,
	P interface {
		*T
		Reset()
		SetTo(u U)
	},
](v *string) (opt T, err error) {
	var zero U
	var n int
	switch any(&zero).(type) {
	case *int8:
		n = 8
	case *int16:
		n = 16
	case *int32:
		n = 32
	case *int64:
		n = 64
	case *int:
		n = 64 // or ... ?
	default:
		panic("unreachable")
	}

	if v == nil {
		P(&opt).Reset()
		return
	}
	val, err := strconv.ParseInt(*v, 10, n)
	if err != nil {
		P(&opt).Reset()
		return
	}
	P(&opt).SetTo(U(val))
	return
}
