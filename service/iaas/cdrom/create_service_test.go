package cdrom

import (
	"io"
	"testing"

	"bytes"
	"os"

	"github.com/sacloud/iaas-service-go/internal/ftps"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestService_Create(t *testing.T) {
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
	filename := "test-cdrom-upload-source.tmp"
	content := []byte("cdrom-upload-test")
	if err := os.WriteFile(filename, content, 0600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename) //nolint:errcheck

	// create
	cdrom, err := svc.Create(&CreateRequest{
		Zone:        zone,
		Name:        testutil.ResourceName("test-cdrom-upload-service"),
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2", "tag3"},
		SizeGB:      5,
		SourcePath:  filename,
	})
	if err != nil {
		t.Fatal(err)
	}

	// downloadして内容検証
	cdromOp := iaas.NewCDROMOp(caller)
	ftpServer, err := cdromOp.OpenFTP(t.Context(), zone, cdrom.ID, &iaas.OpenFTPRequest{
		ChangePassword: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cdromOp.CloseFTP(t.Context(), zone, cdrom.ID)
	}()

	ftpsClient, err := ftps.NewClient(ftpServer)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := ftpsClient.Retr("data.iso")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Close()
	}()
	buf, err := io.ReadAll(resp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf, content) {
		t.Fatalf("unexpected value: got:%s want:%s", buf, string(content))
	}

	if err := ftpsClient.Quit(); err != nil {
		t.Fatalf("closing FTP connection failed: %s", err)
	}

	// delete
	if err := cdromOp.Delete(t.Context(), zone, cdrom.ID); err != nil {
		t.Fatal(err)
	}
}
