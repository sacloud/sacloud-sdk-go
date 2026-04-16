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

package wait

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/defaults"
	"github.com/sacloud/iaas-api-go/types"
)

var (
	// ApplianceNotFoundRetryCount アプライアンスの待ち処理時に404エラーとなった場合のリトライ回数
	ApplianceNotFoundRetryCount = 30
	// InternetNotFoundRetryCount ルータの作成待ち処理時に404エラーとなった場合のリトライ回数
	InternetNotFoundRetryCount = 360
)

// UntilArchiveIsReady コピー完了まで待機
func UntilArchiveIsReady(ctx context.Context, client iaas.ArchiveAPI, zone string, id types.ID) (*iaas.Archive, error) {
	lastState, err := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Archive), err
	}
	return nil, err
}

// UntilDatabaseIsUp 起動まで待機
func UntilDatabaseIsUp(ctx context.Context, client iaas.DatabaseAPI, zone string, id types.ID) (*iaas.Database, error) {
	var database *iaas.Database
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, ApplianceNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		database = lastState.(*iaas.Database)
	}
	if err != nil {
		return nil, err
	}
	// [HACK] データベースアプライアンス場合のみ/appliance/:id/statusも考慮する
	waiter := iaas.WaiterForUp(func() (interface{}, error) {
		return client.Status(ctx, zone, id)
	})
	waiter.(*iaas.StatePollingWaiter).Interval = defaults.DefaultDBStatusPollingInterval // HACK 現状は決め打ち、ユースケースが出たら修正する

	_, err = waiter.WaitForState(ctx)
	return database, err
}

// UntilDatabaseIsDown シャットダウンまで待機
func UntilDatabaseIsDown(ctx context.Context, client iaas.DatabaseAPI, zone string, id types.ID) (*iaas.Database, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Database), err
	}
	return nil, err
}

// UntilDiskIsReady コピー完了/ディスク修正完了まで待機
func UntilDiskIsReady(ctx context.Context, client iaas.DiskAPI, zone string, id types.ID) (*iaas.Disk, error) {
	lastState, err := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Disk), err
	}
	return nil, err
}

// UntilInternetIsReady 準備完了まで待機
func UntilInternetIsReady(ctx context.Context, client iaas.InternetAPI, zone string, id types.ID) (*iaas.Internet, error) {
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, InternetNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Internet), err
	}
	return nil, err
}

// UntilLoadBalancerIsUp 起動完了まで待機
func UntilLoadBalancerIsUp(ctx context.Context, client iaas.LoadBalancerAPI, zone string, id types.ID) (*iaas.LoadBalancer, error) {
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, ApplianceNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.LoadBalancer), err
	}
	return nil, err
}

// UntilLoadBalancerIsDown シャットダウンまで待機
func UntilLoadBalancerIsDown(ctx context.Context, client iaas.LoadBalancerAPI, zone string, id types.ID) (*iaas.LoadBalancer, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.LoadBalancer), err
	}
	return nil, err
}

// UntilMobileGatewayIsReady コピー完了まで待機
func UntilMobileGatewayIsReady(ctx context.Context, client iaas.MobileGatewayAPI, zone string, id types.ID) (*iaas.MobileGateway, error) {
	lastState, err := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.MobileGateway), err
	}
	return nil, err
}

// UntilMobileGatewayIsUp 起動まで待機
func UntilMobileGatewayIsUp(ctx context.Context, client iaas.MobileGatewayAPI, zone string, id types.ID) (*iaas.MobileGateway, error) {
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, ApplianceNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.MobileGateway), err
	}
	return nil, err
}

// UntilMobileGatewayIsDown シャットダウンまで待機
func UntilMobileGatewayIsDown(ctx context.Context, client iaas.MobileGatewayAPI, zone string, id types.ID) (*iaas.MobileGateway, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.MobileGateway), err
	}
	return nil, err
}

// UntilNFSIsUp 起動まで待機
func UntilNFSIsUp(ctx context.Context, client iaas.NFSAPI, zone string, id types.ID) (*iaas.NFS, error) {
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, ApplianceNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.NFS), err
	}
	return nil, err
}

// UntilNFSIsDown シャットダウンまで待機
func UntilNFSIsDown(ctx context.Context, client iaas.NFSAPI, zone string, id types.ID) (*iaas.NFS, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.NFS), err
	}
	return nil, err
}

// UntilServerIsUp 起動まで待機
func UntilServerIsUp(ctx context.Context, client iaas.ServerAPI, zone string, id types.ID) (*iaas.Server, error) {
	lastState, err := iaas.WaiterForUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Server), err
	}
	return nil, err
}

// UntilServerIsDown シャットダウンまで待機
func UntilServerIsDown(ctx context.Context, client iaas.ServerAPI, zone string, id types.ID) (*iaas.Server, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.Server), err
	}
	return nil, err
}

// UntilVPCRouterIsReady コピー完了まで待機
func UntilVPCRouterIsReady(ctx context.Context, client iaas.VPCRouterAPI, zone string, id types.ID) (*iaas.VPCRouter, error) {
	lastState, err := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.VPCRouter), err
	}
	return nil, err
}

// UntilVPCRouterIsUp 起動まで待機
func UntilVPCRouterIsUp(ctx context.Context, client iaas.VPCRouterAPI, zone string, id types.ID) (*iaas.VPCRouter, error) {
	lastState, err := iaas.WaiterForApplianceUp(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}, ApplianceNotFoundRetryCount).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.VPCRouter), err
	}
	return nil, err
}

// UntilVPCRouterIsDown シャットダウンまで待機
func UntilVPCRouterIsDown(ctx context.Context, client iaas.VPCRouterAPI, zone string, id types.ID) (*iaas.VPCRouter, error) {
	lastState, err := iaas.WaiterForDown(func() (interface{}, error) {
		return client.Read(ctx, zone, id)
	}).WaitForState(ctx)
	if lastState != nil {
		return lastState.(*iaas.VPCRouter), err
	}
	return nil, err
}
