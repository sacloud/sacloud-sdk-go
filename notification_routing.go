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
// either express or implied. See the License for the specific language governing permissions and
// limitations under the License.

package monitoringsuite

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type NotificationRoutingAPI interface {
	List(ctx context.Context, projectId string, count, from *int) ([]v1.NotificationRouting, error)
	Create(ctx context.Context, projectId string, params NotificationRoutingCreateParams) (*v1.NotificationRouting, error)
	Read(ctx context.Context, projectId string, id uuid.UUID) (*v1.NotificationRouting, error)
	Update(ctx context.Context, projectId string, id uuid.UUID, params NotificationRoutingUpdateParams) (*v1.NotificationRouting, error)
	Delete(ctx context.Context, projectId string, id uuid.UUID) error

	Reorder(ctx context.Context, projectId string, orders []v1.NotificationRoutingOrder) error
}

var _ NotificationRoutingAPI = (*notificationRoutingOp)(nil)

type notificationRoutingOp struct {
	client *v1.Client
}

func NewNotificationRoutingOp(client *v1.Client) NotificationRoutingAPI {
	return &notificationRoutingOp{client: client}
}

func (op *notificationRoutingOp) List(ctx context.Context, projectId string, count, from *int) ([]v1.NotificationRouting, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationRouting.List", err)
	}
	params := v1.AlertsProjectsNotificationRoutingsListParams{
		ProjectResourceID: id,
		Count:             intoOpt[v1.OptInt](count),
		From:              intoOpt[v1.OptInt](from),
	}
	result, err := op.client.AlertsProjectsNotificationRoutingsList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationRouting.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("NotificationRouting.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationRouting.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

type NotificationRoutingCreateParams struct {
	// Fields based on v1.NotificationRouting
	NotificationTargetUID uuid.UUID
	MatchLabels           []v1.MatchLabelsItem
	ResendIntervalMinutes *int
}

func (op *notificationRoutingOp) Create(ctx context.Context, projectId string, params NotificationRoutingCreateParams) (*v1.NotificationRouting, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationRouting.Create", err)
	}
	createParams := v1.AlertsProjectsNotificationRoutingsCreateParams{ProjectResourceID: id}
	req := v1.NotificationRouting{
		NotificationTargetUID: v1.NewOptUUID(params.NotificationTargetUID),
		MatchLabels:           params.MatchLabels,
		ResendIntervalMinutes: intoOpt[v1.OptInt](params.ResendIntervalMinutes),
	}

	// prevent ogen error (encoder is not accepting empty struct)
	req.NotificationTarget.SetFake()
	req.NotificationTarget.SetServiceType(v1.NotificationTargetServiceTypeSAKURASIMPLENOTICE)

	result, err := op.client.AlertsProjectsNotificationRoutingsCreate(ctx, &req, createParams)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationRouting.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("NotificationRouting.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new routing"))
		default:
			return nil, NewAPIError("NotificationRouting.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationRouting.Create", 0, err)
	} else {
		return result, nil
	}
}

type NotificationRoutingUpdateParams struct {
	// Fields based on v1.PatchedNotificationRouting
	NotificationTargetUID *uuid.UUID
	MatchLabels           []v1.MatchLabelsItem
	ResendIntervalMinutes *int
}

func (op *notificationRoutingOp) Update(ctx context.Context, projectId string, uid uuid.UUID, params NotificationRoutingUpdateParams) (*v1.NotificationRouting, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationRouting.Update", err)
	}
	updateParams := v1.AlertsProjectsNotificationRoutingsPartialUpdateParams{
		ProjectResourceID: id,
		UID:               uid,
	}
	req := v1.PatchedNotificationRouting{
		NotificationTargetUID: intoOpt[v1.OptUUID](params.NotificationTargetUID),
		MatchLabels:           params.MatchLabels,
		ResendIntervalMinutes: intoOpt[v1.OptInt](params.ResendIntervalMinutes),
	}
	result, err := op.client.AlertsProjectsNotificationRoutingsPartialUpdate(ctx, v1.NewOptPatchedNotificationRouting(req), updateParams)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationRouting.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("NotificationRouting.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("NotificationRouting.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationRouting.Update", 0, err)
	} else {
		ret := new(v1.NotificationRouting)
		return Unwrap(ret, result)
	}
}

func (op *notificationRoutingOp) Read(ctx context.Context, projectId string, uid uuid.UUID) (*v1.NotificationRouting, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewError("NotificationRouting.Read", err)
	}
	params := v1.AlertsProjectsNotificationRoutingsRetrieveParams{
		ProjectResourceID: id,
		UID:               uid,
	}
	result, err := op.client.AlertsProjectsNotificationRoutingsRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("NotificationRouting.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("NotificationRouting.Read", e.StatusCode, errors.Wrap(err, "notification routing not found"))
		default:
			return nil, NewAPIError("NotificationRouting.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("NotificationRouting.Read", 0, err)
	} else {
		ret := new(v1.NotificationRouting)
		return Unwrap(ret, result)
	}
}

func (op *notificationRoutingOp) Delete(ctx context.Context, projectId string, uid uuid.UUID) error {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewError("NotificationRouting.Delete", err)
	}
	params := v1.AlertsProjectsNotificationRoutingsDestroyParams{
		ProjectResourceID: id,
		UID:               uid,
	}
	err = op.client.AlertsProjectsNotificationRoutingsDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("NotificationRouting.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("NotificationRouting.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("NotificationRouting.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("NotificationRouting.Delete", 0, err)
	}
	return nil
}

func (op *notificationRoutingOp) Reorder(ctx context.Context, projectId string, orders []v1.NotificationRoutingOrder) error {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewError("NotificationRouting.Reorder", err)
	}
	params := v1.AlertsProjectsNotificationRoutingsReorderUpdateParams{
		ProjectResourceID: id,
	}
	err = op.client.AlertsProjectsNotificationRoutingsReorderUpdate(ctx, orders, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("NotificationRouting.Reorder", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("NotificationRouting.Reorder", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return NewAPIError("NotificationRouting.Reorder", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("NotificationRouting.Reorder", 0, err)
	}
	return nil
}
