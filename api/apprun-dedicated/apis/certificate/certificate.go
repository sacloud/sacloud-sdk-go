// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package certificate

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
)

type CertificateAPI interface {
	// List returns the list of Certificates, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *v1.CertificateID) (list []v1.ReadCertificate, nextCursor *v1.CertificateID, err error)
	Create(ctx context.Context, params CreateParams) (cert *v1.CreatedCertificate, err error)
	Read(ctx context.Context, id v1.CertificateID) (cert *v1.ReadCertificate, err error)
	Update(ctx context.Context, id v1.CertificateID, params UpdateParams) error
	Delete(ctx context.Context, id v1.CertificateID) error
}

type CertificateOp struct {
	client    *v1.Client
	clusterID v1.ClusterID
}

func NewCertificateOp(client *v1.Client, clusterID v1.ClusterID) *CertificateOp {
	return &CertificateOp{
		client:    client,
		clusterID: clusterID,
	}
}

func (op *CertificateOp) List(ctx context.Context, maxItems int64, cursor *v1.CertificateID) (list []v1.ReadCertificate, nextCursor *v1.CertificateID, err error) {
	res, err := common.ErrorFromDecodedResponse("Certificate.List", func() (*v1.ListCertificateResponse, error) {
		return op.client.ListCertificate(ctx, v1.ListCertificateParams{
			ClusterID: op.clusterID,
			Cursor:    common.IntoOpt[v1.OptCertificateID](cursor),
			MaxItems:  maxItems,
		})
	})

	if res != nil {
		list = res.Certificates
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

type CreateParams struct {
	Name                       string
	CertificatePEM             string
	PrivateKeyPEM              string
	IntermediateCertificatePEM *string
}

func (op *CertificateOp) Create(ctx context.Context, req CreateParams) (ret *v1.CreatedCertificate, err error) {
	res, err := common.ErrorFromDecodedResponse("Certificate.Create", func() (*v1.CreateCertificateResponse, error) {
		request := v1.CreateCertificate{
			Name:                       req.Name,
			CertificatePem:             req.CertificatePEM,
			PrivatekeyPem:              req.PrivateKeyPEM,
			IntermediateCertificatePem: common.IntoOpt[v1.OptString](req.IntermediateCertificatePEM),
		}

		return op.client.CreateCertificate(ctx, &request, v1.CreateCertificateParams{ClusterID: op.clusterID})
	})

	if res != nil {
		ret = &res.Certificate
	}

	return
}

func (op *CertificateOp) Read(ctx context.Context, id v1.CertificateID) (cert *v1.ReadCertificate, err error) {
	res, err := common.ErrorFromDecodedResponse("Certificate.Read", func() (*v1.GetCertificateResponse, error) {
		return op.client.GetCertificate(ctx, v1.GetCertificateParams{
			ClusterID:     op.clusterID,
			CertificateID: id,
		})
	})

	if res != nil {
		cert = &res.Certificate
	}
	return
}

type UpdateParams CreateParams

func (op *CertificateOp) Update(ctx context.Context, id v1.CertificateID, request UpdateParams) error {
	return common.ErrorFromDecodedResponseE("Certificate.Update", func() error {
		request := v1.UpdateCertificate{
			Name:                       request.Name,
			CertificatePem:             request.CertificatePEM,
			PrivatekeyPem:              request.PrivateKeyPEM,
			IntermediateCertificatePem: common.IntoOpt[v1.OptString](request.IntermediateCertificatePEM),
		}
		p := v1.UpdateCertificateParams{
			ClusterID:     op.clusterID,
			CertificateID: id,
		}

		return op.client.UpdateCertificate(ctx, &request, p)
	})
}

func (op *CertificateOp) Delete(ctx context.Context, id v1.CertificateID) error {
	return common.ErrorFromDecodedResponseE("Certificate.Delete", func() error {
		return op.client.DeleteCertificate(ctx, v1.DeleteCertificateParams{
			ClusterID:     op.clusterID,
			CertificateID: id,
		})
	})
}

var _ CertificateAPI = (*CertificateOp)(nil)
