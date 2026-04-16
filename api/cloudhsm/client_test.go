// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

package cloudhsm

import (
	"net/url"
	"reflect"
	"testing"
	"unsafe"

	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var testingClient saclient.Client

func TestNewClient(t *testing.T) {
	c, err := NewClient(&testingClient)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewClientWithApiUrl_OtherZone(t *testing.T) {
	c, err := NewClientWithApiUrl("https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1/", &testingClient)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewClient_WithZone(tt *testing.T) {
	client := testingClient.Dup()
	err := client.SetEnviron([]string{"SAKURA_ZONE=tk1a"})
	require.NoError(tt, err)
	err = client.Populate()
	require.NoError(tt, err)
	c, err := NewClient(client)
	require.NoError(tt, err)
	require.NotNil(tt, c)

	// この作成したクライアントが本当にtk1aを向いているかを確認するのがやや困難である
	q := reflect.ValueOf(c)
	w := q.Elem()
	e := w.FieldByName("serverURL")
	r := e.Type()
	t := unsafe.Pointer(e.UnsafeAddr()) //nolint:gosec
	y := reflect.NewAt(r, t)
	u := y.Elem()
	i := u.Interface()
	o, p := i.(*url.URL)

	require.True(tt, p)
	require.Contains(tt, o.String(), "/zone/tk1a/")
}
