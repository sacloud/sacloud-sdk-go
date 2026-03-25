// Copyright 2026- The sacloud/service-endpoint-gateway-api-go authors
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

package service_endpoint_gateway

import (
	"context"
	"errors"

	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

type ConstructAPI interface {
	List(ctx context.Context) (*v1.ModelsApplianceApplianceListResponseBody, error)
	Create(ctx context.Context, request v1.ModelsApplianceApplianceCreateRequest) (*v1.ModelsApplianceApplianceCreateResponseBody, error)
	Read(ctx context.Context, id string) (*v1.ModelsApplianceApplianceGetResponseBody, error)
	Update(ctx context.Context, id string, request v1.ModelsApplianceApplianceUpdateRequest) (*v1.ModelsApplianceApplianceGetResponseBody, error)
	Apply(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	ReadInterface(ctx context.Context, applianceID string, interfaceID string) (*v1.ModelsApplianceApplianceGetInterfaceResponseBody, error)
}

var _ ConstructAPI = (*ConstructOp)(nil)

type ConstructOp struct {
	client *v1.Client
}

func NewConstructOp(client *v1.Client) ConstructAPI {
	return &ConstructOp{client: client}
}

func (op *ConstructOp) List(ctx context.Context) (*v1.ModelsApplianceApplianceListResponseBody, error) {
	const methodName = "Construct.List"
	res, err := op.client.SegList(ctx)
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ConstructOp) Create(ctx context.Context, request v1.ModelsApplianceApplianceCreateRequest) (*v1.ModelsApplianceApplianceCreateResponseBody, error) {
	const methodName = "Construct.Create"

	request.Appliance.Class = v1.ModelsApplianceApplianceCreateBodyClassServiceendpointgateway
	request.Appliance.Plan.ID = v1.ModelsPlanPlanID1

	res, err := op.client.SegCreate(ctx, &request)
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ConstructOp) Read(ctx context.Context, id string) (*v1.ModelsApplianceApplianceGetResponseBody, error) {
	const methodName = "Construct.Read"

	res, err := op.client.SegGet(ctx, v1.SegGetParams{ApplianceID: id})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ConstructOp) Update(ctx context.Context, id string, request v1.ModelsApplianceApplianceUpdateRequest) (*v1.ModelsApplianceApplianceGetResponseBody, error) {
	const methodName = "Construct.Update"

	res, err := op.client.SegUpdate(ctx, &request, v1.SegUpdateParams{ApplianceID: id})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ConstructOp) Apply(ctx context.Context, id string) error {
	const methodName = "Construct.Apply"

	_, err := op.client.SegApply(ctx, v1.OptModelsApplianceApplianceUpdateRequest{Set: false}, v1.SegApplyParams{ApplianceID: id})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, err)
		}
		return err
	}
	return nil
}

func (op *ConstructOp) Delete(ctx context.Context, id string) error {
	const methodName = "Construct.Delete"

	_, err := op.client.SegDelete(ctx, v1.SegDeleteParams{ApplianceID: id})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, err)
		}
		return err
	}
	return nil
}

func (op *ConstructOp) ReadInterface(ctx context.Context, applianceID string, interfaceID string) (*v1.ModelsApplianceApplianceGetInterfaceResponseBody, error) {
	const methodName = "Construct.ReadInterface"

	res, err := op.client.SegInterfaceGetInterface(ctx, v1.SegInterfaceGetInterfaceParams{ApplianceID: applianceID, InterfaceID: interfaceID})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}
