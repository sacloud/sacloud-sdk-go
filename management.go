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
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
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
	ret, err := op.client.GetResourcesLimits(ctx)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusBadRequest:
			return nil, NewAPIError("Management.ResourceLimits", e.StatusCode, errors.Wrap(err, "insufficient privileges to issue this API"))
		default:
			return nil, NewAPIError("Management.ResourceLimits", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("Management.ResourceLimits", 0, err)
	} else {
		return ret, nil
	}
}

func (op *managementOp) ReadProvisioning(ctx context.Context) (*v1.Provisioning, error) {
	ret, err := op.client.GetProvisioningState(ctx)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusBadRequest:
			return nil, NewAPIError("Management.ReadProvisioning", e.StatusCode, errors.Wrap(err, "insufficient privileges to issue this API"))
		default:
			return nil, NewAPIError("Management.ReadProvisioning", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("Management.ReadProvisioning", 0, err)
	} else {
		return ret, nil
	}
}

type ProvisioningCreateParam struct {
	Logs    *v1.ProvisioningExist
	Metrics *v1.ProvisioningExist
}

func (op *managementOp) CreateProvisioning(ctx context.Context, p ProvisioningCreateParam) (*v1.Provisioning, error) {
	request := v1.ProvisioningCreate{
		Logs:    intoOpt[v1.OptProvisioningExist](p.Logs),
		Metrics: intoOpt[v1.OptProvisioningExist](p.Metrics),
	}
	ret, err := op.client.PostProvisioningInitialize(ctx, v1.NewOptProvisioningCreate(request))
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusBadRequest:
			return nil, NewAPIError("Management.CreateProvisioning", e.StatusCode, errors.Wrap(err, "insufficient privileges to issue this API"))
		default:
			return nil, NewAPIError("Management.CreateProvisioning", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("Management.CreateProvisioning", 0, err)
	} else {
		return ret, nil
	}
}
