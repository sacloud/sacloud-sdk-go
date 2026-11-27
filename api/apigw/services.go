// Copyright 2025- The sacloud/apigw-api-go authors
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

package apigw

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
)

type ServiceAPI interface {
	List(ctx context.Context) ([]v1.ServiceDetailResponse, error)
	Create(ctx context.Context, request *v1.ServiceDetailRequest) (*v1.ServiceDetailRequest, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.ServiceDetailResponse, error)
	Update(ctx context.Context, request *v1.ServiceDetail, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ ServiceAPI = (*serviceOp)(nil)

type serviceOp struct {
	client *v1.Client
}

func NewServiceOp(client *v1.Client) ServiceAPI {
	return &serviceOp{client: client}
}

func (op *serviceOp) List(ctx context.Context) ([]v1.ServiceDetailResponse, error) {
	res, err := errorFromDecodedResponse[v1.GetServicesOK]("Service.List", func() (any, error) { return op.client.GetServices(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetServices(), nil
}

func (op *serviceOp) Create(ctx context.Context, request *v1.ServiceDetailRequest) (*v1.ServiceDetailRequest, error) {
	res, err := errorFromDecodedResponse[v1.AddServiceCreated]("Service.Create", func() (any, error) { return op.client.AddService(ctx, request) })
	if err != nil {
		return nil, err
	}

	s, ok := new(res.GetApigw()).GetService().Get()
	if !ok {
		return nil, NewError("Service.Create", errors.New("response does not contain a service"))
	}

	return &s, nil
}

func (op *serviceOp) Read(ctx context.Context, id uuid.UUID) (*v1.ServiceDetailResponse, error) {
	res, err := errorFromDecodedResponse[v1.GetServiceByIdOK]("Service.Read", func() (any, error) { return op.client.GetServiceById(ctx, v1.GetServiceByIdParams{ServiceId: id}) })
	if err != nil {
		return nil, err
	}

	s, ok := new(res.GetApigw()).GetService().Get()
	if !ok {
		return nil, NewError("Service.Read", errors.New("response does not contain a service"))
	}

	return &s, nil
}

func (op *serviceOp) Update(ctx context.Context, request *v1.ServiceDetail, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateServiceNoContent]("Service.Update", func() (any, error) {
		return op.client.UpdateService(ctx, request, v1.UpdateServiceParams{ServiceId: id})
	})
}

func (op *serviceOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteServiceNoContent]("Service.Delete", func() (any, error) {
		return op.client.DeleteService(ctx, v1.DeleteServiceParams{ServiceId: id})
	})
}
