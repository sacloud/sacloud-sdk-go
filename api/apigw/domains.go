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

type DomainAPI interface {
	List(ctx context.Context) ([]v1.Domain, error)
	Create(ctx context.Context, request *v1.Domain) (*v1.Domain, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, request *v1.DomainPUT, id uuid.UUID) error
}

var _ DomainAPI = (*domainOp)(nil)

type domainOp struct {
	client *v1.Client
}

func NewDomainOp(client *v1.Client) DomainAPI {
	return &domainOp{client: client}
}

func (op *domainOp) List(ctx context.Context) ([]v1.Domain, error) {
	res, err := op.client.GetDomains(ctx)
	if err != nil {
		return nil, NewAPIError("Domain.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetDomainsOK:
		return p.Apigw.Domains, nil
	case *v1.GetDomainsBadRequest:
		return nil, NewAPIError("Domain.List", 400, errors.New(p.Message.Value))
	case *v1.GetDomainsUnauthorized:
		return nil, NewAPIError("Domain.List", 401, errors.New(p.Message.Value))
	case *v1.GetDomainsInternalServerError:
		return nil, NewAPIError("Domain.List", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Domain.List", 0, nil)
}

func (op *domainOp) Create(ctx context.Context, request *v1.Domain) (*v1.Domain, error) {
	res, err := op.client.AddDomain(ctx, request)
	if err != nil {
		return nil, NewAPIError("Domain.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.AddDomainCreated:
		return &p.Apigw.Domain.Value, nil
	case *v1.AddDomainBadRequest:
		return nil, NewAPIError("Domain.Create", 400, errors.New(p.Message.Value))
	case *v1.AddDomainUnauthorized:
		return nil, NewAPIError("Domain.Create", 401, errors.New(p.Message.Value))
	case *v1.AddDomainConflict:
		return nil, NewAPIError("Domain.Create", 409, errors.New(p.Message.Value))
	case *v1.AddDomainInternalServerError:
		return nil, NewAPIError("Domain.Create", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Domain.Create", 0, nil)
}

func (op *domainOp) Update(ctx context.Context, request *v1.DomainPUT, id uuid.UUID) error {
	res, err := op.client.UpdateDomain(ctx, request, v1.UpdateDomainParams{DomainId: id})
	if err != nil {
		return NewAPIError("Domain.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateDomainNoContent:
		return nil
	case *v1.UpdateDomainBadRequest:
		return NewAPIError("Domain.Update", 400, errors.New(p.Message.Value))
	case *v1.UpdateDomainUnauthorized:
		return NewAPIError("Domain.Update", 401, errors.New(p.Message.Value))
	case *v1.UpdateDomainNotFound:
		return NewAPIError("Domain.Update", 404, errors.New(p.Message.Value))
	case *v1.UpdateDomainConflict:
		return NewAPIError("Domain.Update", 409, errors.New(p.Message.Value))
	case *v1.UpdateDomainInternalServerError:
		return NewAPIError("Domain.Update", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Domain.Update", 0, nil)
}

func (op *domainOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteDomain(ctx, v1.DeleteDomainParams{DomainId: id})
	if err != nil {
		return NewAPIError("Domain.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteDomainNoContent:
		return nil
	case *v1.DeleteDomainBadRequest:
		return NewAPIError("Domain.Delete", 400, errors.New(p.Message.Value))
	case *v1.DeleteDomainUnauthorized:
		return NewAPIError("Domain.Delete", 401, errors.New(p.Message.Value))
	case *v1.DeleteDomainNotFound:
		return NewAPIError("Domain.Delete", 404, errors.New(p.Message.Value))
	case *v1.DeleteDomainInternalServerError:
		return NewAPIError("Domain.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Domain.Delete", 0, nil)
}
