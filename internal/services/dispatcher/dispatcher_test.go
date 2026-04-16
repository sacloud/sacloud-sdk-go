// Copyright 2022-2025 The sacloud/services Authors
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

package dispatcher

import (
	"testing"

	"github.com/sacloud/services"
	"github.com/sacloud/services/dummy"
	"github.com/stretchr/testify/require"
)

func TestDispatch(t *testing.T) {
	type args struct {
		arguments []string
		parameter interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "invalid platform",
			args: args{
				arguments: []string{"invalid", "dummy", "dummy"},
				parameter: nil,
			},
			wantErr: true,
		},
		{
			name: "operation not found",
			args: args{
				arguments: []string{"dummy", "dummy", "not-found"},
				parameter: map[string]interface{}{
					"Field1": "required",
				},
			},
			wantErr: true,
		},
		{
			name: "with results operation",
			args: args{
				arguments: []string{"dummy", "dummy", "find"},
				parameter: map[string]interface{}{
					"Field1": "required",
				},
			},
			want: []*dummy.FindResult{
				{Dummy: "result1"},
				{Dummy: "result2"},
				{Dummy: "result3"},
			},
			wantErr: false,
		},
		{
			name: "with result operation",
			args: args{
				arguments: []string{"dummy", "dummy", "read"},
				parameter: map[string]interface{}{},
			},
			want:    &dummy.ReadResult{Dummy: "result"},
			wantErr: false,
		},
		{
			name: "with error from operation",
			args: args{
				arguments: []string{"dummy", "dummy", "error-read"},
				parameter: map[string]interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "action without error",
			args: args{
				arguments: []string{"dummy", "dummy", "action"},
				parameter: map[string]interface{}{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "action with error",
			args: args{
				arguments: []string{"dummy", "dummy", "error-action"},
				parameter: map[string]interface{}{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parameter",
			args: args{
				arguments: []string{"dummy", "dummy", "echo"},
				parameter: map[string]interface{}{
					"Field1": "value1",
					"Field2": "value2",
				},
			},
			want: &dummy.EchoRequest{
				Field1: "value1",
				Field2: "value2",
			},
			wantErr: false,
		},
	}

	// 準備
	Register("dummy", services.Services{dummy.New()})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Dispatch(tt.args.arguments, tt.args.parameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dispatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				require.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestDispatch_panic(t *testing.T) {
	Register("dummy", services.Services{dummy.New()})

	defer func() {
		panicMessage := recover()
		require.Equal(t, "invalid arguments", panicMessage)
	}()
	Dispatch([]string{}, nil) //nolint
}

func TestRegister_panic(t *testing.T) {
	defer func() {
		panicMessage := recover()
		require.Equal(t, "platformName is required", panicMessage)
	}()
	Register("", services.Services{dummy.New()})
}
