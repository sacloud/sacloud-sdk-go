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

type SearchAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, params SearchCreateParams) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type SearchOp struct{ *v1.Client }

func NewSearchOp(client *v1.Client) SearchAPI { return &SearchOp{Client: client} }

func (s *SearchOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0105NotFound]("Search.List", func() (any, error) {
		return s.Client.ListAZ0105(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

type SearchCreateParams struct {
	Location       string
	PartitionCount int32
	ReplicaCount   int32
	Sku            v1.SearchSku
}

func (s *SearchOp) Create(ctx context.Context, params SearchCreateParams) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("Search.Create", func() (any, error) {
		return s.Client.CreateAZ0105(ctx, &v1.SearchPostRequestBody{
			Location:       params.Location,
			PartitionCount: params.PartitionCount,
			ReplicaCount:   params.ReplicaCount,
			Sku:            params.Sku,
		})
	})
}

func (s *SearchOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0105NotFound]("Search.Read", func() (any, error) {
		return s.Client.GetAZ0105(ctx, v1.GetAZ0105Params{
			ResourceGroupName: id,
		})
	})
}

func (s *SearchOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0105NoContent, v1.DeleteAZ0105NotFound]("Search.Delete", func() (any, error) {
		return s.Client.DeleteAZ0105(ctx, v1.DeleteAZ0105Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (s *SearchOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("Search.Status", func() (any, error) {
		return s.Client.StatusAZ0105(ctx, v1.StatusAZ0105Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ SearchAPI = (*SearchOp)(nil)
