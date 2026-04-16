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

type CDNAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, params CDNCreateParams) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type CDNOp struct{ *v1.Client }

func NewCDNOp(client *v1.Client) CDNAPI { return &CDNOp{Client: client} }

func (c *CDNOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0301NotFound]("CDN.List", func() (any, error) {
		return c.Client.ListAZ0301(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

type CDNCreateParams struct {
	Location     string
	PricingLevel v1.PricingLevel
	Patterns     []string
	Origin       v1.FrontDoorOrigin
}

func (c *CDNOp) Create(ctx context.Context, params CDNCreateParams) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("CDN.Create", func() (any, error) {
		if len(params.Patterns) == 0 {
			// avoid encode error
			return nil, NewError("CDN.Create", fmt.Errorf("empty Patterns makes no sense"))
		} else {
			return c.Client.CreateAZ0301(ctx, &v1.NetworkRequestBody{
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

func (c *CDNOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0301NotFound]("CDN.Read", func() (any, error) {
		return c.Client.GetAZ0301(ctx, v1.GetAZ0301Params{
			ResourceGroupName: id,
		})
	})
}

func (c *CDNOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0301NoContent, v1.DeleteAZ0301NotFound]("CDN.Delete", func() (any, error) {
		return c.Client.DeleteAZ0301(ctx, v1.DeleteAZ0301Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (c *CDNOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("CDN.Status", func() (any, error) {
		return c.Client.StatusAZ0301(ctx, v1.StatusAZ0301Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ CDNAPI = (*CDNOp)(nil)
