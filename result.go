// Copyright 2025- The sacloud/saclient-go Authors
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

package saclient

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type option[T any] struct {
	some T

	// This is intentionally flipped from Rust's Option<T>
	// to make the zero value mean "not set".
	set bool
}

// This is just Result<Option<T>, E>
type resultOption[T any] struct {
	option[T]
	err error
}

func (m *option[T]) from(src func() (T, bool)) {
	if m != nil {
		value, set := src()
		*m = option[T]{value, set}
	}
}

func (m *option[T]) initialize(v T) { m.from(func() (T, bool) { return v, true }) }

func (m *option[T]) String() string {
	// This is used by flag package
	if m == nil {
		panic("nil dereference")
	}

	if v, ok := m.Get(); ok {
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

func (m *option[T]) Get() (T, bool) {
	var zero T

	if m == nil {
		return zero, false
	} else {
		return m.some, m.set
	}
}

func (m *option[T]) Set(s string) error {
	// This is used by flag package

	switch m := any(m).(type) {
	case *option[string]:
		m.initialize(s)

	case *option[int64]:
		if v, err := strconv.ParseInt(s, 0, 64); err != nil {
			return err
		} else {
			m.initialize(v)
		}

	case *option[[]string]:
		// This behaviour mimics usacloud's --zones= option
		r := strings.NewReader(s)
		c := csv.NewReader(r)
		if v, err := c.Read(); err != nil {
			return err
		} else {
			m.initialize(v)
		}

	default:
		// Should be unreachable because of the constraint,
		// but good practice to guard anyway.
		panic("unsupported type")
	}
	return nil
}

func (m *option[T]) fromEnv(s string) error {
	//nolint:gocritic
	if m == nil {
		return NewErrorf("nil option")
	} else if m.set {
		return nil
	} else {
		return m.Set(s)
	}
}

func (r resultOption[T]) isErr() bool     { return r.err != nil }
func (r resultOption[T]) isNone() bool    { return !r.set }
func (r resultOption[T]) isSome() bool    { return !r.isErr() && !r.isNone() }
func (r resultOption[T]) error() error    { return r.err }
func (r resultOption[T]) ok() *option[T]  { return &r.option }
func (r resultOption[T]) some() (T, bool) { return r.ok().Get() }
func (r resultOption[T]) unwrap() T {
	if ret, ok := r.some(); ok {
		return ret
	} else {
		panic("called unwrap on a None value")
	}
}
func (r resultOption[T]) unwrapOr(zero T) T {
	if ret, ok := r.some(); ok {
		return ret
	} else {
		return zero
	}
}
func (r resultOption[T]) asPtr() *T {
	if r.isSome() {
		return &r.option.some
	} else {
		return nil
	}
}
func (r resultOption[T]) decompose() (v T, ok bool, err error) { return r.option.some, r.set, r.err }

func resultOptionErr[T any](err error) resultOption[T] { return resultOption[T]{err: err} }
func resultOptionNone[T any]() resultOption[T]         { return resultOption[T]{} }
func resultOptionSome[T any](some T) resultOption[T] {
	return resultOption[T]{option: option[T]{some: some, set: true}}
}
