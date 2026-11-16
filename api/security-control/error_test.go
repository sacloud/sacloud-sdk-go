// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/sacloud/sacloud-sdk-go/api/security-control/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/require"
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
			want: "security-control: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "security-control: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "security-control: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.New(t).Equal(tt.want, tt.err.Error())
		})
	}
}

func TestNewError(t *testing.T) {
	assert := require.New(t)
	baseErr := errors.New("base error")

	err := NewError("msg", baseErr)
	assert.Equal("msg", err.msg)
	assert.Equal(baseErr, err.err)

	err2 := NewError("msg only", nil)
	assert.Equal("msg only", err2.msg)
	assert.Nil(err2.err)
}

func TestNewAPIError(t *testing.T) {
	assert := require.New(t)
	baseErr := errors.New("base error")

	err := NewAPIError("msg", 404, baseErr)
	assert.Equal("msg", err.msg)
	assert.True(saclient.IsNotFoundError(err))

	err2 := NewAPIError("msg", 503, nil)
	assert.Equal("msg", err2.msg)
	assert.False(saclient.IsNotFoundError(err2))
}

func TestResponseError(t *testing.T) {
	tests := []struct {
		name     string
		response any
		want     string
		notFound bool
	}{
		{
			name: "problem json",
			response: (*v1.EvaluationRulesReadApplicationProblemJSONNotFound)(&v1.NotFound{
				Status: v1.NewOptNotFoundStatus(v1.NotFoundStatus404),
				Detail: v1.NewOptNotFoundDetail(v1.NotFoundDetailNotFound),
			}),
			want:     "API Error 404: not_found",
			notFound: true,
		},
		{
			name: "missing detail",
			response: (*v1.EvaluationRulesReadApplicationProblemJSONNotFound)(&v1.NotFound{
				Status: v1.NewOptNotFoundStatus(v1.NotFoundStatus404),
			}),
			want:     "API Error 404: API error response missing detail",
			notFound: true,
		},
		{
			name: "missing status",
			response: (*v1.EvaluationRulesReadApplicationProblemJSONNotFound)(&v1.NotFound{
				Detail: v1.NewOptNotFoundDetail(v1.NotFoundDetailNotFound),
			}),
			want:     "API Error 404: not_found",
			notFound: true,
		},
		{
			name:     "empty detail",
			response: &v1.UnexpectedErrorStatusCode{StatusCode: 503, Response: v1.UnexpectedError{Detail: v1.NewOptString("")}},
			want:     "API Error 503: API error response has empty detail",
		},
		{
			name:     "nil typed response",
			response: (*v1.EvaluationRulesReadApplicationProblemJSONNotFound)(nil),
			want:     "nil API error response *v1.EvaluationRulesReadApplicationProblemJSONNotFound",
		},
		{
			name:     "nil unexpected response",
			response: (*v1.UnexpectedErrorStatusCode)(nil),
			want:     "nil API error response *v1.UnexpectedErrorStatusCode",
		},
		{
			name:     "unknown response",
			response: struct{ Secret string }{Secret: "do not expose"},
			want:     "unknown API response type struct { Secret string }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := responseError("test", tt.response)
			require.ErrorContains(t, err, tt.want)
			require.Equal(t, tt.notFound, saclient.IsNotFoundError(err))
			require.NotContains(t, err.Error(), "do not expose")
		})
	}
}

func TestOperationErrorResponses(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		status      int
		body        string
		call        func(*v1.Client) error
		want        string
		notFound    bool
		conflict    bool
	}{
		{
			name: "existing json error", contentType: "application/json", status: 404,
			body: `{"status":404,"detail":"not_found"}`,
			call: func(c *v1.Client) error {
				_, err := NewEvaluationRulesOp(c).Read(t.Context(), "test")
				return err
			},
			want: "API Error 404: not_found", notFound: true,
		},
		{
			name: "problem json not found", contentType: "application/problem+json", status: 404,
			body: `{"status":404,"detail":"not_found"}`,
			call: func(c *v1.Client) error {
				_, err := NewEvaluationRulesOp(c).Read(t.Context(), "test")
				return err
			},
			want: "API Error 404: not_found", notFound: true,
		},
		{
			name: "problem json missing status", contentType: "application/problem+json", status: 404,
			body: `{"detail":"not_found"}`,
			call: func(c *v1.Client) error {
				_, err := NewEvaluationRulesOp(c).Read(t.Context(), "test")
				return err
			},
			want: "API Error 404: not_found", notFound: true,
		},
		{
			name: "problem json conflict", contentType: "application/problem+json", status: 409,
			body: `{"status":409,"detail":"conflict"}`,
			call: func(c *v1.Client) error {
				_, err := NewActivationOp(c).Create(t.Context(), "test")
				return err
			},
			want: "API Error 409: conflict", conflict: true,
		},
		{
			name: "problem json rate limit", contentType: "application/problem+json", status: 429,
			body: `{"status":429,"detail":"too_many_requests"}`,
			call: func(c *v1.Client) error {
				return NewAutomatedActionsOp(c).Delete(t.Context(), "test")
			},
			want: "API Error 429: too_many_requests",
		},
		{
			name: "unexpected status on update", contentType: "application/json", status: 502,
			body: `{"detail":"upstream unavailable"}`,
			call: func(c *v1.Client) error {
				req := SetupEvaluationRuleInput(&EvaluationRuleInputParams{ID: "server-no-public-ip"})
				_, err := NewEvaluationRulesOp(c).Update(t.Context(), "server-no-public-ip", req)
				return err
			},
			want: "API Error 502: upstream unavailable",
		},
		{
			name: "missing unexpected detail", contentType: "application/problem+json", status: 502,
			body: `{}`,
			call: func(c *v1.Client) error {
				_, err := NewEvaluationRulesOp(c).Read(t.Context(), "test")
				return err
			},
			want: "API Error 502: API error response missing detail",
		},
		{
			name: "missing unexpected detail on delete", contentType: "application/json", status: 503,
			body: `{}`,
			call: func(c *v1.Client) error {
				return NewAutomatedActionsOp(c).Delete(t.Context(), "test")
			},
			want: "API Error 503: API error response missing detail",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.status)
				_, _ = fmt.Fprint(w, tt.body)
			}))
			defer server.Close()
			client, err := v1.NewClient(server.URL, dummySecuritySource{})
			require.NoError(t, err)
			err = tt.call(client)
			require.ErrorContains(t, err, tt.want)
			require.Equal(t, tt.notFound, saclient.IsNotFoundError(err))
			require.Equal(t, tt.conflict, saclient.IsConflictError(err))
		})
	}
}

func TestDeleteSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, err := v1.NewClient(server.URL, dummySecuritySource{})
	require.NoError(t, err)
	require.NoError(t, NewAutomatedActionsOp(client).Delete(t.Context(), "test"))
}
