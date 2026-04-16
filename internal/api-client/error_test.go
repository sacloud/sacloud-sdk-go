package client_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	client "github.com/sacloud/api-client-go"
)

type XXXAPIError struct {
	Extra string
	Err   error
}

func (e *XXXAPIError) Error() string {
	return e.Err.Error() + ", extra: " + e.Extra
}

func (e *XXXAPIError) Unwrap() error {
	return e.Err
}

func TestAPIError(t *testing.T) {
	err := &client.APIError{
		Code:    http.StatusNotFound,
		Message: "not found",
		Err:     errors.New("wrapped error"),
	}
	require.Equal(t, "API Error 404 - not found: wrapped error", err.Error())
	assert.True(t, client.IsNotFoundError(err))

	// Unwrap
	assert.Equal(t, "wrapped error", errors.Unwrap(err).Error())

	err2 := client.NewAPIError(http.StatusBadRequest, "", nil)
	require.Equal(t, "API Error 400 - Bad Request", err2.Error())
	assert.False(t, client.IsNotFoundError(err2))
}

func TestXXXAPIError_IsNotFoundError(t *testing.T) {
	// Test with wrapped APIError

	assert.False(t, client.IsNotFoundError(nil))

	xerr := &XXXAPIError{
		Err: &client.APIError{
			Code:    http.StatusNotFound,
			Message: "not found",
			Err:     errors.New("wrapped error"),
		},
		Extra: "extra info",
	}
	assert.True(t, client.IsNotFoundError(xerr))

	xerr2 := &XXXAPIError{
		Err: &client.APIError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("wrapped error 2"),
		},
		Extra: "extra info 2",
	}
	assert.False(t, client.IsNotFoundError(xerr2))
}
