// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

package apprun

import (
	"context"
	"errors"
	"net/http"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
)

type TrafficAPI interface {
	// List アプリケーショントラフィック分散を取得
	List(ctx context.Context, appId string) (*v1.HandlerListTraffic, error)
	// Update アプリケーショントラフィック分散を変更
	Update(ctx context.Context, appId string, params *v1.UpdateTrafficBody) (*v1.HandlerUpdateTraffic, error)
}

var _ TrafficAPI = (*trafficOp)(nil)

type trafficOp struct {
	client *v1.Client
}

// NewTrafficOp アプリケーショントラフィック分散関連API
func NewTrafficOp(client *v1.Client) TrafficAPI {
	return &trafficOp{client: client}
}

func (op *trafficOp) List(ctx context.Context, appId string) (*v1.HandlerListTraffic, error) {
	const methodName = "Traffics.List"
	res, err := op.client.ListApplicationTraffic(ctx, v1.ListApplicationTrafficParams{ID: appId})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerListTraffic:
		return result, nil
	case *v1.ListApplicationTrafficBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.ListApplicationTrafficUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ListApplicationTrafficForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ListApplicationTrafficNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.ListApplicationTrafficInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *trafficOp) Update(ctx context.Context, appId string, params *v1.UpdateTrafficBody) (*v1.HandlerUpdateTraffic, error) {
	const methodName = "Traffics.Update"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	res, err := op.client.UpdateApplicationTraffic(ctx, *params, v1.UpdateApplicationTrafficParams{ID: appId})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerUpdateTraffic:
		return result, nil
	case *v1.UpdateApplicationTrafficBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.UpdateApplicationTrafficUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.UpdateApplicationTrafficForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.UpdateApplicationTrafficNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.UpdateApplicationTrafficInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
