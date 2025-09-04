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

type AlertRuleAPI interface {
	List(ctx context.Context, params v1.AlertsProjectsRulesListParams) ([]v1.AlertRule, error)
	Create(ctx context.Context, projectId string, params *v1.AlertRule) (*v1.AlertRule, error)
	Read(ctx context.Context, projectId string, ruleId string) (*v1.AlertRule, error)
	Update(ctx context.Context, projectId string, ruleId string, params *v1.AlertRule) (*v1.AlertRule, error)
	Delete(ctx context.Context, projectId string, ruleId string) error

	ListHistories(ctx context.Context, params v1.AlertsProjectsRulesHistoriesListParams) ([]v1.History, error)
	ReadHistory(ctx context.Context, projectId string, ruleId string, historyId string) (*v1.History, error)
}

var _ AlertRuleAPI = (*alertRuleOp)(nil)

type alertRuleOp struct {
	client *v1.Client
}

func NewAlertRuleOp(client *v1.Client) AlertRuleAPI {
	return &alertRuleOp{client: client}
}

func (op *alertRuleOp) List(ctx context.Context, params v1.AlertsProjectsRulesListParams) ([]v1.AlertRule, error) {
	result, err := op.client.AlertsProjectsRulesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "project not found"))
		default:
			return nil, NewAPIError("AlertRule.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *alertRuleOp) Create(ctx context.Context, projectId string, params *v1.AlertRule) (*v1.AlertRule, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Create", 0, err)
	}
	query := v1.AlertsProjectsRulesCreateParams{
		ProjectPk: int(intProjectId),
	}
	result, err := op.client.AlertsProjectsRulesCreate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "project not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new alert rule"))
		default:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *alertRuleOp) Read(ctx context.Context, projectId string, ruleId string) (*v1.AlertRule, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Read", 0, err)
	}
	intRuleId, err := strconv.ParseInt(ruleId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Read", 0, err)
	}
	query := v1.AlertsProjectsRulesRetrieveParams{
		ProjectPk: int(intProjectId),
		ID:        int(intRuleId),
	}
	result, err := op.client.AlertsProjectsRulesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		default:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Read", 0, err)
	}
	return result, nil
}

func (op *alertRuleOp) Update(ctx context.Context, projectId string, ruleId string, params *v1.AlertRule) (*v1.AlertRule, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, err)
	}
	intRuleId, err := strconv.ParseInt(ruleId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, err)
	}
	query := v1.AlertsProjectsRulesUpdateParams{
		ProjectPk: int(intProjectId),
		ID:        int(intRuleId),
	}
	result, err := op.client.AlertsProjectsRulesUpdate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, err)
	}
	return result, nil
}

func (op *alertRuleOp) Delete(ctx context.Context, projectId string, ruleId string) error {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewAPIError("AlertRule.Delete", 0, err)
	}
	intRuleId, err := strconv.ParseInt(ruleId, 10, 64)
	if err != nil {
		return NewAPIError("AlertRule.Delete", 0, err)
	}
	query := v1.AlertsProjectsRulesDestroyParams{
		ProjectPk: int(intProjectId),
		ID:        int(intRuleId),
	}
	err = op.client.AlertsProjectsRulesDestroy(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		case http.StatusBadRequest:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("AlertRule.Delete", 0, err)
	}
	return nil
}

func (op *alertRuleOp) ListHistories(ctx context.Context, params v1.AlertsProjectsRulesHistoriesListParams) ([]v1.History, error) {
	result, err := op.client.AlertsProjectsRulesHistoriesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "project or rule not found"))
		default:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.ListHistories", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *alertRuleOp) ReadHistory(ctx context.Context, projectId string, ruleId string, historyId string) (*v1.History, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	intRuleId, err := strconv.ParseInt(ruleId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	intHistoryId, err := strconv.ParseInt(historyId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	query := v1.AlertsProjectsRulesHistoriesRetrieveParams{
		ProjectPk: int(intProjectId),
		RulePk:    int(intRuleId),
		ID:        int(intHistoryId),
	}
	result, err := op.client.AlertsProjectsRulesHistoriesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "history not found"))
		default:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	return result, nil
}
