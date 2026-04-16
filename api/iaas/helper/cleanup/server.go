// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package cleanup

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/types"
)

// DeleteServer サーバの削除を行う。引数に応じて接続されたディスクの削除も同時に行う。
// もし電源がONの場合は強制シャットダウンされる
func DeleteServer(ctx context.Context, caller iaas.APICaller, zone string, id types.ID, withDisks bool) error {
	serverOp := iaas.NewServerOp(caller)
	server, err := serverOp.Read(ctx, zone, id)
	if err != nil {
		return err
	}

	if server.InstanceStatus.IsUp() {
		if err := power.ShutdownServer(ctx, serverOp, zone, id, true); err != nil {
			return err
		}
	}

	if !withDisks {
		return serverOp.Delete(ctx, zone, id)
	}

	var diskIDs []types.ID
	for i := range server.Disks {
		diskIDs = append(diskIDs, server.Disks[i].ID)
	}

	return serverOp.DeleteWithDisks(ctx, zone, id, &iaas.ServerDeleteWithDisksRequest{IDs: diskIDs})
}
