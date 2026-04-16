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
	"github.com/sacloud/iaas-service-go/internal/ftps"
	"github.com/sacloud/packages-go/size"
)

func (s *Service) Create(req *CreateRequest) (*iaas.CDROM, error) {
	return s.CreateWithContext(context.Background(), req)
}

func (s *Service) CreateWithContext(ctx context.Context, req *CreateRequest) (*iaas.CDROM, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var reader io.Reader
	switch req.SourcePath {
	case "":
		reader = req.SourceReader
	default:
		file, err := os.Open(req.SourcePath)
		if err != nil {
			return nil, fmt.Errorf("reading source file[%s] failed: %s", req.SourcePath, err)
		}
		defer file.Close() //nolint:errcheck
		reader = file
	}

	client := iaas.NewCDROMOp(s.caller)
	cdrom, ftpServer, err := client.Create(ctx, req.Zone, &iaas.CDROMCreateRequest{
		SizeMB:      req.SizeGB * size.GiB,
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		IconID:      req.IconID,
	})
	if err != nil {
		return nil, err
	}

	ftpsClient, err := ftps.NewClient(ftpServer)
	if err != nil {
		return nil, fmt.Errorf("cannot create ftp client: %s", err)
	}

	if err := ftpsClient.Stor("data.iso", reader); err != nil {
		return nil, err
	}

	if err := ftpsClient.Quit(); err != nil {
		return nil, fmt.Errorf("closing FTP connection failed: %s", err)
	}

	if err := client.CloseFTP(ctx, req.Zone, cdrom.ID); err != nil {
		return nil, err
	}

	// refresh
	return client.Read(ctx, req.Zone, cdrom.ID)
}
