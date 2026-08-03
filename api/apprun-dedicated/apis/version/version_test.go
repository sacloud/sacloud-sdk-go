// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package version_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/version"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (*assert.Assertions, VersionAPI) {
	c, e := apprun_test.NewTestClient(v, s...)
	require.NoError(t, e)

	return assert.New(t), NewVersionOp(c, v1.ApplicationID(uuid.New()))
}

func TestList(t *testing.T) {
	var expected v1.ListApplicationVersionResponse
	expected.SetFake()
	expected.NextCursor.SetTo(2)
	expected.SetVersions(make([]v1.ApplicationVersionDeploymentStatus, 3))
	for i := 0; i < len(expected.GetVersions()); i++ {
		expected.Versions[i].SetFake()
		//nolint: gosec // this never integer overflows
		expected.Versions[i].SetVersion(v1.ApplicationVersionNumber(i + 1))
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	// cursor may be nil if there are no more results
	if cursor != nil {
		assert.GreaterOrEqual(*cursor, v1.ApplicationVersionNumber(1))
	}
	assert.Equal(expected.GetVersions(), actual)
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
	var expected v1.CreateApplicationVersionResponse
	expected.SetFake()
	expected.ApplicationVersion.SetVersion(1)
	assert, api := setup(t, &expected)

	actual, err := api.Create(t.Context(), CreateParams{
		Image:                  "nginx:latest",
		CPU:                    1000,
		Memory:                 512,
		ScalingMode:            v1.ScalingModeManual,
		FixedScale:             saclient.Ptr(int32(27)),
		Cmd:                    []string{"/bin/sh"},
		RegistryUsername:       saclient.Ptr("username"),
		RegistryPassword:       saclient.Ptr("password"),
		RegistryPasswordAction: v1.RegistryPasswordActionKeep,
		ExposedPorts: []ExposedPort{{
			TargetPort:       v1.Port(8080),
			LoadBalancerPort: saclient.Ptr(v1.Port(80)),
			UseLetsEncrypt:   true,
			Host:             []string{apprun_test.FakeCN()},
			HealthCheck: &v1.HealthCheck{
				Path:            "/status",
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
			},
		}},
		EnvVars: []EnvironmentVariable{
			{
				Key:    "SAKURA_ACCESS_TOKEN",
				Value:  saclient.Ptr("token"),
				Secret: false,
			},
			{
				Key:    "SAKURA_ACCESS_TOKEN_SECRET",
				Value:  nil,
				Secret: true,
			},
		},
	})

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetApplicationVersion(), *actual)
}

func TestCreate_failed(t *testing.T) {
	expected := apprun_test.Fake400Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.Create(t.Context(), CreateParams{})

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetApplicationVersionResponse
	ver := v1.ApplicationVersionNumber(1)
	expected.SetFake()
	expected.SetApplicationVersion(apprun_test.FakeApplicationVersion())

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), ver)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(ver, actual.Version)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	ver := v1.ApplicationVersionNumber(1)
	actual, err := api.Read(t.Context(), ver)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	ver := v1.ApplicationVersionNumber(1)
	err := api.Delete(t.Context(), ver)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	ver := v1.ApplicationVersionNumber(1)
	err := api.Delete(t.Context(), ver)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)

	assert.NotNil(client)

	cid, clusterDeleter := apprun_test.IntegratedCluster(t.Context(), assert, client)
	defer clusterDeleter()

	aid, appDeleter := apprun_test.IntegratedApplication(t.Context(), assert, client, cid)
	defer appDeleter()

	t.Run("Create", func(t *testing.T) {
		api := NewVersionOp(client, aid)
		assert.NotNil(api)

		version, err := api.Create(t.Context(), CreateParams{
			Image:                  "nginx:latest",
			CPU:                    1000,
			Memory:                 512,
			ScalingMode:            v1.ScalingModeManual,
			FixedScale:             saclient.Ptr(int32(1)),
			Cmd:                    []string{"/bin/sh"},
			RegistryUsername:       nil,
			RegistryPassword:       nil,
			RegistryPasswordAction: v1.RegistryPasswordActionRemove,
			ExposedPorts:           nil,
			EnvVars:                nil,
		})
		assert.NoError(err)
		assert.NotNil(version)

		vid := version.Version

		defer t.Run("Delete", func(t *testing.T) {
			err := api.Delete(t.Context(), vid)
			assert.NoError(err)
		})

		t.Run("List", func(t *testing.T) {
			list := apprun_test.RepeatedList(func(cursor *v1.ApplicationVersionNumber) (res []v1.ApplicationVersionDeploymentStatus, next *v1.ApplicationVersionNumber) {
				res, next, err := api.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(list)
		})

		t.Run("Read", func(t *testing.T) {
			actual, err := api.Read(t.Context(), vid)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(vid, actual.Version)
		})
	})
}

func TestJSONTags(t *testing.T) {
	assert := assert.New(t)

	params := CreateParams{
		CPU:     100,
		Memory:  128,
		Image:   "nginx:latest",
		Cmd:     []string{"/bin/sh"},
		EnvVars: []EnvironmentVariable{{Key: "K", Value: saclient.Ptr("V")}},
		ExposedPorts: []ExposedPort{{
			TargetPort:     8080,
			UseLetsEncrypt: true,
			Host:           []string{"example.com"},
			HealthCheck: &v1.HealthCheck{
				Path:            "/health",
				IntervalSeconds: 10,
				TimeoutSeconds:  5,
			},
		}},
	}

	b, err := json.Marshal(params)
	assert.NoError(err)

	var raw map[string]any
	assert.NoError(json.Unmarshal(b, &raw))

	assert.Equal(float64(100), raw["cpu"])
	assert.Equal(float64(128), raw["memory"])
	assert.Equal("nginx:latest", raw["image"])
	assert.Equal([]any{"/bin/sh"}, raw["cmd"])

	envVars, ok := raw["envVars"].([]any)
	assert.True(ok)
	assert.Len(envVars, 1)
	env := envVars[0].(map[string]any)
	assert.Equal("K", env["key"])
	assert.Equal("V", env["value"])
	// bool fields are not omitempty: false is a meaningful value and is
	// always serialized, so unset/false can be distinguished from omission.
	assert.Equal(false, env["secret"])

	ports, ok := raw["exposedPorts"].([]any)
	assert.True(ok)
	assert.Len(ports, 1)
	port := ports[0].(map[string]any)
	assert.Equal(float64(8080), port["targetPort"])
	assert.Equal(true, port["useLetsEncrypt"])
	assert.Equal([]any{"example.com"}, port["host"])

	hc := port["healthCheck"].(map[string]any)
	assert.Equal("/health", hc["path"])
	assert.Equal(float64(10), hc["intervalSeconds"])
	assert.Equal(float64(5), hc["timeoutSeconds"])

	// optional fields are omitted when unset (omitempty)
	assert.NotContains(raw, "fixedScale")
	assert.NotContains(raw, "registryUsername")
	assert.NotContains(raw, "registryPassword")

	// round-trip
	var restored CreateParams
	assert.NoError(json.Unmarshal(b, &restored))
	assert.Equal(params.CPU, restored.CPU)
	assert.Equal(params.Memory, restored.Memory)
	assert.Equal(params.Image, restored.Image)
	assert.Equal(params.Cmd, restored.Cmd)
	assert.Len(restored.EnvVars, 1)
	assert.Equal("K", restored.EnvVars[0].Key)
	assert.Equal("V", *restored.EnvVars[0].Value)
}
