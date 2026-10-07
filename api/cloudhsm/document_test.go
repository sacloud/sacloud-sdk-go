// Copyright 2026- The sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package cloudhsm_test

import (
	"context"
	"net/url"
	"testing"

	. "github.com/sacloud/sacloud-sdk-go/api/cloudhsm"
	v1 "github.com/sacloud/sacloud-sdk-go/api/cloudhsm/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestDocumentOp_List(t *testing.T) {
	assert := require.New(t)
	expected := v1.PaginatedCloudHSMDocumentList{
		Count:             1,
		From:              10,
		Total:             11,
		CloudHSMDocuments: []v1.CloudHSMDocument{TemplateDocument},
		IsOk:              true,
	}
	client := newTestClient(expected)
	api := NewDocumentOp(client, "license-1")
	ctx := context.Background()

	documents, err := api.List(ctx, nil, nil)

	assert.NoError(err)
	assert.Len(documents, 1)
	assert.Equal(TemplateDocument, documents[0])
}

func TestDocumentOp_Download(t *testing.T) {
	assert := require.New(t)
	expected := map[string]any{
		"Document": map[string]string{
			"URL": "https://example.com/document.zip",
		},
		"is_ok": true,
	}
	client := newTestClient(expected)
	api := NewDocumentOp(client, "license-1")

	document, err := api.Download(context.Background(), "document-1")

	assert.NoError(err)
	assert.Equal(url.URL{Scheme: "https", Host: "example.com", Path: "/document.zip"}, document.GetURL())
}

func TestDocumentOp_List_404(t *testing.T) {
	assert := require.New(t)
	client := newTestClient(newErrorResponse("No documents found."), 404)
	api := NewDocumentOp(client, "missing-license")

	documents, err := api.List(context.Background(), nil, nil)

	assert.Nil(documents)
	assert.ErrorContains(err, "not found")
}

func TestDocumentOp_Download_422(t *testing.T) {
	assert := require.New(t)
	client := newTestClient(newErrorResponse("Invalid document."), 422)
	api := NewDocumentOp(client, "license-1")

	document, err := api.Download(context.Background(), "invalid-document")

	assert.Nil(document)
	assert.ErrorContains(err, "invalid parameter")
}
