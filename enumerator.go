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

package client

import "iter"

// See also how xiter was rejected: https://github.com/golang/go/issues/61898#issuecomment-2877009539

// mimics ruby's Enumerator#to_h
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

// returns first such element that satisfies f
// returns (..., false) if nothing holds
func findFirst[T comparable, U any](seq iter.Seq2[T, U], f func(T, U) bool) (T, U, bool) {
	var k T
	var v U
	for k, v = range seq {
		if f(k, v) {
			return k, v, true
		}
	}
	return k, v, false
}
