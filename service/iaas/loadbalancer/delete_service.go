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

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/wait"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
)

func (s *Service) Delete(req *DeleteRequest) error {
	return s.DeleteWithContext(context.Background(), req)
}

func (s *Service) DeleteWithContext(ctx context.Context, req *DeleteRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	client := iaas.NewLoadBalancerOp(s.caller)
	target, err := client.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return serviceutil.HandleNotFoundError(err, !req.FailIfNotFound)
	}

	if !req.Force && target.InstanceStatus.IsUp() {
		return fmt.Errorf("target %s:%q has not yet shut down", req.Zone, req.ID)
	}

	if target.InstanceStatus.IsUp() {
		if err := client.Shutdown(ctx, req.Zone, req.ID, &iaas.ShutdownOption{Force: true}); err != nil {
			return err
		}
	}

	// 元の状態がUnknownでなければwait
	if target.InstanceStatus != types.ServerInstanceStatuses.Unknown {
		if _, err := wait.UntilLoadBalancerIsDown(ctx, client, req.Zone, req.ID); err != nil {
			return err
		}
	}

	if err := client.Delete(ctx, req.Zone, req.ID); err != nil {
		return serviceutil.HandleNotFoundError(err, !req.FailIfNotFound)
	}
	return nil
}
