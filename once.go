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

import "sync"

// type-generic variant of sync.Once
type once[T any] struct {
	o sync.Once
	d bool
	t T
	e error
}

func (o *once[T]) Do(f func(*T) error) (*T, error) {
	o.o.Do(func() {
		o.e = f(&o.t)
		o.d = true
	})
	return &o.t, o.e
}

func (o *once[T]) Done() bool { return o.d }
