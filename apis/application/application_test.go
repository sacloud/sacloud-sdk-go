// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package application_test

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/application"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api ApplicationAPI) {
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewApplicationOp(c)

	return
}

func TestList(t *testing.T) {
	next := "next-cursor"
	var expected v1.ListApplicationResponse
	expected.SetFake()
	expected.NextCursor.SetTo(next)
	expected.SetApplications(make([]v1.ReadApplicationDetail, 3))
	for i := 0; i < len(expected.GetApplications()); i++ {
		expected.Applications[i] = apprun_test.FakeApplication()
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Equal(expected.GetApplications(), actual)
}

func TestList_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, cursor, err := api.List(t.Context(), 0, nil)

	assert.Error(err)
	assert.Nil(actual)
	assert.Nil(cursor)
	assert.False(saclient.IsNotFoundError(err))
}

func TestCreate(t *testing.T) {
	var expected v1.CreateApplicationResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Create(t.Context(), "test-app", v1.ClusterID(uuid.New()))

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetApplication(), *actual)
}

func TestCreate_failed(t *testing.T) {
	expected := apprun_test.Fake400Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.Create(t.Context(), "", v1.ClusterID{})

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetApplicationResponse
	fake := apprun_test.FakeApplication()
	id := fake.ApplicationID
	expected.SetApplication(fake)

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(id, actual.ApplicationID)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ApplicationID(uuid.New())
	actual, err := api.Read(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.ApplicationID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ApplicationID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestContainers(t *testing.T) {
	var expected v1.GetApplicationContainersResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	id := v1.ApplicationID(uuid.New())
	actual, err := api.Containers(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Len(actual, 0)
}

func TestContainers_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ApplicationID(uuid.New())
	actual, err := api.Containers(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestUpdate(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.ApplicationID(uuid.New())
	version := int32(3)
	err := api.Update(t.Context(), id, &version)

	assert.NoError(err)
}

func TestUpdate_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ApplicationID(uuid.New())
	err := api.Update(t.Context(), id, nil)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)
	api := NewApplicationOp(client)

	assert.NotNil(api)

	cid, deleter := apprun_test.IntegratedCluster(t.Context(), assert, client)
	defer deleter()

	t.Run("Create", func(t *testing.T) {
		appName := super.RandomName("test-", 15, super.CharSetAlphaNum)
		app, err := api.Create(t.Context(), appName, cid)
		assert.NoError(err)
		assert.NotNil(app)

		aid := app.ApplicationID

		defer t.Run("Delete", func(t *testing.T) {
			err := api.Delete(t.Context(), aid)
			assert.NoError(err)
		})

		t.Run("List", func(t *testing.T) {
			list := apprun_test.RepeatedList(func(cursor *string) (res []v1.ReadApplicationDetail, next *string) {
				res, next, err := api.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(list)
		})

		t.Run("Read", func(t *testing.T) {
			actual, err := api.Read(t.Context(), aid)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(aid, actual.ApplicationID)
			assert.Equal(appName, actual.Name)
			assert.Equal(cid, actual.ClusterID)
		})

		t.Run("Update", func(t *testing.T) {
			err := api.Update(t.Context(), aid, nil)
			assert.NoError(err)
		})

		t.Run("Containers", func(t *testing.T) {
			containers, err := api.Containers(t.Context(), aid)
			assert.NoError(err)
			assert.NotNil(containers)
			// TODO: need to create containers and check the result
		})
	})
}
