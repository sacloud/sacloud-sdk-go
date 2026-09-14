// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package into

import "strconv"

// Generic-ish type cast helper function
//
// This function helps conversion from a pointer of type X into an OptX.
//
// ```golang
//
//	opt := into.Opt[OptInt](new(10))
//
// ```
func Opt[T, U any, P interface {
	*T
	Reset()
	SetTo(u U)
}](v *U) (opt T) {
	if v == nil {
		P(&opt).Reset()
	} else {
		P(&opt).SetTo(*v)
	}
	return
}

// Generic-ish type cast helper function
//
// This function helps conversion from a pointer of type X into a NilX.
//
// ```golang
//
//	opt := into.Nil[NilInt](new(10))
//
// ```
func Nil[T, U any, P interface {
	*T
	SetTo(u U)
	SetToNull()
}](v *U) (opt T) {
	if v == nil {
		P(&opt).SetToNull()
	} else {
		P(&opt).SetTo(*v)
	}
	return
}

// Generic-ish type cast helper function
//
// This function helps conversion from a pointer of type X into an OptNilX.
//
// ```golang
//
//	opt := into.OptNil[OptNilInt](new(10))
//
// ```
func OptNil[T, U any, P interface {
	*T
	SetTo(u U)
	SetToNull()
	Reset()
}](v *U) T {
	return Opt[T, U, P](v)
}

// string parser
//
// ```golang
//
//	opt, err := into.FromStringPtr[OptInt64, int64](new("10"))
//
// ```
func FromStringPtr[
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
	} else if val, err := strconv.ParseInt(*v, 10, n); err != nil {
		P(&opt).Reset()
	} else {
		P(&opt).SetTo(U(val))
	}
	return
}
