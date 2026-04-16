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

package server

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
)

// APIClient builderが利用するAPIクライアント群
type APIClient struct {
	Disk         DiskHandler
	Interface    InterfaceHandler
	PacketFilter PacketFilterReader
	Server       CreateServerHandler
	ServerPlan   query.ServerPlanFinder
	Switch       SwitchReader
}

// DiskHandler ディスクの接続/切断のためのインターフェース
type DiskHandler interface {
	ConnectToServer(ctx context.Context, zone string, id types.ID, serverID types.ID) error
	DisconnectFromServer(ctx context.Context, zone string, id types.ID) error
}

// SwitchReader スイッチ参照のためのインターフェース
type SwitchReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Switch, error)
}

// InterfaceHandler NIC操作のためのインターフェース
type InterfaceHandler interface {
	Create(ctx context.Context, zone string, param *iaas.InterfaceCreateRequest) (*iaas.Interface, error)
	Update(ctx context.Context, zone string, id types.ID, param *iaas.InterfaceUpdateRequest) (*iaas.Interface, error)
	Delete(ctx context.Context, zone string, id types.ID) error
	ConnectToSharedSegment(ctx context.Context, zone string, id types.ID) error
	ConnectToSwitch(ctx context.Context, zone string, id types.ID, switchID types.ID) error
	DisconnectFromSwitch(ctx context.Context, zone string, id types.ID) error
	ConnectToPacketFilter(ctx context.Context, zone string, id types.ID, packetFilterID types.ID) error
	DisconnectFromPacketFilter(ctx context.Context, zone string, id types.ID) error
}

// PacketFilterReader パケットフィルタ参照のためのインターフェース
type PacketFilterReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.PacketFilter, error)
}

// CreateServerHandler サーバ操作のためのインターフェース
type CreateServerHandler interface {
	Create(ctx context.Context, zone string, param *iaas.ServerCreateRequest) (*iaas.Server, error)
	Update(ctx context.Context, zone string, id types.ID, param *iaas.ServerUpdateRequest) (*iaas.Server, error)
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Server, error)
	InsertCDROM(ctx context.Context, zone string, id types.ID, insertParam *iaas.InsertCDROMRequest) error
	EjectCDROM(ctx context.Context, zone string, id types.ID, ejectParam *iaas.EjectCDROMRequest) error
	Boot(ctx context.Context, zone string, id types.ID) error
	BootWithVariables(ctx context.Context, zone string, id types.ID, param *iaas.ServerBootVariables) error
	Shutdown(ctx context.Context, zone string, id types.ID, shutdownOption *iaas.ShutdownOption) error
	ChangePlan(ctx context.Context, zone string, id types.ID, plan *iaas.ServerChangePlanRequest) (*iaas.Server, error)
}

// NewBuildersAPIClient APIクライアントの作成
func NewBuildersAPIClient(caller iaas.APICaller) *APIClient {
	return &APIClient{
		Disk:         iaas.NewDiskOp(caller),
		Interface:    iaas.NewInterfaceOp(caller),
		PacketFilter: iaas.NewPacketFilterOp(caller),
		Server:       iaas.NewServerOp(caller),
		ServerPlan:   iaas.NewServerPlanOp(caller),
		Switch:       iaas.NewSwitchOp(caller),
	}
}
