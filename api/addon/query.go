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

type QueryAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type QueryOp struct{ *v1.Client }

func NewQueryOp(client *v1.Client) QueryAPI { return &QueryOp{Client: client} }

func (q *QueryOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListQueryResourcesNotFound]("Query.List", func() (any, error) {
		return q.Client.ListQueryResources(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

func (q *QueryOp) Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("Query.Create", func() (any, error) {
		return q.Client.CreateQueryResource(ctx, &v1.QueryPostRequestBody{
			Location: v1.QueryPostRequestBodyLocation(location),
		})
	})
}

func (q *QueryOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.ReadQueryResourceNotFound]("Query.Read", func() (any, error) {
		return q.Client.ReadQueryResource(ctx, v1.ReadQueryResourceParams{
			ResourceGroupName: id,
		})
	})
}

func (q *QueryOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteQueryResourceNoContent, v1.DeleteQueryResourceNotFound]("Query.Delete", func() (any, error) {
		return q.Client.DeleteQueryResource(ctx, v1.DeleteQueryResourceParams{
			ResourceGroupName: id,
		})
	})

	return
}

func (q *QueryOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("Query.Status", func() (any, error) {
		return q.Client.ReadQueryResourceStatus(ctx, v1.ReadQueryResourceStatusParams{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ QueryAPI = (*QueryOp)(nil)
