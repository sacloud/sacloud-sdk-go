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
	"strconv"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type NotificationTargetAPI interface {
	List(ctx context.Context, params v1.AlertsProjectsNotificationTargetsListParams) ([]v1.NotificationTarget, error)
	Create(ctx context.Context, params v1.NotificationTarget) (*v1.NotificationTarget, error)
	Read(ctx context.Context, id string) (*v1.NotificationTarget, error)
	Update(ctx context.Context, id string, request *v1.NotificationTarget) (*v1.NotificationTarget, error)
	Delete(ctx context.Context, id string) error
}

var _ NotificationTargetAPI = (*notificationTargetOp)(nil)

type notificationTargetOp struct {
	client *v1.Client
}

func NewNotificationTargetOp(client *v1.Client) NotificationTargetAPI {
	return &notificationTargetOp{client: client}
}

func (op *notificationTargetOp) List(ctx context.Context, params v1.AlertsProjectsNotificationTargetsListParams) ([]v1.NotificationTarget, error) {
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

func (op *notificationTargetOp) Read(ctx context.Context, id string) (*v1.NotificationTarget, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("NotificationTarget.Read", 0, err)
	}
	// :TODO: AlertsProjectsNotificationTargetsRetrieveParams() taking int instead of int64 can be subject to change
	params := v1.AlertsProjectsNotificationTargetsRetrieveParams{ID: int(intId)}
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

func (op *notificationTargetOp) Create(ctx context.Context, params v1.NotificationTarget) (*v1.NotificationTarget, error) {
	// project_pk is required for creation, extract from params.ProjectID
	createParams := v1.AlertsProjectsNotificationTargetsCreateParams{ProjectPk: params.ProjectID}
	result, err := op.client.AlertsProjectsNotificationTargetsCreate(ctx, &params, createParams)
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

func (op *notificationTargetOp) Update(ctx context.Context, id string, resource *v1.NotificationTarget) (*v1.NotificationTarget, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("NotificationTarget.Update", 0, err)
	}
	// :TODO: AlertsProjectsNotificationTargetsUpdateParams() taking int instead of int64 can be subject to change
	params := v1.AlertsProjectsNotificationTargetsUpdateParams{ID: int(intId)}
	result, err := op.client.AlertsProjectsNotificationTargetsUpdate(ctx, resource, params)
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

func (op *notificationTargetOp) Delete(ctx context.Context, id string) error {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("NotificationTarget.Delete", 0, err)
	}
	// :TODO: AlertsProjectsNotificationTargetsDestroyParams() taking int instead of int64 can be subject to change
	params := v1.AlertsProjectsNotificationTargetsDestroyParams{ID: int(intId)}
	err = op.client.AlertsProjectsNotificationTargetsDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("NotificationTarget.Remove", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("NotificationTarget.Remove", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("NotificationTarget.Remove", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("NotificationTarget.Remove", 0, err)
	}
	return nil
}
