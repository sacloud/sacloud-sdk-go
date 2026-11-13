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

type DWHAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type DWHOp struct{ *v1.Client }

func NewDWHOp(client *v1.Client) DWHAPI { return &DWHOp{Client: client} }

func (d *DWHOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListDwhResourcesNotFound]("DWH.List", func() (any, error) {
		return d.Client.ListDwhResources(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}
	return
}

func (d *DWHOp) Create(ctx context.Context, location string) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("DWH.Create", func() (any, error) {
		return d.Client.CreateDwhResource(ctx, &v1.DatawarehousePostRequestBody{
			Location: v1.DatawarehousePostRequestBodyLocation(location),
		})
	})
}

func (d *DWHOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.ReadDwhResourceNotFound]("DWH.Read", func() (any, error) {
		return d.Client.ReadDwhResource(ctx, v1.ReadDwhResourceParams{
			ResourceGroupName: id,
		})
	})
}

func (d *DWHOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteDwhResourceNoContent, v1.DeleteDwhResourceNotFound]("DWH.Delete", func() (any, error) {
		return d.Client.DeleteDwhResource(ctx, v1.DeleteDwhResourceParams{
			ResourceGroupName: id,
		})
	})
	return
}

func (d *DWHOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("DWH.Status", func() (any, error) {
		return d.Client.ReadDwhResourceStatus(ctx, v1.ReadDwhResourceStatusParams{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ DWHAPI = (*DWHOp)(nil)
