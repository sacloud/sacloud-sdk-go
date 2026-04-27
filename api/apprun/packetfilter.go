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

	"github.com/google/uuid"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

type PacketFilterAPI interface {
	// Read パケットフィルタ詳細を取得
	Read(ctx context.Context, appId string) (*v1.HandlerGetPacketFilter, error)
	// Update パケットフィルタを部分的に変更
	Update(ctx context.Context, appId string, params *v1.PatchPacketFilter) (*v1.HandlerPatchPacketFilter, error)
}

var _ PacketFilterAPI = (*packetFilterOp)(nil)

type packetFilterOp struct {
	client *v1.Client
}

// NewPacketFilterOp アプリケーショントラフィック分散関連API
func NewPacketFilterOp(client *v1.Client) PacketFilterAPI {
	return &packetFilterOp{client: client}
}

func (op *packetFilterOp) Read(ctx context.Context, appId string) (*v1.HandlerGetPacketFilter, error) {
	const methodName = "PacketFilter.Read"
	id, err := uuid.Parse(appId)
	if err != nil {
		return nil, NewError(methodName, err)
	}

	res, err := op.client.GetPacketFilter(ctx, v1.GetPacketFilterParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerGetPacketFilter:
		return result, nil
	case *v1.GetPacketFilterBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.GetPacketFilterUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.GetPacketFilterForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.GetPacketFilterNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.GetPacketFilterInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *packetFilterOp) Update(ctx context.Context, appId string, params *v1.PatchPacketFilter) (*v1.HandlerPatchPacketFilter, error) {
	const methodName = "PacketFilter.Update"
	id, err := uuid.Parse(appId)
	if err != nil {
		return nil, NewError(methodName, err)
	}
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	res, err := op.client.PatchPacketFilter(ctx, params, v1.PatchPacketFilterParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerPatchPacketFilter:
		return result, nil
	case *v1.PatchPacketFilterBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.PatchPacketFilterUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.PatchPacketFilterForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.PatchPacketFilterNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.PatchPacketFilterInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
