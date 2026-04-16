// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package archive

import (
	"bytes"
	"os"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestArchiveService_downloadAfterBuild(t *testing.T) {
	if !testutil.IsAccTest() {
		t.SkipNow()
	}

	caller := api.NewCallerWithOptions(&api.CallerOptions{
		Options: &client.Options{
			AccessToken:       os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
			UserAgent:         "test-" + iaas.DefaultUserAgent,
			RetryMax:          20,
			Trace:             testutil.IsEnableTrace() || testutil.IsEnableHTTPTrace(),
		},
		TraceAPI: testutil.IsEnableTrace() || testutil.IsEnableHTTPTrace(),
	})
	zone := testutil.TestZone()
	svc := New(caller)

	// file
	filename := "test-archive-source.tmp"
	if err := os.WriteFile(filename, []byte("test"), 0600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename) //nolint:errcheck

	// create
	archive, err := svc.Create(&CreateRequest{
		Zone:        zone,
		Name:        testutil.ResourceName("test-archive-service"),
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2"},
		SizeGB:      20,
		SourcePath:  filename,
	})
	if err != nil {
		t.Fatal(err)
	}

	// update
	updName := archive.Name + "-upd"
	updArchive, err := svc.Update(&UpdateRequest{
		Zone: zone,
		ID:   archive.ID,
		Name: &updName,
	})
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, updName, updArchive.Name)
	require.Equal(t, archive.Description, updArchive.Description)
	require.Equal(t, archive.Tags, updArchive.Tags)
	require.Equal(t, archive.IconID, updArchive.IconID)

	// download
	buf := bytes.NewBuffer([]byte{})
	err = svc.Download(&DownloadRequest{
		Zone:   zone,
		ID:     archive.ID,
		Writer: buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	if buf.String() != "test" {
		t.Fatalf("unexpected value: got:%s want:%s", buf.String(), "test")
	}

	// delete
	if err := svc.Delete(&DeleteRequest{
		Zone: zone,
		ID:   archive.ID,
	}); err != nil {
		t.Fatal(err)
	}
}
