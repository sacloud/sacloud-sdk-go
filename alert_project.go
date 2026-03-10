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
	"strconv"
	"time"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type AlertProjectAPI interface {
	List(ctx context.Context, count *int, from *int) ([]v1.AlertProject, error)
	Create(ctx context.Context, params AlertProjectCreateParams) (*v1.AlertProject, error)
	Read(ctx context.Context, id string) (*v1.AlertProject, error)
	Update(ctx context.Context, id string, request AlertProjectUpdateParams) (*v1.AlertProject, error)
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

func (op *alertProjectOp) List(ctx context.Context, count *int, from *int) (ret []v1.AlertProject, err error) {
	res, err := errorFromDecodedResponse("AlertProject.List", func() (*v1.PaginatedAlertProjectList, error) {
		return op.client.AlertsProjectsList(ctx, v1.AlertsProjectsListParams{
			Count: intoOpt[v1.OptInt](count),
			From:  intoOpt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *alertProjectOp) Read(ctx context.Context, id string) (*v1.AlertProject, error) {
	res, err := errorFromDecodedResponse("AlertProject.Read", func() (*v1.WrappedAlertProject, error) {
		if intId, err := strconv.ParseInt(id, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsRetrieve(ctx, v1.AlertsProjectsRetrieveParams{ResourceID: intId})
		}
	})

	return unwrapE[*v1.AlertProject](res, err)
}

type AlertProjectCreateParams struct {
	Name        string
	Description *string
}

func (op *alertProjectOp) Create(ctx context.Context, params AlertProjectCreateParams) (*v1.AlertProject, error) {
	return errorFromDecodedResponse("AlertProject.Create", func() (*v1.AlertProject, error) {
		return op.client.AlertsProjectsCreate(ctx, &v1.AlertProjectCreate{
			Name:        params.Name,
			Description: intoOpt[v1.OptString](params.Description),
		})
	})
}

type AlertProjectUpdateParams struct {
	Name        *string
	Description *string
}

func (op *alertProjectOp) Update(ctx context.Context, id string, params AlertProjectUpdateParams) (*v1.AlertProject, error) {
	res, err := errorFromDecodedResponse("AlertProject.Update", func() (*v1.WrappedAlertProject, error) {
		if intId, err := strconv.ParseInt(id, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsPartialUpdate(
				ctx,
				v1.NewOptPatchedAlertProject(v1.PatchedAlertProject{
					Name:        intoOpt[v1.OptString](params.Name),
					Description: intoOpt[v1.OptString](params.Description),
				}),
				v1.AlertsProjectsPartialUpdateParams{ResourceID: intId},
			)
		}
	})

	return unwrapE[*v1.AlertProject](res, err)
}

func (op *alertProjectOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponse1("AlertProject.Delete", func() error {
		if intId, err := strconv.ParseInt(id, 10, 64); err != nil {
			return err
		} else {
			return op.client.AlertsProjectsDestroy(ctx, v1.AlertsProjectsDestroyParams{ResourceID: intId})
		}
	})
}

type AlertsProjectsHistoriesListParams struct {
	ProjectID string // mandatory
	Count     *int
	From      *int
	Open      *bool
	Severity  *v1.AlertsProjectsHistoriesListSeverity
	StartsAt  *time.Time
}

func (op *alertProjectOp) ListHistories(ctx context.Context, params AlertsProjectsHistoriesListParams) (ret []v1.History, err error) {
	res, err := errorFromDecodedResponse("AlertProject.ListHistories", func() (*v1.PaginatedHistoryList, error) {
		if intProjectId, err := strconv.ParseInt(params.ProjectID, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsHistoriesList(ctx, v1.AlertsProjectsHistoriesListParams{
				ProjectResourceID: intProjectId,
				Count:             intoOpt[v1.OptInt](params.Count),
				From:              intoOpt[v1.OptInt](params.From),
				Open:              intoOpt[v1.OptBool](params.Open),
				Severity:          intoOpt[v1.OptAlertsProjectsHistoriesListSeverity](params.Severity),
				StartsAt:          intoOpt[v1.OptDateTime](params.StartsAt),
			})
		}
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *alertProjectOp) ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error) {
	return errorFromDecodedResponse("AlertProject.ReadHistory", func() (*v1.History, error) {
		if intProjectId, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsHistoriesRetrieve(ctx, v1.AlertsProjectsHistoriesRetrieveParams{
				UID:               historyId,
				ProjectResourceID: intProjectId,
			})
		}
	})
}
