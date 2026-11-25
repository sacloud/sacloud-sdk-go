// Copyright 2026- The sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package cloudhsm

import (
	"context"

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
	resp, err := errorFromDecodedResponse("Document.List", func() (*v1.PaginatedCloudHSMDocumentList, error) {
		return op.client.ListCloudHSMDocuments(ctx, request)
	})
	if err != nil {
		return nil, err
	}

	return resp.GetCloudHSMDocuments(), nil
}

func (op *DocumentOp) Download(ctx context.Context, id string) (*v1.CloudHSMDocumentDownload, error) {
	resp, err := errorFromDecodedResponse("Document.Download", func() (*v1.WrappedCloudHSMDocumentDownload, error) {
		return op.client.DownloadCloudHSMDocument(ctx, v1.DownloadCloudHSMDocumentParams{
			ID:                id,
			LicenseResourceID: op.licenseResourceID,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetDocument()), nil
}
