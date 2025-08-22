// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func newTestClient(v any, s ...int) *v1.Client {
	s = append(s, http.StatusOK)
	j, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(s[0])
		w.Write(j)
	})
	sv := httptest.NewServer(h)
	c, e := NewClientWithApiUrlAndClient(sv.URL, sv.Client())
	if e != nil {
		panic(e)
	}
	return c
}

func TestPublisherOp_List(t *testing.T) {
	expected := v1.PaginatedPublisherList{
		IsOk:  v1.NewOptBool(true),
		Count: 1,
		From:  0,
		Results: []v1.Publisher{
			{
				Code:        "test-publisher",
				Description: v1.NewOptString("This is a test publisher"),
				Variants: []v1.PublisherVariant{
					{
						Label:   "test-variant",
						Name:    "Test Variant",
						Storage: "metrics",
						System:  "disallow",
					},
				},
			},
		},
	}
	client := newTestClient(expected)
	api := NewPublisherOp(client)
	ctx := context.Background()

	publishers, err := api.List(ctx, 1, 0)
	require.NoError(t, err)
	require.NotNil(t, publishers)
	require.Equal(t, 1, len(publishers))

	p := publishers[0]
	require.Equal(t, "test-publisher", p.GetCode())
	require.Equal(t, v1.OptString{Value: "This is a test publisher", Set: true}, p.GetDescription())
	require.Equal(t, 1, len(p.GetVariants()))

	e := p.GetVariants()[0]
	require.Equal(t, "test-variant", e.GetLabel())
	require.Equal(t, "Test Variant", e.GetName())
	require.Equal(t, v1.PublisherVariantStorageMetrics, e.GetStorage())
	require.Equal(t, v1.PublisherVariantSystemDisallow, e.GetSystem())
}

func TestPublisherOp_Read_200(t *testing.T) {
	expected := v1.WrappedPublisher{
		IsOk:        true,
		Code:        "test-publisher",
		Description: v1.NewOptString("This is a test publisher"),
		Variants: []v1.PublisherVariant{
			{
				Label:   "test-variant",
				Name:    "Test Variant",
				Storage: "metrics",
				System:  "disallow",
			},
		},
	}
	client := newTestClient(expected)
	api := NewPublisherOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, expected.GetCode())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.GetCode(), actual.GetCode())
	require.Equal(t, expected.GetDescription(), actual.GetDescription())
	require.Equal(t, len(expected.GetVariants()), len(actual.GetVariants()))
	for i := range expected.GetVariants() {
		ve := expected.GetVariants()[i]
		va := actual.GetVariants()[i]
		require.Equal(t, ve.GetLabel(), va.GetLabel())
		require.Equal(t, ve.GetName(), va.GetName())
		require.Equal(t, ve.GetStorage(), va.GetStorage())
		require.Equal(t, ve.GetSystem(), va.GetSystem())
	}
}

func TestPublisherOp_Read_404(t *testing.T) {
	expected := struct {
		Code    string `json:"error_code"`
		Message string `json:"error_msg"`
		IsOk    bool   `json:"is_ok"`
		Status  int    `json:"status"`
	}{
		Code:    "not_found",
		Message: "No Publisher matches the given query.",
		IsOk:    false,
		Status:  404,
	}
	client := newTestClient(expected, http.StatusNotFound)
	api := NewPublisherOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "nonexistent")
	require.Nil(t, actual)
	require.ErrorContains(t, err, "publisher not found")
}
