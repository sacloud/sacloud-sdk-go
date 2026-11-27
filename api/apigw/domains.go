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
	res, err := errorFromDecodedResponse[v1.GetDomainsOK]("Domain.List", func() (any, error) { return op.client.GetDomains(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetDomains(), nil
}

func (op *domainOp) Create(ctx context.Context, request *v1.Domain) (*v1.Domain, error) {
	res, err := errorFromDecodedResponse[v1.AddDomainCreated]("Domain.Create", func() (any, error) { return op.client.AddDomain(ctx, request) })
	if err != nil {
		return nil, err
	}

	d, ok := new(res.GetApigw()).GetDomain().Get()
	if !ok {
		return nil, NewError("Domain.Create", errors.New("response does not contain a domain"))
	}

	return &d, nil
}

func (op *domainOp) Update(ctx context.Context, request *v1.DomainPUT, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateDomainNoContent]("Domain.Update", func() (any, error) {
		return op.client.UpdateDomain(ctx, request, v1.UpdateDomainParams{DomainId: id})
	})
}

func (op *domainOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteDomainNoContent]("Domain.Delete", func() (any, error) {
		return op.client.DeleteDomain(ctx, v1.DeleteDomainParams{DomainId: id})
	})
}
