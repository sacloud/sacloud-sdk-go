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

type StreamingAPI interface {
	List(ctx context.Context) ([]v1.ResourceGroupResource, error)
	Create(ctx context.Context, location, unitCount string) (*v1.PostDeploymentResponse, error)
	Read(ctx context.Context, id string) (*v1.GetResourceResponse, error)
	Delete(ctx context.Context, id string) error

	Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error)
}

type StreamingOp struct{ *v1.Client }

func NewStreamingOp(client *v1.Client) StreamingAPI { return &StreamingOp{Client: client} }

func (s *StreamingOp) List(ctx context.Context) (ret []v1.ResourceGroupResource, err error) {
	if res, e := ErrorFromDecodedResponse[v1.ListResourcesResponse, v1.ListAZ0107NotFound]("Streaming.List", func() (any, error) {
		return s.Client.ListAZ0107(ctx)
	}); e != nil {
		err = e
	} else if resources, ok := res.GetResources().Get(); ok {
		ret = resources
	}

	return
}

func (s *StreamingOp) Create(ctx context.Context, location, unitCount string) (*v1.PostDeploymentResponse, error) {
	return ErrorFromDecodedResponse[v1.PostDeploymentResponse, void]("Streaming.Create", func() (any, error) {
		return s.Client.CreateAZ0107(ctx, &v1.StreamingRequestBody{
			Location:  location,
			UnitCount: unitCount,
		})
	})
}

func (s *StreamingOp) Read(ctx context.Context, id string) (*v1.GetResourceResponse, error) {
	return ErrorFromDecodedResponse[v1.GetResourceResponse, v1.GetAZ0107NotFound]("Streaming.Read", func() (any, error) {
		return s.Client.GetAZ0107(ctx, v1.GetAZ0107Params{
			ResourceGroupName: id,
		})
	})
}

func (s *StreamingOp) Delete(ctx context.Context, id string) (err error) {
	_, err = ErrorFromDecodedResponse[v1.DeleteAZ0107NoContent, v1.DeleteAZ0107NotFound]("Streaming.Delete", func() (any, error) {
		return s.Client.DeleteAZ0107(ctx, v1.DeleteAZ0107Params{
			ResourceGroupName: id,
		})
	})

	return
}

func (s *StreamingOp) Status(ctx context.Context, resourceGroupName, deploymentName string) (*v1.DeploymentStatus, error) {
	return ErrorFromDecodedResponse[v1.DeploymentStatus, void]("Streaming.Status", func() (any, error) {
		return s.Client.StatusAZ0107(ctx, v1.StatusAZ0107Params{
			ResourceGroupName: resourceGroupName,
			DeploymentName:    deploymentName,
		})
	})
}

var _ StreamingAPI = (*StreamingOp)(nil)
