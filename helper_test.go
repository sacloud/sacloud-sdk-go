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
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	client "github.com/sacloud/api-client-go"
	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	Code    string `json:"error_code"`
	Message string `json:"error_msg"`
	IsOk    bool   `json:"is_ok"`
	Status  int    `json:"status"`
}

func newErrorResponse(status int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    http.StatusText(status),
		Message: message,
		IsOk:    false,
		Status:  status,
	}
}

func newTestClient(v any, s ...int) *v1.Client {
	s = append(s, http.StatusOK)
	j, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st := s[0]

		w.WriteHeader(st)
		if st == http.StatusNoContent {
			return
		}
		if _, e = w.Write(j); e != nil {
			panic(e)
		}
	})
	sv := httptest.NewServer(h)
	c, e := NewClientWithApiUrlAndClient(sv.URL, sv.Client())
	if e != nil {
		panic(e)
	}
	return c
}

func IntegratedClient(t *testing.T, params ...client.ClientParam) (*v1.Client, error) {
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
	)(t)

	apiUrl := DefaultAPIRootURL
	if root, ok := os.LookupEnv("SAKURACLOUD_LOCAL_ENDPOINT_MONITORINGSUITE"); ok {
		apiUrl = root
	}
	return NewClientWithApiUrl(apiUrl, append(params, client.WithApiKeys(
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
	))...)
}

func WithAlertProject(t *testing.T, cli *v1.Client, ctx context.Context) *v1.AlertProject {
	op := NewAlertProjectOp(cli)

	ret, err := op.Create(ctx, AlertProjectCreateParams{
		Name:        testutil.RandomName("test-alert-project-", 16, testutil.CharSetAlphaNum),
		Description: ref(testutil.Random(128, testutil.CharSetAlphaNum)),
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	id, ok := ret.GetResourceID().Get()
	require.True(t, ok)

	t.Cleanup(func() {
		aid := fmt.Sprintf("%d", id)
		err := op.Delete(ctx, aid)
		require.NoError(t, err)
	})
	return ret
}

func WithMetricsStorage(t *testing.T, cli *v1.Client, ctx context.Context) *v1.MetricsStorage {
	op := NewMetricsStorageOp(cli)

	ret, err := op.Create(ctx, MetricsStorageCreateParams{
		Name:        testutil.RandomName("test-metrics-storage-", 16, testutil.CharSetAlphaNum),
		Description: ref(testutil.Random(128, testutil.CharSetAlphaNum)),
		IsSystem:    false,
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	id, ok := ret.GetResourceID().Get()
	require.True(t, ok)

	t.Cleanup(func() {
		mid := fmt.Sprintf("%d", id)
		err := op.Delete(ctx, mid)
		require.NoError(t, err)
	})
	return ret
}

func WithLogStorage(t *testing.T, cli *v1.Client, ctx context.Context) *v1.LogStorage {
	op := NewLogsStorageOp(cli)

	ret, err := op.Create(ctx, LogStorageCreateParams{
		Name:           testutil.RandomName("test-log-storage-", 16, testutil.CharSetAlphaNum),
		Description:    ref(testutil.Random(128, testutil.CharSetAlphaNum)),
		IsSystem:       false,
		Classification: ref(v1.LogStorageCreateClassificationShared),
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	id, ok := ret.GetResourceID().Get()
	require.True(t, ok)

	t.Cleanup(func() {
		lid := fmt.Sprintf("%d", id)
		err := op.Delete(ctx, lid)
		require.NoError(t, err)
	})
	return ret
}

func WithTraceStorage(t *testing.T, cli *v1.Client, ctx context.Context) *v1.TraceStorage {
	op := NewTracesStorageOp(cli)

	ret, err := op.Create(ctx, TracesStorageCreateParams{
		Name:           testutil.RandomName("test-trace-storage-", 16, testutil.CharSetAlphaNum),
		Description:    ref(testutil.Random(128, testutil.CharSetAlphaNum)),
		Classification: ref(v1.TraceStorageCreateClassificationShared),
	})
	require.NoError(t, err)
	require.NotNil(t, ret)

	id := ret.GetResourceID()

	t.Cleanup(func() {
		tid := fmt.Sprintf("%d", id)
		err := op.Delete(ctx, tid)
		require.NoError(t, err)
	})
	return ret
}

// generic-ish type cast helper function
func ref[T any](v T) *T { return &v }

// time.Now() をexpectationに使うのは筋悪である(SetFakeのままだとそうなる)
var TemplateTime time.Time = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

var TemplateTags = []string{"tag1", "tag2"}

var TemplatePublisher = func() v1.Publisher {
	var ret v1.Publisher

	ret.SetFake()
	for range 3 {
		var v v1.PublisherVariant

		v.SetFake()
		ret.Variants = append(ret.Variants, v)
	}
	return ret
}()

var TemplateMetricsStorage = func() v1.MetricsStorage {
	var ret v1.MetricsStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	ret.SetUpdatedAt(TemplateTime)
	return ret
}()

var TemplateWrappedMetricsStorage = func() v1.WrappedMetricsStorage {
	var ret v1.WrappedMetricsStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	ret.SetUpdatedAt(TemplateTime)
	return ret
}()

var TemplateMetricsStorageAccessKey = func() v1.MetricsStorageAccessKey {
	var ret v1.MetricsStorageAccessKey

	ret.SetFake()
	return ret
}()

var TemplateMetricsRouting = func() v1.MetricsRouting {
	var r v1.MetricsRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetMetricsStorage(TemplateMetricsStorage)
	return r
}()

var TemplateWrappedMetricsRouting = func() v1.WrappedMetricsRouting {
	var r v1.WrappedMetricsRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetMetricsStorage(TemplateMetricsStorage)
	r.SetIsOk(true) // それはそう
	return r
}()

var TemplateWrappedAccessKey = func() v1.WrappedMetricsStorageAccessKey {
	var ret v1.WrappedMetricsStorageAccessKey

	ret.SetFake()
	return ret
}()

var TemplateLogStorageAccessKey = func() v1.LogStorageAccessKey {
	var ret v1.LogStorageAccessKey

	ret.SetFake()
	return ret
}()

var TemplateLogStorageEndpoints = func() v1.LogStorageEndpoints {
	var ret v1.LogStorageEndpoints

	ret.SetFake()
	return ret
}()

var TemplateWrappedLogStorageEndpoints = func() v1.WrappedLogStorageEndpoints {
	var ret v1.WrappedLogStorageEndpoints

	ret.SetFake()
	return ret
}()

var TemplateLogStorageUsage = func() v1.LogStorageUsage {
	var ret v1.LogStorageUsage

	ret.SetFake()
	return ret
}()

var TemplateWrappedLogStorageUsage = func() v1.WrappedLogStorageUsage {
	var ret v1.WrappedLogStorageUsage

	ret.SetFake()
	return ret
}()

var TemplateLogStorage = func() v1.LogStorage {
	var ret v1.LogStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetEndpoints(TemplateLogStorageEndpoints)
	ret.SetUsage(TemplateLogStorageUsage)
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	return ret
}()

var TemplateWrappedLogStorage = func() v1.WrappedLogStorage {
	var ret v1.WrappedLogStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetEndpoints(TemplateWrappedLogStorageEndpoints)
	ret.SetUsage(TemplateWrappedLogStorageUsage)
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	return ret
}()

var TemplateAlertProject = func() v1.AlertProject {
	var ret v1.AlertProject

	ret.SetFake()
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	return ret
}()

var TemplateWrappedAlertProject = func() v1.WrappedAlertProject {
	var ret v1.WrappedAlertProject

	ret.SetFake()
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	ret.SetIsOk(true)
	return ret
}()

var TemplateLogRouting = func() v1.LogRouting {
	var r v1.LogRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetLogStorage(TemplateLogStorage)
	return r
}()

var TemplateWrappedLogRouting = func() v1.WrappedLogRouting {
	var r v1.WrappedLogRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetLogStorage(TemplateLogStorage)
	r.SetIsOk(true)
	return r
}()

var TemplateDashboardProject = func() v1.DashboardProject {
	var ret v1.DashboardProject

	ret.SetFake()
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	return ret
}()

var TemplateWrappedDashboardProject = func() v1.WrappedDashboardProject {
	var ret v1.WrappedDashboardProject

	ret.SetFake()
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	ret.SetIsOk(true)
	return ret
}()

var TemplateNotificationTarget = func() v1.NotificationTarget {
	var ret v1.NotificationTarget

	ret.SetFake()
	ret.SetProjectID(v1.NewNilInt64(^0))
	return ret
}()

var TemplateHistory = func() v1.History {
	var ret v1.History

	ret.SetFake()
	return ret
}()

var TemplateAlertRule = func() v1.AlertRule {
	var ret v1.AlertRule

	ret.SetFake()
	return ret
}()

var TemplateTraceStorage = func() v1.TraceStorage {
	var ret v1.TraceStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	return ret
}()

var TemplateWrappedTraceStorage = func() v1.WrappedTraceStorage {
	var ret v1.WrappedTraceStorage

	ret.SetFake()
	ret.SetID(^0)
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	return ret
}()

var TemplateTraceStorageAccessKey = func() v1.TraceStorageAccessKey {
	var ret v1.TraceStorageAccessKey

	ret.SetFake()
	return ret
}()

var TemplateWrappedTraceStorageAccessKey = func() v1.WrappedTraceStorageAccessKey {
	var ret v1.WrappedTraceStorageAccessKey

	ret.SetFake()
	return ret
}()

var TemplateLogMeasureRule = func() v1.LogMeasureRule {
	var ret v1.LogMeasureRule

	ret.SetFake()
	ret.SetLogStorage(TemplateLogStorage)
	ret.SetMetricsStorage(TemplateMetricsStorage)
	ret.Rule.Query.SetMatchers([]v1.FieldMatcher{})
	return ret
}()
