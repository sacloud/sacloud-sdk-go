// Copyright 2026- The sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package cloudhsm

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/cloudhsm/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type DocumentAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.CloudHSMDocument, error)
	Download(ctx context.Context, id string) (*v1.CloudHSMDocumentDownload, error)
}

var _ DocumentAPI = (*DocumentOp)(nil)

type DocumentOp struct {
	client            *v1.Client
	licenseResourceID string
}

func NewDocumentOp(client *v1.Client, licenseResourceID string) DocumentAPI {
	return &DocumentOp{
		client:            client,
		licenseResourceID: licenseResourceID,
	}
}

func (op *DocumentOp) List(ctx context.Context, count, from *int) ([]v1.CloudHSMDocument, error) {
	request := v1.ListCloudHSMDocumentsParams{
		LicenseResourceID: op.licenseResourceID,
		Count:             into.Opt[v1.OptInt](count),
		From:              into.Opt[v1.OptInt](from),
	}

	resp, err := op.client.ListCloudHSMDocuments(ctx, request)
	if err == nil {
		return resp.GetCloudHSMDocuments(), nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("Document.List", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("Document.List", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return nil, NewAPIError("Document.List", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *DocumentOp) Download(ctx context.Context, id string) (*v1.CloudHSMDocumentDownload, error) {
	resp, err := op.client.DownloadCloudHSMDocument(
		ctx,
		v1.DownloadCloudHSMDocumentParams{
			ID:                id,
			LicenseResourceID: op.licenseResourceID,
		},
	)
	//nolint:gocritic
	if err == nil {
		document := resp.GetDocument()
		return &document, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("Document.Download", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("Document.Download", e.StatusCode, errors.Wrap(err, "not found"))
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("Document.Download", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("Document.Download", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}
