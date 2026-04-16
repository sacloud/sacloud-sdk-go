// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"context"

	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type ManagementAPI interface {
	ResourceLimits(ctx context.Context) (*v1.ResourcesLimits, error)
	ReadProvisioning(ctx context.Context) (*v1.Provisioning, error)
	CreateProvisioning(ctx context.Context, request ProvisioningCreateParam) (*v1.Provisioning, error)
}

var _ ManagementAPI = (*managementOp)(nil)

type managementOp struct {
	client *v1.Client
}

func NewManagementOp(client *v1.Client) ManagementAPI {
	return &managementOp{client: client}
}

func (op *managementOp) ResourceLimits(ctx context.Context) (*v1.ResourcesLimits, error) {
	return errorFromDecodedResponse("Management.ResourceLimits", func() (*v1.ResourcesLimits, error) {
		return op.client.GetResourcesLimits(ctx)
	})
}

func (op *managementOp) ReadProvisioning(ctx context.Context) (*v1.Provisioning, error) {
	return errorFromDecodedResponse("Management.ReadProvisioning", func() (*v1.Provisioning, error) {
		return op.client.GetProvisioningState(ctx)
	})
}

type ProvisioningCreateParam struct {
	Logs    *v1.ProvisioningExist
	Metrics *v1.ProvisioningExist
}

func (op *managementOp) CreateProvisioning(ctx context.Context, p ProvisioningCreateParam) (*v1.Provisioning, error) {
	return errorFromDecodedResponse("Management.CreateProvisioning", func() (ret *v1.Provisioning, err error) {
		ret = new(v1.Provisioning)

		res, err := op.client.PostProvisioningInitialize(ctx, v1.NewOptProvisioningCreateRequest(v1.ProvisioningCreateRequest{
			Logs: intoOpt[v1.OptProvisioningExistRequest](func() *v1.ProvisioningExistRequest {
				if p.Logs == nil {
					return nil
				}
				r := v1.ProvisioningExistRequest{
					SystemExist: p.Logs.SystemExist,
					UserExist:   p.Logs.UserExist,
				}
				return &r
			}()),
			Metrics: intoOpt[v1.OptProvisioningExistRequest](func() *v1.ProvisioningExistRequest {
				if p.Metrics == nil {
					return nil
				}
				r := v1.ProvisioningExistRequest{
					SystemExist: p.Metrics.SystemExist,
					UserExist:   p.Metrics.UserExist,
				}
				return &r
			}()),
		}))
		switch rem := res.(type) {
		case *v1.PostProvisioningInitializeOK:
			ret.SetLogs(rem.Logs)
			ret.SetMetrics(rem.Metrics)
		case *v1.PostProvisioningInitializeCreated:
			ret.SetLogs(rem.Logs)
			ret.SetMetrics(rem.Metrics)
		default:
			ret = nil
		}
		return
	})
}
