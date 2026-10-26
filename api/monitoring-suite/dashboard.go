// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"context"
	"strconv"

	v1 "github.com/sacloud/sacloud-sdk-go/api/monitoring-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type DashboardProjectAPI interface {
	List(ctx context.Context, count *int, from *int) ([]v1.DashboardProject, error)
	Create(ctx context.Context, request DashboardProjectCreateParams) (*v1.DashboardProject, error)
	Read(ctx context.Context, id string) (*v1.DashboardProject, error)
	Update(ctx context.Context, id string, params DashboardProjectUpdateParams) (*v1.DashboardProject, error)
	Delete(ctx context.Context, id string) error
}

var _ DashboardProjectAPI = (*dashboardProjectOp)(nil)

type dashboardProjectOp struct {
	client *v1.Client
}

func NewDashboardOp(client *v1.Client) DashboardProjectAPI {
	return &dashboardProjectOp{client: client}
}

func (op *dashboardProjectOp) List(ctx context.Context, count *int, from *int) (ret []v1.DashboardProject, err error) {
	res, err := errorFromDecodedResponse("DashboardProject.List", func() (*v1.PaginatedDashboardProjectList, error) {
		return op.client.ListDashboardProjects(ctx, v1.ListDashboardProjectsParams{
			Count: into.Opt[v1.OptInt](count),
			From:  into.Opt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

type DashboardProjectCreateParams struct {
	Name        string
	Description *string
}

func (op *dashboardProjectOp) Create(ctx context.Context, p DashboardProjectCreateParams) (*v1.DashboardProject, error) {
	return errorFromDecodedResponse("DashboardProject.Create", func() (*v1.DashboardProject, error) {
		return op.client.CreateDashboardProject(ctx, &v1.DashboardProjectCreateRequest{
			Name:        p.Name,
			Description: into.Opt[v1.OptString](p.Description),
		})
	})
}

func (op *dashboardProjectOp) Read(ctx context.Context, id string) (*v1.DashboardProject, error) {
	res, err := errorFromDecodedResponse("DashboardProject.Read", func() (*v1.WrappedDashboardProject, error) {
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.ReadDashboardProject(ctx, v1.ReadDashboardProjectParams{ResourceId: intId})
	})
	return unwrapE[*v1.DashboardProject](res, err)
}

type DashboardProjectUpdateParams struct {
	Name        *string
	Description *string
}

func (op *dashboardProjectOp) Update(ctx context.Context, id string, params DashboardProjectUpdateParams) (*v1.DashboardProject, error) {
	res, err := errorFromDecodedResponse("DashboardProject.Update", func() (*v1.WrappedDashboardProject, error) {
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.PatchDashboardProject(ctx, v1.NewOptPatchedDashboardProjectRequest(v1.PatchedDashboardProjectRequest{
			Name:        into.Opt[v1.OptString](params.Name),
			Description: into.Opt[v1.OptString](params.Description),
		}), v1.PatchDashboardProjectParams{ResourceId: intId})
	})
	return unwrapE[*v1.DashboardProject](res, err)
}

func (op *dashboardProjectOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponse1("DashboardProject.Delete", func() error {
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		return op.client.DeleteDashboardProject(ctx, v1.DeleteDashboardProjectParams{
			ResourceId: intId,
		})
	})
}
