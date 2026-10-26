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
	v1 "github.com/sacloud/sacloud-sdk-go/api/monitoring-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type NotificationRoutingAPI interface {
	List(ctx context.Context, projectId string, params NotificationRoutingsListParams) ([]v1.NotificationRouting, error)
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

type NotificationRoutingsListParams struct {
	Count  *int
	From   *int
	Target *int
}

func (op *notificationRoutingOp) List(ctx context.Context, projectId string, params NotificationRoutingsListParams) (ret []v1.NotificationRouting, err error) {
	res, err := errorFromDecodedResponse("NotificationRouting.List", func() (*v1.PaginatedNotificationRoutingList, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.ListAlertProjectNotificationRoutings(ctx, v1.ListAlertProjectNotificationRoutingsParams{
			ProjectResourceId: id,
			Count:             into.Opt[v1.OptInt](params.Count),
			From:              into.Opt[v1.OptInt](params.From),
			Target:            into.Opt[v1.OptInt](params.Target),
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
				ret := make([]v1.MatchLabelsItemRequest, len(params.MatchLabels))
				for i, item := range params.MatchLabels {
					ret[i] = v1.MatchLabelsItemRequest{
						Name:  item.Name,
						Value: item.Value,
					}
				}
				return ret
			}(),
			ResendIntervalMinutes: into.Opt[v1.OptInt](params.ResendIntervalMinutes),
		}
		return op.client.CreateAlertProjectNotificationRouting(ctx, &req, v1.CreateAlertProjectNotificationRoutingParams{ProjectResourceId: id})
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
		return op.client.PatchAlertProjectNotificationRouting(ctx, v1.NewOptPatchedNotificationRoutingRequest(v1.PatchedNotificationRoutingRequest{
			NotificationTargetUID: into.Opt[v1.OptUUID](params.NotificationTargetUID),
			MatchLabels: func() []v1.MatchLabelsItemRequest {
				ret := make([]v1.MatchLabelsItemRequest, len(params.MatchLabels))
				for i, item := range params.MatchLabels {
					ret[i] = v1.MatchLabelsItemRequest{
						Name:  item.Name,
						Value: item.Value,
					}
				}
				return ret
			}(),
			ResendIntervalMinutes: into.Opt[v1.OptInt](params.ResendIntervalMinutes),
		}), v1.PatchAlertProjectNotificationRoutingParams{
			ProjectResourceId: id,
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
		return op.client.ReadAlertProjectNotificationRouting(ctx, v1.ReadAlertProjectNotificationRoutingParams{
			ProjectResourceId: id,
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
		return op.client.DeleteAlertProjectNotificationRouting(ctx, v1.DeleteAlertProjectNotificationRoutingParams{
			ProjectResourceId: id,
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
		return op.client.ReorderAlertProjectNotificationRoutings(ctx, orders, v1.ReorderAlertProjectNotificationRoutingsParams{
			ProjectResourceId: id,
		})
	})
}
