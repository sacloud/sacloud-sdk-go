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
	"strconv"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
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

func (op *dashboardProjectOp) List(ctx context.Context, count *int, from *int) ([]v1.DashboardProject, error) {
	resp, err := op.client.DashboardsProjectsList(ctx, v1.DashboardsProjectsListParams{
		Count: intoOpt[v1.OptInt](count),
		From:  intoOpt[v1.OptInt](from),
	})
	if err != nil {
		return nil, NewAPIError("DashboardProject.List", 0, err)
	}
	return resp.Results, nil
}

type DashboardProjectCreateParams struct {
	Name        string
	Description *string
}

func (op *dashboardProjectOp) Create(ctx context.Context, p DashboardProjectCreateParams) (*v1.DashboardProject, error) {
	request := v1.DashboardProjectCreate{
		Name:        p.Name,
		Description: intoOpt[v1.OptString](p.Description),
	}
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

func (op *dashboardProjectOp) Read(ctx context.Context, id string) (*v1.DashboardProject, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewError("DashboardProject.Read", err)
	}
	resp, err := op.client.DashboardsProjectsRetrieve(ctx, v1.DashboardsProjectsRetrieveParams{ResourceID: intId})
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

type DashboardProjectUpdateParams struct {
	Name        *string
	Description *string
}

func (op *dashboardProjectOp) Update(ctx context.Context, id string, params DashboardProjectUpdateParams) (*v1.DashboardProject, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewError("DashboardProject.Update", err)
	}
	request := v1.PatchedDashboardProject{
		Name:        intoOpt[v1.OptString](params.Name),
		Description: intoOpt[v1.OptString](params.Description),
	}
	req := v1.NewOptPatchedDashboardProject(request)
	resp, err := op.client.DashboardsProjectsPartialUpdate(ctx, req, v1.DashboardsProjectsPartialUpdateParams{ResourceID: intId})
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

func (op *dashboardProjectOp) Delete(ctx context.Context, id string) error {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewError("DashboardProject.Delete", err)
	}
	err = op.client.DashboardsProjectsDestroy(ctx, v1.DashboardsProjectsDestroyParams{ResourceID: intId})
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
