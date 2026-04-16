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

package builder

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/internal/ftps"
	"github.com/sacloud/packages-go/size"
)

// BlankArchiveBuilder ブランクアーカイブの作成〜FTPSでのファイルアップロードを行う
type BlankArchiveBuilder struct {
	Name         string
	Description  string
	Tags         types.Tags
	IconID       types.ID
	SizeGB       int
	SourceReader io.Reader

	NoWait bool

	Client *APIClient
}

// Validate 設定値の検証
func (b *BlankArchiveBuilder) Validate(ctx context.Context, zone string) error {
	if b.NoWait {
		return errors.New("NoWait=true is not supported when uploading files and creating archives")
	}
	requiredValues := map[string]bool{
		"Name":         b.Name == "",
		"SizeGB":       b.SizeGB == 0,
		"SourceReader": b.SourceReader == nil,
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}

// Build ブランクアーカイブの作成〜FTPSでのファイルアップロードを行う
func (b *BlankArchiveBuilder) Build(ctx context.Context, zone string) (*iaas.Archive, error) {
	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	archive, ftpServer, err := b.Client.Archive.CreateBlank(ctx, zone,
		&iaas.ArchiveCreateBlankRequest{
			Name:        b.Name,
			Description: b.Description,
			Tags:        b.Tags,
			IconID:      b.IconID,
			SizeMB:      b.SizeGB * size.GiB,
		})
	if err != nil {
		return nil, err
	}

	// upload sources via FTPS
	ftpsClient, err := ftps.NewClient(ftpServer)
	if err != nil {
		return nil, fmt.Errorf("cannot create ftp client: %s", err)
	}

	if err := ftpsClient.Stor("data.raw", b.SourceReader); err != nil {
		return archive, fmt.Errorf("uploading file via FTPS is failed: %s", err)
	}

	if err := ftpsClient.Quit(); err != nil {
		return nil, fmt.Errorf("closing FTP connection failed: %s", err)
	}

	// close FTP
	if err := b.Client.Archive.CloseFTP(ctx, zone, archive.ID); err != nil {
		return archive, err
	}

	// reload
	return b.Client.Archive.Read(ctx, zone, archive.ID)
}
