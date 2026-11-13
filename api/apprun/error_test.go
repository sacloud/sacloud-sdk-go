// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

package apprun

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

func TestError_Error(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with msg and err",
			err:  &Error{msg: "something failed", err: baseErr},
			want: "apprun_shared: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "apprun_shared: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "apprun_shared: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	baseErr := errors.New("base error")

	err := NewError("msg", baseErr)
	if err.msg != "msg" || err.err != baseErr {
		t.Errorf("NewError() did not set fields correctly")
	}

	err2 := NewError("msg only", nil)
	if err2.msg != "msg only" || err2.err != nil {
		t.Errorf("NewError() with nil err did not set fields correctly")
	}
}

func TestNewAPIError(t *testing.T) {
	baseErr := errors.New("base error")

	err := NewAPIError("msg", 404, baseErr)
	if !saclient.IsNotFoundError(err) {
		t.Errorf("IsNotFoundError is true for NewAPIError with 404")
	}

	err2 := NewAPIError("msg", 503, nil)
	if saclient.IsNotFoundError(err2) {
		t.Errorf("IsNotFoundError is false for NewAPIError with 503")
	}
}

func TestAPIErrorFromModel(t *testing.T) {
	defaultError := v1.NewModelDefaultErrorReadApplicationNotFound(v1.ModelDefaultError{
		Error: v1.ModelDefaultErrorError{Message: "application not found"},
	})
	cloudError := v1.NewModelCloudctrlErrorCreateApplicationConflict(v1.ModelCloudctrlError{
		ErrorMsg: "application already exists",
	})
	emptyDefault := v1.NewModelDefaultErrorReadApplicationNotFound(v1.ModelDefaultError{
		Error: v1.ModelDefaultErrorError{Code: 42},
	})
	emptyCloud := v1.NewModelCloudctrlErrorCreateApplicationConflict(v1.ModelCloudctrlError{
		ErrorCode: "E42",
	})
	unknownModel := v1.ReadApplicationNotFound{Type: "future-model"}
	var nilModel *v1.ReadApplicationNotFound

	tests := []struct {
		name     string
		code     int
		response modelErrorGetter
		want     []string
		notFound bool
		conflict bool
	}{
		{"default model", http.StatusNotFound, &defaultError, []string{"application not found"}, true, false},
		{"cloud model", http.StatusConflict, &cloudError, []string{"application already exists"}, false, true},
		{"empty default message", http.StatusNotFound, &emptyDefault, []string{"empty message", "42"}, true, false},
		{"empty cloud message", http.StatusConflict, &emptyCloud, []string{"empty message", "E42"}, false, true},
		{"unknown model", http.StatusNotFound, &unknownModel, []string{"unknown API error model", "future-model"}, true, false},
		{"nil response", 0, nil, []string{"nil API response"}, false, false},
		{"typed nil response", http.StatusNotFound, nilModel, []string{"nil API response", "*v1.ReadApplicationNotFound"}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiErrorFromModel("test", tt.code, tt.response)
			for _, want := range tt.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err, want)
				}
			}

			if got := saclient.IsNotFoundError(err); got != tt.notFound {
				t.Errorf("IsNotFoundError = %v, want %v", got, tt.notFound)
			}
			if got := saclient.IsConflictError(err); got != tt.conflict {
				t.Errorf("IsConflictError = %v, want %v", got, tt.conflict)
			}
		})
	}
}

func TestUnknownAPIError(t *testing.T) {
	model := v1.NewModelDefaultErrorReadApplicationNotFound(v1.ModelDefaultError{
		Error: v1.ModelDefaultErrorError{Message: "application not found"},
	})
	var nilModel *v1.ReadApplicationNotFound
	var nilResponse *v1.HandlerListApplications

	tests := []struct {
		name     string
		response any
		want     []string
	}{
		{"model", &model, []string{"application not found"}},
		{"nil", nil, []string{"nil API response"}},
		{"typed nil model", nilModel, []string{"nil API response", "*v1.ReadApplicationNotFound"}},
		{"typed nil response", nilResponse, []string{"nil API response", "*v1.HandlerListApplications"}},
		{"unknown response", &v1.HandlerListApplications{}, []string{"unknown API response type", "*v1.HandlerListApplications"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unknownAPIError("test", tt.response)
			for _, want := range tt.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err, want)
				}
			}
			if saclient.IsNotFoundError(err) || saclient.IsConflictError(err) {
				t.Errorf("unexpected HTTP status classification: %v", err)
			}
		})
	}
}
