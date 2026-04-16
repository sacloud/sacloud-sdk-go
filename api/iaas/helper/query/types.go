// Copyright 2016-2022 The sacloud/iaas-api-go Authors
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

package query

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// ArchiveFinder アーカイブ検索インターフェース
type ArchiveFinder interface {
	Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.ArchiveFindResult, error)
}

// NoteFinder スタートアップスクリプト(Note)検索インターフェース
type NoteFinder interface {
	Find(ctx context.Context, conditions *iaas.FindCondition) (*iaas.NoteFindResult, error)
}

// ArchiveSourceReader アーカイブソースを取得するためのインターフェース
type ArchiveSourceReader struct {
	ArchiveReader ArchiveReader
	DiskReader    DiskReader
}

// NewArchiveSourceReader デフォルトのリーダーを返す
func NewArchiveSourceReader(caller iaas.APICaller) *ArchiveSourceReader {
	return &ArchiveSourceReader{
		ArchiveReader: iaas.NewArchiveOp(caller),
		DiskReader:    iaas.NewDiskOp(caller),
	}
}

// ServerPlanFinder .
type ServerPlanFinder interface {
	Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.ServerPlanFindResult, error)
}

// ServerSourceReader サーバのコピー元情報を参照するためのリーダー
type ServerSourceReader struct {
	ServerReader  ServerReader
	ArchiveReader ArchiveReader
	DiskReader    DiskReader
}

// NewServerSourceReader デフォルトのリーダーを返す
func NewServerSourceReader(caller iaas.APICaller) *ServerSourceReader {
	return &ServerSourceReader{
		ServerReader:  iaas.NewServerOp(caller),
		ArchiveReader: iaas.NewArchiveOp(caller),
		DiskReader:    iaas.NewDiskOp(caller),
	}
}

// ServerReader サーバ参照インターフェース
type ServerReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Server, error)
}

// ArchiveReader アーカイブ参照インターフェース
type ArchiveReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Archive, error)
}

// DiskReader ディスク参照インターフェース
type DiskReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Disk, error)
}
