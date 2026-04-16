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

	v1 "github.com/sacloud/addon-api-go/apis/v1"
)

type AIAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, location string, sku v1.AiServiceSku) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type AIOp struct{ *v1.Client }

func NewAIOp(client *v1.Client) AIAPI { return &AIOp{Client: client} }

func (a *AIOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0501NotFound]("AI.List", func() (any, error) {
		return a.Client.ListAZ0501(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

func (a *AIOp) Create(ctx context.Context, location string, sku v1.AiServiceSku) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("AI.Create", func() (any, error) {
		return a.Client.CreateAZ0501(ctx, &v1.AiRequestBody{
			Location: location,
			Sku:      sku,
		})
	})
}

func (a *AIOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0501NotFound]("AI.Read", func() (any, error) {
		return a.Client.GetAZ0501(ctx, v1.GetAZ0501Params{
			ResourceGroupName: id,
		})
	})
}

func (a *AIOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0501NoContent, v1.DeleteAZ0501NotFound]("AI.Delete", func() (any, error) {
		return a.Client.DeleteAZ0501(ctx, v1.DeleteAZ0501Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (a *AIOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("AI.Status", func() (any, error) {
		return a.Client.StatusAZ0501(ctx, v1.StatusAZ0501Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ AIAPI = (*AIOp)(nil)
