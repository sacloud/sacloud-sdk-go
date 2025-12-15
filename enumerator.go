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

import "iter"

// See also how xiter was rejected: https://github.com/golang/go/issues/61898#issuecomment-2877009539

// a Seq that yields 0 times
func nonceSeq[T any]() iter.Seq[T] { return func(yield func(T) bool) {} }

// `(mapcar f seq)`
func intoSeq2[T comparable](seq iter.Seq[T], f func(T) (T, T, bool)) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		for v := range seq {
			if k, w, ok := f(v); !ok {
				continue
			} else if !yield(k, w) {
				return
			}
		}
	}
}

// `(find-if f seq)`, returns (..., false) if nothing holds
func findFirst[
	T comparable,
	U any,
](
	seq iter.Seq2[T, U],
	f func(T, U) bool,
) (
	k T,
	v U,
	ok bool,
) {
	for k, v = range seq {
		if ok = f(k, v); ok {
			return
		}
	}
	return
}

// variant of `findFirst` that only cares about error
func findFirstError[T any](seq iter.Seq[T], f func(T) error) error {
	for v := range seq {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

// transforms Seq2 by f.
// `4` is the number of type parameters
func transformSeq4[
	Q comparable,
	W any,
	E comparable,
	R any,
](
	seq iter.Seq2[Q, W],
	f func(Q, W) (E, R, bool),
) iter.Seq2[E, R] {
	return func(yield func(E, R) bool) {
		for k, v := range seq {
			if nk, nv, ok := f(k, v); !ok {
				continue
			} else if !yield(nk, nv) {
				return
			}
		}
	}
}

// a variant of transformSeq4 that keeps keys as-is.
// `3` is the number of type parameters
func transformSeq3[K comparable, V, W any](seq iter.Seq2[K, V], f func(K, V) W) iter.Seq2[K, W] {
	return transformSeq4(seq, func(k K, v V) (K, W, bool) { return k, f(k, v), true })
}

// `transformSeq2` can be thinkable of course,
// but `transformSeq3[K,V,V](seq, f)` seems just enough.

// `(remove-if f seq)`
func rejectSeq2[T comparable, U any](seq iter.Seq2[T, U], f func(T, U) bool) iter.Seq2[T, U] {
	return transformSeq4(seq, func(k T, v U) (T, U, bool) { return k, v, !f(k, v) })
}

// transforms Seq by f.
func transformSeq[T, U any](seq iter.Seq[T], f func(T) (U, bool)) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if w, ok := f(v); !ok {
				continue
			} else if !yield(w) {
				return
			}
		}
	}
}

// (intuitive)
func selectSeq[T any](seq iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return transformSeq(seq, func(v T) (T, bool) { return v, f(v) })
}

// (intuitive)
func mapSeq[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return transformSeq(seq, func(v T) (U, bool) { return f(v), true })
}

// `(mapcar #cdr seq)`
func valuesOfSeq2[K comparable, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				return
			}
		}
	}
}
