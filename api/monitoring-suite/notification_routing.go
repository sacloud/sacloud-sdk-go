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
	"strconv"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type NotificationRoutingAPI interface {
	List(ctx context.Context, projectId string, count, from *int) ([]v1.NotificationRouting, error)
	Create(ctx context.Context, projectId string, params NotificationRoutingCreateParams) (*v1.NotificationRouting, error)
	Read(ctx context.Context, projectId string, id uuid.UUID) (*v1.NotificationRouting, error)
	Update(ctx context.Context, projectId string, id uuid.UUID, params NotificationRoutingUpdateParams) (*v1.NotificationRouting, error)
	Delete(ctx context.Context, projectId string, id uuid.UUID) error

	Reorder(ctx context.Context, projectId string, orders []v1.NotificationRoutingOrderRequest) error
}

var _ NotificationRoutingAPI = (*notificationRoutingOp)(nil)

type notificationRoutingOp struct {
	client *v1.Client
}

func NewNotificationRoutingOp(client *v1.Client) NotificationRoutingAPI {
	return &notificationRoutingOp{client: client}
}

func (op *notificationRoutingOp) List(ctx context.Context, projectId string, count, from *int) (ret []v1.NotificationRouting, err error) {
	res, err := errorFromDecodedResponse("NotificationRouting.List", func() (*v1.PaginatedNotificationRoutingList, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationRoutingsList(ctx, v1.AlertsProjectsNotificationRoutingsListParams{
			ProjectResourceID: id,
			Count:             intoOpt[v1.OptInt](count),
			From:              intoOpt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

type NotificationRoutingCreateParams struct {
	// Fields based on v1.NotificationRouting
	NotificationTargetUID uuid.UUID
	MatchLabels           []v1.MatchLabelsItem
	ResendIntervalMinutes *int
}

func (op *notificationRoutingOp) Create(ctx context.Context, projectId string, params NotificationRoutingCreateParams) (*v1.NotificationRouting, error) {
	res, err := errorFromDecodedResponse("NotificationRouting.Create", func() (*v1.NotificationRouting, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		req := v1.NotificationRoutingRequest{
			NotificationTargetUID: params.NotificationTargetUID,
			MatchLabels: func() []v1.MatchLabelsItemRequest {
				if len(params.MatchLabels) == 0 {
					return nil
				}
				ret := make([]v1.MatchLabelsItemRequest, len(params.MatchLabels))
				for i, item := range params.MatchLabels {
					ret[i] = v1.MatchLabelsItemRequest{
						Name:  item.Name,
						Value: item.Value,
					}
				}
				return ret
			}(),
			ResendIntervalMinutes: intoOpt[v1.OptInt](params.ResendIntervalMinutes),
		}
		return op.client.AlertsProjectsNotificationRoutingsCreate(ctx, &req, v1.AlertsProjectsNotificationRoutingsCreateParams{ProjectResourceID: id})
	})
	return unwrapE[*v1.NotificationRouting](res, err)
}

type NotificationRoutingUpdateParams struct {
	// Fields based on v1.PatchedNotificationRouting
	NotificationTargetUID *uuid.UUID
	MatchLabels           []v1.MatchLabelsItem
	ResendIntervalMinutes *int
}

func (op *notificationRoutingOp) Update(ctx context.Context, projectId string, uid uuid.UUID, params NotificationRoutingUpdateParams) (*v1.NotificationRouting, error) {
	return errorFromDecodedResponse("NotificationRouting.Update", func() (*v1.NotificationRouting, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationRoutingsPartialUpdate(ctx, v1.NewOptPatchedNotificationRoutingRequest(v1.PatchedNotificationRoutingRequest{
			NotificationTargetUID: intoOpt[v1.OptUUID](params.NotificationTargetUID),
			MatchLabels: func() []v1.MatchLabelsItemRequest {
				if len(params.MatchLabels) == 0 {
					return nil
				}
				ret := make([]v1.MatchLabelsItemRequest, len(params.MatchLabels))
				for i, item := range params.MatchLabels {
					ret[i] = v1.MatchLabelsItemRequest{
						Name:  item.Name,
						Value: item.Value,
					}
				}
				return ret
			}(),
			ResendIntervalMinutes: intoOpt[v1.OptInt](params.ResendIntervalMinutes),
		}), v1.AlertsProjectsNotificationRoutingsPartialUpdateParams{
			ProjectResourceID: id,
			UID:               uid,
		})
	})
}

func (op *notificationRoutingOp) Read(ctx context.Context, projectId string, uid uuid.UUID) (*v1.NotificationRouting, error) {
	return errorFromDecodedResponse("NotificationRouting.Read", func() (*v1.NotificationRouting, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationRoutingsRetrieve(ctx, v1.AlertsProjectsNotificationRoutingsRetrieveParams{
			ProjectResourceID: id,
			UID:               uid,
		})
	})
}

func (op *notificationRoutingOp) Delete(ctx context.Context, projectId string, uid uuid.UUID) error {
	return errorFromDecodedResponse1("NotificationRouting.Delete", func() error {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.AlertsProjectsNotificationRoutingsDestroy(ctx, v1.AlertsProjectsNotificationRoutingsDestroyParams{
			ProjectResourceID: id,
			UID:               uid,
		})
	})
}

func (op *notificationRoutingOp) Reorder(ctx context.Context, projectId string, orders []v1.NotificationRoutingOrderRequest) error {
	return errorFromDecodedResponse1("NotificationRouting.Reorder", func() error {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.AlertsProjectsNotificationRoutingsReorderUpdate(ctx, orders, v1.AlertsProjectsNotificationRoutingsReorderUpdateParams{
			ProjectResourceID: id,
		})
	})
}
