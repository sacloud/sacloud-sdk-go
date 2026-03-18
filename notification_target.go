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
	"net/url"
	"strconv"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
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

func (op *notificationTargetOp) List(ctx context.Context, projectId string, p NotificationTargetsListParams) (ret []v1.NotificationTarget, err error) {
	res, err := errorFromDecodedResponse("NotificationTarget.List", func() (*v1.PaginatedNotificationTargetList, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationTargetsList(ctx, v1.AlertsProjectsNotificationTargetsListParams{
			ProjectResourceID: id,
			Count:             intoOpt[v1.OptInt](p.Count),
			From:              intoOpt[v1.OptInt](p.From),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *notificationTargetOp) Read(ctx context.Context, projectId string, uid uuid.UUID) (*v1.NotificationTarget, error) {
	return errorFromDecodedResponse("NotificationTarget.Read", func() (*v1.NotificationTarget, error) {
		pid, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationTargetsRetrieve(ctx, v1.AlertsProjectsNotificationTargetsRetrieveParams{
			ProjectResourceID: pid,
			UID:               uid,
		})
	})
}

type NotificationTargetCreateParams struct {
	ServiceType v1.NotificationTargetServiceType
	URL         *url.URL
	Description *string
}

func (cp *NotificationTargetCreateParams) urlstr() (ret *string) {
	if cp.URL != nil {
		ret = saclient.Ptr(cp.URL.String())
	}
	return
}

func (op *notificationTargetOp) Create(ctx context.Context, projectId string, params NotificationTargetCreateParams) (*v1.NotificationTarget, error) {
	return errorFromDecodedResponse("NotificationTarget.Create", func() (*v1.NotificationTarget, error) {
		pid, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationTargetsCreate(ctx, &v1.NotificationTargetRequest{
			ServiceType: v1.NotificationTargetRequestServiceType(params.ServiceType),
			URL:         intoOpt[v1.OptString](params.urlstr()),
			Description: intoOpt[v1.OptString](params.Description),
		}, v1.AlertsProjectsNotificationTargetsCreateParams{ProjectResourceID: pid})
	})
}

type NotificationTargetUpdateParams struct {
	ServiceType *v1.PatchedNotificationTargetRequestServiceType
	URL         *string
	Description *string
}

func (op *notificationTargetOp) Update(ctx context.Context, projectId string, uid uuid.UUID, params NotificationTargetUpdateParams) (*v1.NotificationTarget, error) {
	return errorFromDecodedResponse("NotificationTarget.Update", func() (*v1.NotificationTarget, error) {
		pid, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsNotificationTargetsPartialUpdate(ctx, v1.NewOptPatchedNotificationTargetRequest(v1.PatchedNotificationTargetRequest{
			ServiceType: intoOpt[v1.OptPatchedNotificationTargetRequestServiceType](params.ServiceType),
			URL:         intoOpt[v1.OptString](params.URL),
			Description: intoOpt[v1.OptString](params.Description),
		}), v1.AlertsProjectsNotificationTargetsPartialUpdateParams{
			ProjectResourceID: pid,
			UID:               uid,
		})
	})
}

func (op *notificationTargetOp) Delete(ctx context.Context, projectId string, uid uuid.UUID) error {
	return errorFromDecodedResponse1("NotificationTarget.Delete", func() error {
		pid, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.AlertsProjectsNotificationTargetsDestroy(ctx, v1.AlertsProjectsNotificationTargetsDestroyParams{
			ProjectResourceID: pid,
			UID:               uid,
		})
	})
}
