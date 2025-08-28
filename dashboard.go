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
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type DashboardProjectAPI interface {
	List(ctx context.Context, count int, from int) ([]v1.DashboardProject, error)
	Create(ctx context.Context, request v1.DashboardProjectCreate) (*v1.DashboardProject, error)
	Read(ctx context.Context, id int64) (*v1.DashboardProject, error)
	Update(ctx context.Context, id int64, request *v1.DashboardProject) (*v1.DashboardProject, error)
	Delete(ctx context.Context, id int64) error
}

var _ DashboardProjectAPI = (*dashboardOp)(nil)

type dashboardOp struct {
	client *v1.Client
}

func NewDashboardOp(client *v1.Client) DashboardProjectAPI {
	return &dashboardOp{client: client}
}

func (op *dashboardOp) List(ctx context.Context, count int, from int) ([]v1.DashboardProject, error) {

	resp, err := op.client.DashboardsProjectsList(ctx, v1.DashboardsProjectsListParams{
		Count: v1.NewOptInt(count),
		From:  v1.NewOptInt(from),
	})
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (op *dashboardOp) Create(ctx context.Context, request v1.DashboardProjectCreate) (*v1.DashboardProject, error) {
	resp, err := op.client.DashboardsProjectsCreate(ctx, &request)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("DashboardProject.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("DashboardProject.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("DashboardProject.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("DashboardProject.Create", 0, err)
	} else {
		return resp, nil
	}
}

func (op *dashboardOp) Read(ctx context.Context, id int64) (*v1.DashboardProject, error) {
	resp, err := op.client.DashboardsProjectsRetrieve(ctx, v1.DashboardsProjectsRetrieveParams{ID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("DashboardProject.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("DashboardProject.Read", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("DashboardProject.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("DashboardProject.Read", 0, err)
	} else {
		ret := new(v1.DashboardProject)
		return Unwrap(ret, resp)
	}
}
func (op *dashboardOp) Update(ctx context.Context, id int64, request *v1.DashboardProject) (*v1.DashboardProject, error) {
	req := v1.NewOptDashboardProject(*request)
	resp, err := op.client.DashboardsProjectsUpdate(ctx, req, v1.DashboardsProjectsUpdateParams{ID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("DashboardProject.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("DashboardProject.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("DashboardProject.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("DashboardProject.Update", 0, err)
	} else {
		ret := new(v1.DashboardProject)
		return Unwrap(ret, resp)
	}
}

func (op *dashboardOp) Delete(ctx context.Context, id int64) error {
	err := op.client.DashboardsProjectsDestroy(ctx, v1.DashboardsProjectsDestroyParams{ID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("DashboardProject.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("DashboardProject.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("DashboardProject.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("DashboardProject.Delete", 0, err)
	}
	return nil
}
