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
	"net/http"
	"testing"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestPublisherOp_List(t *testing.T) {
	expected := v1.PaginatedPublisherList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.Publisher{TemplatePublisher},
	}
	client := newTestClient(expected)
	api := NewPublisherOp(client)
	ctx := context.Background()

	publishers, err := api.List(ctx, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, publishers)
	require.Equal(t, 1, len(publishers))

	p := publishers[0]
	require.Equal(t, TemplatePublisher.GetCode(), p.GetCode())
	require.Equal(t, TemplatePublisher.GetDescription(), p.GetDescription())
	require.Equal(t, 3, len(p.GetVariants()))
	for i, v := range p.GetVariants() {
		ev := TemplatePublisher.GetVariants()[i]
		require.Equal(t, ev.GetLabel(), v.GetLabel())
		require.Equal(t, ev.GetName(), v.GetName())
		require.Equal(t, ev.GetStorage(), v.GetStorage())
		require.Equal(t, ev.GetSystem(), v.GetSystem())
	}
}

func TestPublisherOp_Read_200(t *testing.T) {
	var expected v1.WrappedPublisher
	var variant v1.PublisherVariant
	expected.SetFake()
	variant.SetFake()
	expected.SetVariants([]v1.PublisherVariant{variant})
	client := newTestClient(expected)
	api := NewPublisherOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, expected.GetCode())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, expected.GetCode(), actual.GetCode())
	require.Equal(t, expected.GetDescription(), actual.GetDescription())
	require.Equal(t, len(expected.GetVariants()), len(actual.GetVariants()))
	for i, va := range expected.GetVariants() {
		ve := expected.GetVariants()[i]
		require.Equal(t, ve.GetLabel(), va.GetLabel())
		require.Equal(t, ve.GetName(), va.GetName())
		require.Equal(t, ve.GetStorage(), va.GetStorage())
		require.Equal(t, ve.GetSystem(), va.GetSystem())
	}
}

func TestPublisherOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No Publisher matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewPublisherOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "nonexistent")
	require.Nil(t, actual)
	require.ErrorContains(t, err, "publisher not found")
}

func TestPublisherIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)

	api := NewPublisherOp(client)
	ctx := context.Background()

	// List all publishers
	publishers, err := api.List(ctx, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, publishers)
	require.NotEmpty(t, publishers)

	// Pick the first publisher and read it by code
	listed := publishers[0]
	code := listed.GetCode()
	read, err := api.Read(ctx, code)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, code, read.GetCode())
	require.Equal(t, listed.GetDescription(), read.GetDescription())
	require.Equal(t, len(listed.GetVariants()), len(read.GetVariants()))
}
