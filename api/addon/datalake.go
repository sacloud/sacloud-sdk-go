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

type DataLakeAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, params DataLakeCreateParams) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error
	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type DataLakeOp struct{ *v1.Client }

func NewDataLakeOp(client *v1.Client) DataLakeAPI { return &DataLakeOp{Client: client} }

func (d *DataLakeOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0101NotFound]("DataLake.List", func() (any, error) {
		return d.Client.ListAZ0101(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}
	return
}

type DataLakeCreateParams struct {
	Location    string
	Performance v1.DataLakePerformance
	Redundancy  v1.DataLakeRedundancy
}

func (d *DataLakeOp) Create(ctx context.Context, params DataLakeCreateParams) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("DataLake.Create", func() (any, error) {
		return d.Client.CreateAZ0101(ctx, &v1.DatalakePostRequestBody{
			Location:    params.Location,
			Performance: params.Performance,
			Redundancy:  params.Redundancy,
		})
	})
}

func (d *DataLakeOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0101NotFound]("DataLake.Read", func() (any, error) {
		return d.Client.GetAZ0101(ctx, v1.GetAZ0101Params{
			ResourceGroupName: id,
		})
	})
}

func (d *DataLakeOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0101NoContent, v1.DeleteAZ0101NotFound]("DataLake.Delete", func() (any, error) {
		return d.Client.DeleteAZ0101(ctx, v1.DeleteAZ0101Params{
			ResourceGroupName: id,
		})
	})
	return
}

func (d *DataLakeOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("DataLake.Status", func() (any, error) {
		return d.Client.StatusAZ0101(ctx, v1.StatusAZ0101Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ DataLakeAPI = (*DataLakeOp)(nil)
