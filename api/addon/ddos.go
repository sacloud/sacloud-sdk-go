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

	v1 "github.com/sacloud/sacloud-sdk-go/api/addon/apis/v1"
)

type DDoSAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, params DDoSCreateParams) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type DDoSOp struct{ *v1.Client }

func NewDDoSOp(client *v1.Client) DDoSAPI { return &DDoSOp{Client: client} }

func (d *DDoSOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListDdosResourcesNotFound]("DDoS.List", func() (any, error) {
		return d.Client.ListDdosResources(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

type DDoSCreateParams struct {
	Location     string
	PricingLevel v1.PricingLevel
	Patterns     []string
	Origin       v1.FrontDoorOrigin
}

func (d *DDoSOp) Create(ctx context.Context, params DDoSCreateParams) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("DDoS.Create", func() (any, error) {
		if len(params.Patterns) == 0 {
			// avoid encode error
			return nil, NewError("DDoS.Create", fmt.Errorf("empty Patterns makes no sense"))
		} else {
			return d.Client.CreateDdosResource(ctx, &v1.DdosRequestBody{
				Location: v1.DdosRequestBodyLocation(params.Location),
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

func (d *DDoSOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.ReadDdosResourceNotFound]("DDoS.Read", func() (any, error) {
		return d.Client.ReadDdosResource(ctx, v1.ReadDdosResourceParams{
			ResourceGroupName: id,
		})
	})
}

func (d *DDoSOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteDdosResourceNoContent, v1.DeleteDdosResourceNotFound]("DDoS.Delete", func() (any, error) {
		return d.Client.DeleteDdosResource(ctx, v1.DeleteDdosResourceParams{
			ResourceGroupName: id,
		})
	})

	return
}

func (d *DDoSOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("DDoS.Status", func() (any, error) {
		return d.Client.ReadDdosResourceStatus(ctx, v1.ReadDdosResourceStatusParams{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ DDoSAPI = (*DDoSOp)(nil)
