package eventbus

import (
	"errors"
	"strings"
	"testing"

	v1 "github.com/sacloud/sacloud-sdk-go/api/eventbus/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

func TestUnknownError(t *testing.T) {
	tests := []struct {
		name string
		res  any
		code int
		want string
	}{
		{
			name: "error message",
			res: &v1.ListCommonServiceItemsBadRequest{
				ErrorMsg: v1.NewOptString("bad request"),
				Status:   v1.NewOptString("400 Bad Request"),
			},
			code: 400,
			want: "bad request",
		},
		{
			name: "missing message",
			res:  (*v1.ListCommonServiceItemsBadRequest)(nil),
			want: "nil API error response",
		},
		{
			name: "unknown response",
			res:  123,
			want: "unknown API response type 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := unknownError(tt.res)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("unknownError() error = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

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
			want: "eventbus: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "eventbus: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "eventbus: base error",
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
