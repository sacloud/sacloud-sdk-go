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
	params "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type AlertProjectAPI interface {
	List(ctx context.Context, count int, from int) ([]v1.AlertProject, error)
	Create(ctx context.Context, params v1.AlertProjectCreate) (*v1.AlertProject, error)
	Read(ctx context.Context, id string) (*v1.WrappedAlertProject, error)
	Update(ctx context.Context, id string, request *v1.AlertProject) (*v1.WrappedAlertProject, error)
	Delete(ctx context.Context, id string) error

	ListHistories(ctx context.Context, params v1.AlertsProjectsHistoriesListParams) ([]v1.History, error)
	ReadHistory(ctx context.Context, projectId string, historyId string) (*v1.History, error)
}

var _ AlertProjectAPI = (*alertProjectOp)(nil)

type alertProjectOp struct {
	client *v1.Client
}

func NewAlertProjectOp(client *v1.Client) AlertProjectAPI {
	return &alertProjectOp{client: client}
}

func (op *alertProjectOp) List(ctx context.Context, count int, from int) ([]v1.AlertProject, error) {
	result, err := op.client.AlertsProjectsList(ctx, params.AlertsProjectsListParams{
		Count: v1.NewOptInt(count),
		From:  v1.NewOptInt(from),
	})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("AlertProject.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *alertProjectOp) Read(ctx context.Context, id string) (*v1.WrappedAlertProject, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.Read", 0, err)
	}
	p := params.AlertsProjectsRetrieveParams{ID: intId}
	result, err := op.client.AlertsProjectsRetrieve(ctx, p)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertProject.Read", e.StatusCode, errors.Wrap(err, "alert project not found"))
		default:
			return nil, NewAPIError("AlertProject.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.Read", 0, err)
	} else {
		ret := new(v1.WrappedAlertProject)
		return Unwrap(ret, result)
	}
}

func (op *alertProjectOp) Create(ctx context.Context, params v1.AlertProjectCreate) (*v1.AlertProject, error) {
	result, err := op.client.AlertsProjectsCreate(ctx, &params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertProject.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new alert project"))
		default:
			return nil, NewAPIError("AlertProject.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *alertProjectOp) Update(ctx context.Context, id string, resource *v1.AlertProject) (*v1.WrappedAlertProject, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.Update", 0, err)
	}
	p := params.AlertsProjectsUpdateParams{ID: intId}
	body := v1.NewOptAlertProject(*resource)
	result, err := op.client.AlertsProjectsUpdate(ctx, body, p)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertProject.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("AlertProject.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.Update", 0, err)
	} else {
		ret := new(v1.WrappedAlertProject)
		return Unwrap(ret, result)
	}
}

func (op *alertProjectOp) Delete(ctx context.Context, id string) error {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("AlertProject.Delete", 0, err)
	}
	p := params.AlertsProjectsDestroyParams{ID: intId}
	err = op.client.AlertsProjectsDestroy(ctx, p)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("AlertProject.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("AlertProject.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("AlertProject.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("AlertProject.Delete", 0, err)
	}
	return nil
}

func (op *alertProjectOp) ListHistories(ctx context.Context, params v1.AlertsProjectsHistoriesListParams) ([]v1.History, error) {
	result, err := op.client.AlertsProjectsHistoriesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.ListHistories", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertProject.ListHistories", e.StatusCode, errors.Wrap(err, "project not found"))
		default:
			return nil, NewAPIError("AlertProject.ListHistories", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.ListHistories", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *alertProjectOp) ReadHistory(ctx context.Context, projectId string, historyId string) (*v1.History, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.ReadHistory", 0, err)
	}
	intHistoryId, err := strconv.ParseInt(historyId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.ReadHistory", 0, err)
	}
	p := v1.AlertsProjectsHistoriesRetrieveParams{
		ID:        int(intHistoryId),
		ProjectPk: int(intProjectId),
	}
	result, err := op.client.AlertsProjectsHistoriesRetrieve(ctx, p)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertProject.ReadHistory", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertProject.ReadHistory", e.StatusCode, errors.Wrap(err, "alert history not found"))
		default:
			return nil, NewAPIError("AlertProject.ReadHistory", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertProject.ReadHistory", 0, err)
	} else {
		return result, nil
	}
}
