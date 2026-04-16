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
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
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
	res, err := op.client.GetServices(ctx)
	if err != nil {
		return nil, NewAPIError("Service.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetServicesOK:
		return p.Apigw.Services, nil
	case *v1.GetServicesBadRequest:
		return nil, NewAPIError("Service.List", 400, errors.New(p.Message.Value))
	case *v1.GetServicesUnauthorized:
		return nil, NewAPIError("Service.List", 401, errors.New(p.Message.Value))
	case *v1.GetServicesInternalServerError:
		return nil, NewAPIError("Service.List", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Service.List", 0, nil)
}

func (op *serviceOp) Create(ctx context.Context, request *v1.ServiceDetailRequest) (*v1.ServiceDetailRequest, error) {
	res, err := op.client.AddService(ctx, request)
	if err != nil {
		return nil, NewAPIError("Service.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.AddServiceCreated:
		return &p.Apigw.Service.Value, nil
	case *v1.AddServiceBadRequest:
		return nil, NewAPIError("Service.Create", 400, errors.New(p.Message.Value))
	case *v1.AddServiceNotFound:
		return nil, NewAPIError("Service.Create", 404, errors.New(p.Message.Value))
	case *v1.AddServiceConflict:
		return nil, NewAPIError("Service.Create", 409, errors.New(p.Message.Value))
	case *v1.AddServiceInternalServerError:
		return nil, NewAPIError("Service.Create", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Service.Create", 0, nil)
}

func (op *serviceOp) Read(ctx context.Context, id uuid.UUID) (*v1.ServiceDetailResponse, error) {
	res, err := op.client.GetServiceById(ctx, v1.GetServiceByIdParams{ServiceId: id})
	if err != nil {
		return nil, NewAPIError("Service.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetServiceByIdOK:
		return &p.Apigw.Service.Value, nil
	case *v1.GetServiceByIdBadRequest:
		return nil, NewAPIError("Service.Read", 400, errors.New(p.Message.Value))
	case *v1.GetServiceByIdNotFound:
		return nil, NewAPIError("Service.Read", 404, errors.New(p.Message.Value))
	case *v1.GetServiceByIdInternalServerError:
		return nil, NewAPIError("Service.Read", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Service.Read", 0, nil)
}

func (op *serviceOp) Update(ctx context.Context, request *v1.ServiceDetail, id uuid.UUID) error {
	res, err := op.client.UpdateService(ctx, request, v1.UpdateServiceParams{ServiceId: id})
	if err != nil {
		return NewAPIError("Service.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateServiceNoContent:
		return nil
	case *v1.UpdateServiceBadRequest:
		return NewAPIError("Service.Update", 400, errors.New(p.Message.Value))
	case *v1.UpdateServiceNotFound:
		return NewAPIError("Service.Update", 404, errors.New(p.Message.Value))
	case *v1.UpdateServiceInternalServerError:
		return NewAPIError("Service.Update", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Service.Update", 0, nil)
}

func (op *serviceOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteService(ctx, v1.DeleteServiceParams{ServiceId: id})
	if err != nil {
		return NewAPIError("Service.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteServiceNoContent:
		return nil
	case *v1.DeleteServiceBadRequest:
		return NewAPIError("Service.Delete", 400, errors.New(p.Message.Value))
	case *v1.DeleteServiceUnauthorized:
		return NewAPIError("Service.Delete", 401, errors.New(p.Message.Value))
	case *v1.DeleteServiceNotFound:
		return NewAPIError("Service.Delete", 404, errors.New(p.Message.Value))
	case *v1.DeleteServiceInternalServerError:
		return NewAPIError("Service.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Service.Delete", 0, nil)
}
