// Copyright 2025- The sacloud/addon-api-go Authors
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

package addon

import (
	"context"
	"fmt"

	v1 "github.com/sacloud/addon-api-go/apis/v1"
)

type WAFAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, params WAFCreateParams) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type WAFOp struct{ *v1.Client }

func NewWAFOp(client *v1.Client) WAFAPI { return &WAFOp{Client: client} }

func (w *WAFOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0402NotFound]("WAF.List", func() (any, error) {
		return w.Client.ListAZ0402(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

type WAFCreateParams struct {
	Location     string
	PricingLevel v1.PricingLevel
	Patterns     []string
	Origin       v1.FrontDoorOrigin
}

func (w *WAFOp) Create(ctx context.Context, params WAFCreateParams) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("WAF.Create", func() (any, error) {
		if len(params.Patterns) == 0 {
			// avoid encode error
			return nil, NewError("WAF.Create", fmt.Errorf("empty Patterns makes no sense"))
		} else {
			return w.Client.CreateAZ0402(ctx, &v1.WafRequestBody{
				Location: params.Location,
				Profile: v1.FrontDoorProfile{
					Level: params.PricingLevel,
				},
				Endpoint: v1.FrontDoorEndpoint{
					Route: v1.FrontDoorRoute{
						Patterns: params.Patterns,
						OriginGroup: v1.FrontDoorOriginGroup{
							Origin: params.Origin,
						},
					},
				},
			})
		}
	})
}

func (w *WAFOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0402NotFound]("WAF.Read", func() (any, error) {
		return w.Client.GetAZ0402(ctx, v1.GetAZ0402Params{
			ResourceGroupName: id,
		})
	})
}

func (w *WAFOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0402NoContent, v1.DeleteAZ0402NotFound]("WAF.Delete", func() (any, error) {
		return w.Client.DeleteAZ0402(ctx, v1.DeleteAZ0402Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (w *WAFOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("WAF.Status", func() (any, error) {
		return w.Client.StatusAZ0402(ctx, v1.StatusAZ0402Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ WAFAPI = (*WAFOp)(nil)
