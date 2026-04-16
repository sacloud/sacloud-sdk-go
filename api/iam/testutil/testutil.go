// Copyright 2025- The sacloud/iam-api-go Authors
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

package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/sacloud/iam-api-go"
	"github.com/sacloud/iam-api-go/apis/folder"
	"github.com/sacloud/iam-api-go/apis/project"
	"github.com/sacloud/iam-api-go/apis/serviceprincipal"
	"github.com/sacloud/iam-api-go/apis/user"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var theClient saclient.Client

func NewTestClient(v any, s ...int) *v1.Client {
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
	sa, e := theClient.DupWith(saclient.WithTestServer(sv))
	if e != nil {
		panic(e)
	}
	c, e := iam.NewClientWithAPIRootURL(sa, sv.URL)
	if e != nil {
		panic(e)
	}
	if e := sa.Populate(); e != nil {
		panic(e)
	}
	return c
}

func IntegratedClient(
	t *testing.T,
) (
	assert *require.Assertions,
	client *v1.Client,
) {
	acc, ok := os.LookupEnv("TESTACC")
	if !ok {
		t.Skip("environment variable TESTACC is not set. skip")
	}

	switch acc {
	case "1", "true", "TRUE", "True":
		// pass
	default:
		t.Skip("environment variable TESTACC is not set. skip")
	}

	assert = require.New(t)
	err := theClient.SetEnviron(os.Environ())
	assert.NoError(err)

	err = theClient.Populate()
	assert.NoError(err)

	authen, err := theClient.DupWith(
		saclient.WithFavouringBearerAuthentication(),
		saclient.WithTraceMode("error"),
	)
	assert.NoError(err)

	err = authen.Populate()
	assert.NoError(err)

	client, err = iam.NewClient(authen)
	assert.NoError(err)
	return assert, client
}

func Myself() (
	ret struct {
		PrincipalID int
		KID         string
	},
) {
	j := theClient.JSON()

	if v, ok := j["ServicePrincipalID"]; !ok {
		msg := fmt.Sprintf("failed to parse ServicePrincipalID: %#+v", j)
		panic(msg)
	} else if w, ok := v.(string); !ok {
		msg := fmt.Sprintf("failed to parse ServicePrincipalKeyID: %#+v", v)
		panic(msg)
	} else if x, err := strconv.ParseInt(w, 0, 64); err != nil {
		msg := fmt.Sprintf("failed to parse ServicePrincipalID: %s", err)
		panic(msg)
	} else {
		ret.PrincipalID = int(x)
	}

	if v, ok := j["ServicePrincipalKeyID"]; !ok {
		msg := fmt.Sprintf("failed to parse ServicePrincipalKeyID: %#+v", j)
		panic(msg)
	} else if w, ok := v.(string); !ok {
		msg := fmt.Sprintf("failed to parse ServicePrincipalKeyID: %#+v", v)
		panic(msg)
	} else {
		ret.KID = w
	}
	return ret
}

func NewFolder(
	t *testing.T,
	client *v1.Client,
) (
	ret *v1.Folder,
	deleter func(),
) {
	var err error
	api := iam.NewFolderOp(client)
	params := folder.CreateParams{
		Name: super.RandomName("folder", 32, super.CharSetAlphaNum),
	}
	if ret, err = api.Create(t.Context(), params); err != nil {
		t.Fatalf("Project.Create() failed: %s", err)
	} else {
		deleter = func() { _ = api.Delete(t.Context(), ret.GetID()) }
	}
	return
}

func NewProject(
	t *testing.T,
	client *v1.Client,
) (
	ret *v1.Project,
	deleter func(),
) {
	var err error
	api := iam.NewProjectOp(client)
	params := project.CreateParams{
		Name: super.RandomName("project", 32, super.CharSetAlphaNum),
		Code: super.RandomName("code", 16, super.CharSetAlphaNum),
	}
	if ret, err = api.Create(t.Context(), params); err != nil {
		t.Fatalf("Project.Create() failed: %s", err)
	} else {
		deleter = func() { _ = api.Delete(t.Context(), ret.GetID()) }
	}
	return
}

func NewUser(
	t *testing.T,
	client *v1.Client,
) (
	ret *v1.User,
	deleter func(),
) {
	var err error
	api := iam.NewUserOp(client)
	params := user.CreateParams{
		Name:     super.RandomName("user", 32, super.CharSetAlphaNum),
		Password: super.Random(64, super.CharSetAlphaNum),
		Code:     super.RandomName("code", 16, super.CharSetAlphaNum),
	}
	if ret, err = api.Create(t.Context(), params); err != nil {
		t.Fatalf("User.Create() failed: %s", err)
	} else {
		deleter = func() { _ = api.Delete(t.Context(), ret.GetID()) }
	}
	return
}

func NewPrincipal(
	t *testing.T,
	client *v1.Client,
	project int,
) (
	ret *v1.ServicePrincipal,
	deleter func(),
) {
	var err error
	api := iam.NewServicePrincipalOp(client)
	params := serviceprincipal.CreateParams{
		ProjectID: project,
		Name:      super.RandomName("user", 32, super.CharSetAlphaNum),
	}
	if ret, err = api.Create(t.Context(), params); err != nil {
		t.Fatalf("User.Create() failed: %s", err)
	} else {
		deleter = func() { _ = api.Delete(t.Context(), ret.GetID()) }
	}
	return
}
