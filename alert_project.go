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
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type AlertProjectAPI interface {
	List(ctx context.Context, count *int, from *int) ([]v1.AlertProject, error)
	Create(ctx context.Context, params AlertProjectCreateParams) (*v1.AlertProject, error)
	Read(ctx context.Context, id string) (*v1.WrappedAlertProject, error)
	Update(ctx context.Context, id string, request AlertProjectUpdateParams) (*v1.WrappedAlertProject, error)
	Delete(ctx context.Context, id string) error

	ListHistories(ctx context.Context, params AlertsProjectsHistoriesListParams) ([]v1.History, error)
	ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error)
}

var _ AlertProjectAPI = (*alertProjectOp)(nil)

type alertProjectOp struct {
	client *v1.Client
}

func NewAlertProjectOp(client *v1.Client) AlertProjectAPI {
	return &alertProjectOp{client: client}
}

func (op *alertProjectOp) List(ctx context.Context, count *int, from *int) ([]v1.AlertProject, error) {
	result, err := op.client.AlertsProjectsList(ctx, v1.AlertsProjectsListParams{
		Count: intoOpt[v1.OptInt](count),
		From:  intoOpt[v1.OptInt](from),
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
	p := v1.AlertsProjectsRetrieveParams{ResourceID: intId}
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

type AlertProjectCreateParams struct {
	Name        string
	Description *string
}

func (op *alertProjectOp) Create(ctx context.Context, params AlertProjectCreateParams) (*v1.AlertProject, error) {
	result, err := op.client.AlertsProjectsCreate(ctx, &v1.AlertProjectCreate{
		Name:        params.Name,
		Description: intoOpt[v1.OptString](params.Description),
	})
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

type AlertProjectUpdateParams struct {
	Name        *string
	Description *string
}

func (op *alertProjectOp) Update(ctx context.Context, id string, params AlertProjectUpdateParams) (*v1.WrappedAlertProject, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.Update", 0, err)
	}
	p := v1.AlertsProjectsPartialUpdateParams{ResourceID: intId}
	body := v1.NewOptPatchedAlertProject(v1.PatchedAlertProject{
		Name:        intoOpt[v1.OptString](params.Name),
		Description: intoOpt[v1.OptString](params.Description),
		// other fields are golang's zero values, meaning "not set"
	})
	result, err := op.client.AlertsProjectsPartialUpdate(ctx, body, p)
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
	p := v1.AlertsProjectsDestroyParams{ResourceID: intId}
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

type AlertsProjectsHistoriesListParams struct {
	ProjectID string // mandatory
	Count     *int
	From      *int
	Open      *bool
	Severity  *v1.AlertsProjectsHistoriesListSeverity
	StartsAt  *time.Time
}

func (op *alertProjectOp) ListHistories(ctx context.Context, params AlertsProjectsHistoriesListParams) ([]v1.History, error) {
	intProjectId, err := strconv.ParseInt(params.ProjectID, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.ListHistories", 0, err)
	}
	p := v1.AlertsProjectsHistoriesListParams{
		ProjectResourceID: intProjectId,
		Count:             intoOpt[v1.OptInt](params.Count),
		From:              intoOpt[v1.OptInt](params.From),
		Open:              intoOpt[v1.OptBool](params.Open),
		Severity:          intoOpt[v1.OptAlertsProjectsHistoriesListSeverity](params.Severity),
		StartsAt:          intoOpt[v1.OptDateTime](params.StartsAt),
	}
	result, err := op.client.AlertsProjectsHistoriesList(ctx, p)
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

func (op *alertProjectOp) ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertProject.ReadHistory", 0, err)
	}
	p := v1.AlertsProjectsHistoriesRetrieveParams{
		UID:               historyId,
		ProjectResourceID: intProjectId,
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
