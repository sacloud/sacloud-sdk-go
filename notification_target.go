// Copyright 2025- The sacloud/monitoring-suite-api-go Contributors
//
// This software is licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You can obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is provided "AS IS", without warranties or conditions of any kind,
// either express or implied. See the License for specific language governing permissions and
// limitations under the License.

package monitoringsuite

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type NotificationTargetAPI interface {
	List(ctx context.Context, projectId string, params NotificationTargetsListParams) ([]v1.NotificationTarget, error)
	Create(ctx context.Context, projectId string, params NotificationTargetCreateParams) (*v1.NotificationTarget, error)
	Read(ctx context.Context, projectId string, id uuid.UUID) (*v1.NotificationTarget, error)
	Update(ctx context.Context, projectId string, id uuid.UUID, params NotificationTargetUpdateParams) (*v1.NotificationTarget, error)
	Delete(ctx context.Context, projectId string, id uuid.UUID) error
}

var _ NotificationTargetAPI = (*notificationTargetOp)(nil)

type notificationTargetOp struct {
	client *v1.Client
}

func NewNotificationTargetOp(client *v1.Client) NotificationTargetAPI {
	return &notificationTargetOp{client: client}
}

type NotificationTargetsListParams struct {
	Count *int
	From  *int
}

func (op *notificationTargetOp) List(ctx context.Context, projectId string, p NotificationTargetsListParams) ([]v1.NotificationTarget, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationTarget.List", err)
	}
	params := v1.AlertsProjectsNotificationTargetsListParams{
		ProjectResourceID: id,
		Count:             intoOpt[v1.OptInt](p.Count),
		From:              intoOpt[v1.OptInt](p.From),
	}
	result, err := op.client.AlertsProjectsNotificationTargetsList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationTarget.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("NotificationTarget.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationTarget.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *notificationTargetOp) Read(ctx context.Context, projectId string, uid uuid.UUID) (*v1.NotificationTarget, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationTarget.Read", err)
	}
	params := v1.AlertsProjectsNotificationTargetsRetrieveParams{
		ProjectResourceID: pid,
		UID:               uid,
	}
	result, err := op.client.AlertsProjectsNotificationTargetsRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationTarget.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("NotificationTarget.Read", e.StatusCode, errors.Wrap(err, "notification target not found"))
		default:
			return nil, NewAPIError("NotificationTarget.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationTarget.Read", 0, err)
	} else {
		ret := new(v1.NotificationTarget)
		return Unwrap(ret, result)
	}
}

type NotificationTargetCreateParams struct {
	ServiceType v1.NotificationTargetServiceType
	URL         url.URL
	Description *string
}

func (op *notificationTargetOp) Create(ctx context.Context, projectId string, params NotificationTargetCreateParams) (*v1.NotificationTarget, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationTarget.Create", err)
	}
	createParams := v1.AlertsProjectsNotificationTargetsCreateParams{ProjectResourceID: pid}
	req := v1.NotificationTarget{
		ServiceType: params.ServiceType,
		URL:         params.URL.String(),
		Description: intoOpt[v1.OptString](params.Description),
	}
	result, err := op.client.AlertsProjectsNotificationTargetsCreate(ctx, &req, createParams)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationTarget.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("NotificationTarget.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new target"))
		default:
			return nil, NewAPIError("NotificationTarget.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationTarget.Create", 0, err)
	} else {
		return result, nil
	}
}

type NotificationTargetUpdateParams struct {
	ServiceType *v1.PatchedNotificationTargetServiceType
	URL         *string
	Description *string
}

func (op *notificationTargetOp) Update(ctx context.Context, projectId string, uid uuid.UUID, params NotificationTargetUpdateParams) (*v1.NotificationTarget, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationTarget.Update", err)
	}
	req := v1.PatchedNotificationTarget{
		ServiceType: intoOpt[v1.OptPatchedNotificationTargetServiceType](params.ServiceType),
		URL:         intoOpt[v1.OptString](params.URL),
		Description: intoOpt[v1.OptString](params.Description),
	}
	updateParams := v1.AlertsProjectsNotificationTargetsPartialUpdateParams{
		ProjectResourceID: pid,
		UID:               uid,
	}
	result, err := op.client.AlertsProjectsNotificationTargetsPartialUpdate(ctx, v1.NewOptPatchedNotificationTarget(req), updateParams)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationTarget.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("NotificationTarget.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("NotificationTarget.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationTarget.Update", 0, err)
	} else {
		ret := new(v1.NotificationTarget)
		return Unwrap(ret, result)
	}
}

func (op *notificationTargetOp) Delete(ctx context.Context, projectId string, uid uuid.UUID) error {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewError("NotificationTarget.Delete", err)
	}
	params := v1.AlertsProjectsNotificationTargetsDestroyParams{
		ProjectResourceID: pid,
		UID:               uid,
	}
	err = op.client.AlertsProjectsNotificationTargetsDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("NotificationTarget.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("NotificationTarget.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("NotificationTarget.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("NotificationTarget.Delete", 0, err)
	}
	return nil
}
