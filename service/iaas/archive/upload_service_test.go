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
)

func TestArchiveService_uploadAndDownload(t *testing.T) {
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
	filename := "test-archive-upload-source.tmp"
	content := []byte("upload-test")
	if err := os.WriteFile(filename, content, 0600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename) //nolint:errcheck

	// create (空のアーカイブ)
	// 本当に blank なアーカイブを作成することは service 経由だとできないので、低レベル API を呼ぶ。
	archiveOp := iaas.NewArchiveOp(caller)
	archive, _, err := archiveOp.CreateBlank(t.Context(), zone,
		&iaas.ArchiveCreateBlankRequest{
			SizeMB:      1048576,
			Name:        testutil.ResourceName("test-archive-upload-service"),
			Description: "desc",
			Tags:        types.Tags{"tag1", "tag2", "tag3"},
		})
	if err != nil {
		t.Fatal(err)
	}

	// upload
	err = svc.Upload(&UploadRequest{
		Zone: zone,
		ID:   archive.ID,
		Path: filename,
	})
	if err != nil {
		t.Fatal(err)
	}

	// downloadして内容検証
	buf := bytes.NewBuffer([]byte{})
	err = svc.Download(&DownloadRequest{
		Zone:   zone,
		ID:     archive.ID,
		Writer: buf,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), content) {
		t.Fatalf("unexpected value: got:%s want:%s", buf.String(), string(content))
	}

	// delete
	if err := svc.Delete(&DeleteRequest{
		Zone: zone,
		ID:   archive.ID,
	}); err != nil {
		t.Fatal(err)
	}
}
