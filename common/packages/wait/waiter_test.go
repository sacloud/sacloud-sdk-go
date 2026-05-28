// Copyright 2022-2025 The sacloud/packages-go Authors
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

package wait

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type dummyState struct {
	state interface{}
	err   error
}

func testStateCheckFunc(target interface{}) (bool, error) {
	state, ok := target.(*dummyState)
	if !ok {
		return false, fmt.Errorf("got invalid state type: %+v", target)
	}
	return state.state != nil, state.err
}

func TestStatePollingWaiter_withStateCheckFunc(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		waiter := &PollingWaiter{
			ReadFunc: func() (interface{}, error) {
				return &dummyState{}, nil
			},
			StateCheckFunc: testStateCheckFunc,
			Timeout:        5 * time.Millisecond,
			Interval:       1 * time.Millisecond,
		}
		ctx := context.Background()
		_, err := waiter.WaitForState(ctx)
		require.Error(t, err)
		require.EqualError(t, err, "context deadline exceeded")
	})

	t.Run("parent context was canceled", func(t *testing.T) {
		waiter := &PollingWaiter{
			ReadFunc: func() (interface{}, error) {
				return &dummyState{}, nil
			},
			StateCheckFunc: testStateCheckFunc,
			Timeout:        100 * time.Millisecond,
			Interval:       1 * time.Millisecond,
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		_, err := waiter.WaitForState(ctx)
		go func() {
			time.Sleep(5 * time.Millisecond)
			cancel()
		}()

		require.Error(t, err)
		require.EqualError(t, err, "context deadline exceeded")
	})
}
