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

type ETLAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type ETLOp struct{ *v1.Client }

func NewETLOp(client *v1.Client) ETLAPI { return &ETLOp{Client: client} }

func (e *ETLOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0103NotFound]("ETL.List", func() (any, error) {
		return e.Client.ListAZ0103(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

func (e *ETLOp) Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("ETL.Create", func() (any, error) {
		return e.Client.CreateAZ0103(ctx, &v1.EtlPostRequestBody{
			Location: location,
		})
	})
}

func (e *ETLOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0103NotFound]("ETL.Read", func() (any, error) {
		return e.Client.GetAZ0103(ctx, v1.GetAZ0103Params{
			ResourceGroupName: id,
		})
	})
}

func (e *ETLOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0103NoContent, v1.DeleteAZ0103NotFound]("ETL.Delete", func() (any, error) {
		return e.Client.DeleteAZ0103(ctx, v1.DeleteAZ0103Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (e *ETLOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("ETL.Status", func() (any, error) {
		return e.Client.StatusAZ0103(ctx, v1.StatusAZ0103Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ ETLAPI = (*ETLOp)(nil)
