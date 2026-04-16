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
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0104NotFound]("Query.List", func() (any, error) {
		return q.Client.ListAZ0104(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

func (q *QueryOp) Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("Query.Create", func() (any, error) {
		return q.Client.CreateAZ0104(ctx, &v1.QueryPostRequestBody{
			Location: location,
		})
	})
}

func (q *QueryOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0104NotFound]("Query.Read", func() (any, error) {
		return q.Client.GetAZ0104(ctx, v1.GetAZ0104Params{
			ResourceGroupName: id,
		})
	})
}

func (q *QueryOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0104NoContent, v1.DeleteAZ0104NotFound]("Query.Delete", func() (any, error) {
		return q.Client.DeleteAZ0104(ctx, v1.DeleteAZ0104Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (q *QueryOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("Query.Status", func() (any, error) {
		return q.Client.StatusAZ0104(ctx, v1.StatusAZ0104Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ QueryAPI = (*QueryOp)(nil)
