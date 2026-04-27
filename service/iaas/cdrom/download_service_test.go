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

package cdrom

import (
	"bytes"
	"os"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestService_Download(t *testing.T) {
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
	filename := "test-cdrom-download-source.tmp"
	content := []byte("cdrom-download-test")
	if err := os.WriteFile(filename, content, 0600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename) //nolint:errcheck

	// create
	cdrom, err := svc.Create(&CreateRequest{
		Zone:        zone,
		Name:        testutil.ResourceName("test-cdrom-download-service"),
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2", "tag3"},
		SizeGB:      5,
		SourcePath:  filename,
	})
	if err != nil {
		t.Fatal(err)
	}

	// downloadして内容検証
	buf := bytes.NewBuffer([]byte{})
	err = svc.Download(&DownloadRequest{
		Zone:   zone,
		ID:     cdrom.ID,
		Writer: buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), content) {
		t.Fatalf("unexpected value: got:%s want:%s", buf.String(), string(content))
	}

	// delete
	cdromOp := iaas.NewCDROMOp(caller)
	if err := cdromOp.Delete(t.Context(), zone, cdrom.ID); err != nil {
		t.Fatal(err)
	}
}
