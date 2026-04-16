// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package service_class_test

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/service_class"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api ServiceClassAPI) {
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewServiceClassOp(c)

	return
}

func TestListLB(t *testing.T) {
	var expected v1.ListLbServiceClassResponse
	expected.SetFake()
	expected.SetLbServiceClasses(make([]v1.ReadLbServiceClass, 2))
	for i := 0; i < len(expected.GetLbServiceClasses()); i++ {
		expected.LbServiceClasses[i].SetFake()
	}
	assert, api := setup(t, &expected)

	actual, err := api.ListLB(t.Context())

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Len(actual, 2)
}

func TestListLB_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.ListLB(t.Context())

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestListWorker(t *testing.T) {
	var expected v1.ListWorkerServiceClassResponse
	expected.SetFake()
	expected.SetWorkerServiceClasses(make([]v1.ReadWorkerServiceClass, 3))
	for i := 0; i < len(expected.GetWorkerServiceClasses()); i++ {
		expected.WorkerServiceClasses[i].SetFake()
	}
	assert, api := setup(t, &expected)

	actual, err := api.ListWorker(t.Context())

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Len(actual, 3)
}

func TestListWorker_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.ListWorker(t.Context())

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)

	api := NewServiceClassOp(client)
	assert.NotNil(api)

	t.Run("ListLB", func(t *testing.T) {
		actual, err := api.ListLB(t.Context())
		assert.NoError(err)
		assert.NotNil(actual)
	})

	t.Run("ListWorker", func(t *testing.T) {
		actual, err := api.ListWorker(t.Context())
		assert.NoError(err)
		assert.NotNil(actual)
	})
}
