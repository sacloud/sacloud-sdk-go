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

	v1 "github.com/sacloud/sacloud-sdk-go/api/addon/apis/v1"
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
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListSearchResourcesNotFound]("Search.List", func() (any, error) {
		return s.Client.ListSearchResources(ctx)
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
		return s.Client.CreateSearchResource(ctx, &v1.SearchPostRequestBody{
			Location:       v1.SearchPostRequestBodyLocation(params.Location),
			PartitionCount: v1.SearchPostRequestBodyPartitionCount(params.PartitionCount),
			ReplicaCount:   params.ReplicaCount,
			Sku:            params.Sku,
		})
	})
}

func (s *SearchOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.ReadSearchResourceNotFound]("Search.Read", func() (any, error) {
		return s.Client.ReadSearchResource(ctx, v1.ReadSearchResourceParams{
			ResourceGroupName: id,
		})
	})
}

func (s *SearchOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteSearchResourceNoContent, v1.DeleteSearchResourceNotFound]("Search.Delete", func() (any, error) {
		return s.Client.DeleteSearchResource(ctx, v1.DeleteSearchResourceParams{
			ResourceGroupName: id,
		})
	})

	return
}

func (s *SearchOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("Search.Status", func() (any, error) {
		return s.Client.ReadSearchResourceStatus(ctx, v1.ReadSearchResourceStatusParams{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ SearchAPI = (*SearchOp)(nil)
