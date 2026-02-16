// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package common

// generic-ish type cast helper function
func IntoOpt[
	T, U any,
	P interface {
		*T
		Reset()
		SetTo(u U)
	},
](v *U) (ret T) {
	if v == nil {
		P(&ret).Reset()
	} else {
		P(&ret).SetTo(*v)
	}

	return
}

// generic-ish type cast helper function
func FromOpt[T any, P interface{ Get() (T, bool) }](v P) (ret *T) {
	val, ok := v.Get()
	if ok {
		ret = &val
	}
	return
}

// generic-ish type cast helper function
func IntoNullable[
	T, U any,
	P interface {
		*T
		SetToNull()
		SetTo(u U)
	},
](
	v *U,
) (
	ret T,
) {
	if v == nil {
		P(&ret).SetToNull()
	} else {
		P(&ret).SetTo(*v)
	}

	return
}

func MapSlice[T, U any](src []T, f func(T) U) (ret []U) {
	ret = make([]U, len(src))
	for i, v := range src {
		ret[i] = f(v)
	}
	return
}

func ConvertFrom[
	T, U any,
	P interface {
		*U
		From(*T)
	},
]() func(T) U {
	return func(t T) (u U) {
		P(&u).From(&t)
		return
	}
}
