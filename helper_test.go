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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
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

var TemplateMetricsTank = func() v1.MetricsTank {
	var ret v1.MetricsTank

	ret.SetFake()
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	ret.SetUpdatedAt(TemplateTime)
	return ret
}()

var TemplateWrappedMetricsTank = func() v1.WrappedMetricsTank {
	var ret v1.WrappedMetricsTank

	ret.SetFake()
	ret.SetTags(TemplateTags)
	ret.SetCreatedAt(TemplateTime)
	ret.SetUpdatedAt(TemplateTime)
	return ret
}()

var TemplateMetricsTankAccessKey = func() v1.MetricsTankAccessKey {
	var ret v1.MetricsTankAccessKey

	ret.SetFake()
	return ret
}()

var TemplateMetricsRouting = func() v1.MetricsRouting {
	var r v1.MetricsRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetMetricsStorage(TemplateMetricsTank)
	return r
}()

var TemplateWrappedMetricsRouting = func() v1.WrappedMetricsRouting {
	var r v1.WrappedMetricsRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetMetricsStorage(TemplateMetricsTank)
	r.SetIsOk(true) // それはそう
	return r
}()

var TemplateWrappedAccessKey = func() v1.WrappedMetricsTankAccessKey {
	var ret v1.WrappedMetricsTankAccessKey

	ret.SetFake()
	return ret
}()

var TemplateLogTableAccessKey = func() v1.LogTableAccessKey {
	var ret v1.LogTableAccessKey

	ret.SetFake()
	return ret
}()

var TemplateLogTableEndpoints = func() v1.LogTableEndpoints {
	var ret v1.LogTableEndpoints

	ret.SetFake()
	return ret
}()

var TemplateWrappedLogTableEndpoints = func() v1.WrappedLogTableEndpoints {
	var ret v1.WrappedLogTableEndpoints

	ret.SetFake()
	return ret
}()

var TemplateLogTableUsage = func() v1.LogTableUsage {
	var ret v1.LogTableUsage

	ret.SetFake()
	return ret
}()

var TemplateWrappedLogTableUsage = func() v1.WrappedLogTableUsage {
	var ret v1.WrappedLogTableUsage

	ret.SetFake()
	return ret
}()

var TemplateLogTable = func() v1.LogTable {
	var ret v1.LogTable

	ret.SetFake()
	ret.SetEndpoints(TemplateLogTableEndpoints)
	ret.SetUsage(TemplateLogTableUsage)
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	return ret
}()

var TemplateWrappedLogTable = func() v1.WrappedLogTable {
	var ret v1.WrappedLogTable

	ret.SetFake()
	ret.SetEndpoints(TemplateWrappedLogTableEndpoints)
	ret.SetUsage(TemplateWrappedLogTableUsage)
	ret.SetCreatedAt(TemplateTime)
	ret.SetTags(TemplateTags)
	return ret
}()

var TemplateLogRouting = func() v1.LogRouting {
	var r v1.LogRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetLogStorage(TemplateLogTable)
	return r
}()

var TemplateWrappedLogRouting = func() v1.WrappedLogRouting {
	var r v1.WrappedLogRouting

	r.SetFake()
	r.SetPublisher(TemplatePublisher)
	r.SetLogStorage(TemplateLogTable)
	r.IsOk = true
	return r
}()

var TemplateDashboardProject = func() v1.DashboardProject {
	var ret v1.DashboardProject

	ret.SetFake()
	for _, tag := range []string{"tag1", "tag2"} {
		ret.Tags = append(ret.Tags, tag)
	}
	ret.SetCreatedAt(TemplateTime)
	return ret
}()

var TemplateWrappedDashboardProject = func() v1.WrappedDashboardProject {
	var ret v1.WrappedDashboardProject

	ret.SetFake()
	for _, tag := range []string{"tag1", "tag2"} {
		ret.Tags = append(ret.Tags, tag)
	}
	ret.SetCreatedAt(TemplateTime)
	ret.SetIsOk(true)
	return ret
}()
