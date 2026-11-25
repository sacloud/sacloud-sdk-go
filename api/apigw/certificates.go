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

type CertificateAPI interface {
	List(ctx context.Context) ([]v1.Certificate, error)
	Create(ctx context.Context, request *v1.Certificate) (*v1.Certificate, error)
	Update(ctx context.Context, request *v1.Certificate, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ CertificateAPI = (*certificateOp)(nil)

type certificateOp struct {
	client *v1.Client
}

func NewCertificateOp(client *v1.Client) CertificateAPI {
	return &certificateOp{client: client}
}

func (op *certificateOp) List(ctx context.Context) ([]v1.Certificate, error) {
	res, err := errorFromDecodedResponse[v1.GetCertificatesOK]("Certificate.List", func() (any, error) { return op.client.GetCertificates(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetCertificates(), nil
}

func (op *certificateOp) Create(ctx context.Context, request *v1.Certificate) (*v1.Certificate, error) {
	res, err := errorFromDecodedResponse[v1.AddCertificateCreated]("Certificate.Create", func() (any, error) { return op.client.AddCertificate(ctx, request) })
	if err != nil {
		return nil, err
	}

	c, ok := new(res.GetApigw()).GetCertificate().Get()
	if !ok {
		return nil, NewError("Certificate.Create", errors.New("response does not contain a certificate"))
	}

	return &c, nil
}

func (op *certificateOp) Update(ctx context.Context, request *v1.Certificate, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateCertificateNoContent]("Certificate.Update", func() (any, error) {
		return op.client.UpdateCertificate(ctx, request, v1.UpdateCertificateParams{CertificateId: id})
	})
}

func (op *certificateOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteCertificateNoContent]("Certificate.Delete", func() (any, error) {
		return op.client.DeleteCertificate(ctx, v1.DeleteCertificateParams{CertificateId: id})
	})
}
