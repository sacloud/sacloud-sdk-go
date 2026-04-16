package eventbus

import (
	"errors"
	"testing"

	"github.com/sacloud/saclient-go"
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
