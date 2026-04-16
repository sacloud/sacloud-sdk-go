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
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/internal/ftps"
)

func (s *Service) Download(req *DownloadRequest) error {
	return s.DownloadWithContext(context.Background(), req)
}

func (s *Service) DownloadWithContext(ctx context.Context, req *DownloadRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	client := iaas.NewCDROMOp(s.caller)
	resource, err := client.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return fmt.Errorf("reading CDROM[%s] failed: %s", req.ID, err)
	}

	if resource.Scope != types.Scopes.User {
		return fmt.Errorf("CDROM[%s] is not allowed to download", req.ID)
	}

	ftpServer, err := client.OpenFTP(ctx, req.Zone, req.ID, &iaas.OpenFTPRequest{ChangePassword: req.ChangePassword})
	if err != nil {
		return fmt.Errorf("requesting FTP server information failed: %s", err)
	}

	ftpsClient, err := ftps.NewClient(ftpServer)
	if err != nil {
		return fmt.Errorf("cannot create ftp client: %s", err)
	}
	switch req.Path {
	case "":
		var out io.Writer = os.Stdout
		if req.Writer != nil {
			out = req.Writer
		}
		if err := ftps.DownloadWriter(ftpsClient, out); err != nil {
			return fmt.Errorf("downloading via FTP failed: %s", err)
		}
	default:
		if err := ftps.Download(ftpsClient, req.Path); err != nil {
			return fmt.Errorf("downloading via FTP failed: %s", err)
		}
	}

	if err := ftpsClient.Quit(); err != nil {
		return fmt.Errorf("closing FTP connection failed: %s", err)
	}

	// close
	if err := client.CloseFTP(ctx, req.Zone, req.ID); err != nil {
		return fmt.Errorf("closing FTP server failed: %s", err)
	}
	return nil
}
