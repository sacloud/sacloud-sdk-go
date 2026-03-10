// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

package cloudhsm_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func TestCloudHSMClientOp_List(t *testing.T) {
	assert := require.New(t)
	expected := v1.PaginatedCloudHSMClientList{
		Count:   1,
		From:    v1.NewOptInt(0),
		Total:   v1.NewOptInt(1),
		Clients: []v1.CloudHSMClient{TemplateCloudHSMClient},
	}
	client := newTestClient(expected)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()
	clients, err := api.List(ctx)

	assert.NoError(err)
	assert.NotNil(clients)
	assert.Equal(1, len(clients))
}

func TestCloudHSMClientOp_Create(t *testing.T) {
	assert := require.New(t)
	client := newTestClient(TemplateWrappedCreateCloudHSMClient, http.StatusCreated)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	res, err := api.Create(ctx, CloudHSMClientCreateParams{
		Name:        "client-name",
		Certificate: "cert-1",
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateCreateCloudHSMClient.GetName(), res.GetName())
	assert.Equal(TemplateCreateCloudHSMClient.GetCertificate(), res.GetCertificate())
}

func TestCloudHSMClientOp_Create_422(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid request body.")
	client := newTestClient(expected, http.StatusUnprocessableEntity)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	clientObj, err := api.Create(ctx, CloudHSMClientCreateParams{})
	assert.Nil(clientObj)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestCloudHSMClientOp_Read(t *testing.T) {
	assert := require.New(t)
	client := newTestClient(&v1.WrappedCloudHSMClient{Client: TemplateCloudHSMClient})
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	res, err := api.Read(ctx, "client-1")
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateCloudHSMClient.ID, res.ID)
	assert.Equal(TemplateCloudHSMClient.Name, res.Name)
}

func TestCloudHSMClientOp_Read_Error(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("not found")
	client := newTestClient(expected, http.StatusNotFound)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	res, err := api.Read(ctx, "client-1")
	assert.Nil(res)
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestCloudHSMClientOp_Update(t *testing.T) {
	assert := require.New(t)
	updated := TemplateCloudHSMClient
	updated.Name = "updated-name"
	updated.Certificate = "updated-cert"
	client := newTestClient(&v1.WrappedCloudHSMClient{Client: updated})
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	res, err := api.Update(ctx, "client-1", CloudHSMClientUpdateParams{
		Name: "updated-name",
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal("updated-name", res.Name)
	assert.Equal("updated-cert", res.Certificate)
}

func TestCloudHSMClientOp_Update_422(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid request body.")
	client := newTestClient(expected, http.StatusUnprocessableEntity)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	res, err := api.Update(ctx, "client-1", CloudHSMClientUpdateParams{})
	assert.Nil(res)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestCloudHSMClientOp_Delete(t *testing.T) {
	assert := require.New(t)
	client := newTestClient(nil, http.StatusNoContent)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Delete(ctx, "12345")
	assert.NoError(err)
}

func TestCloudHSMClientOp_Delete_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Not found")
	client := newTestClient(expected, http.StatusNotFound)
	api, err := NewClientOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Delete(ctx, "0")
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestCloudHSMClientIntegrated(t *testing.T) {
	assert := require.New(t)
	client := newIntegratedClient(t)

	testutil.PreCheckEnvsFunc("SAKURA_CLOUDHSM_ID")(t)

	ctx := context.Background()
	hsm, err := NewCloudHSMOp(client).Read(ctx, os.Getenv("SAKURA_CLOUDHSM_ID"))
	assert.NoError(err)
	assert.NotNil(hsm)
	assert.Equal(v1.AvailabilityEnumAvailable, hsm.GetAvailability())
	api, err := NewClientOp(client, hsm)
	assert.NoError(err)

	cn := testutil.RandomName("client-integrated-", 16, testutil.CharSetAlphaNum)
	pem, _, err := certpem(cn)
	assert.NoError(err)

	// Create
	created, err := api.Create(ctx, CloudHSMClientCreateParams{
		Name:        cn,
		Certificate: string(pem),
	})
	assert.NoError(err)
	assert.NotNil(created)

	// Delete
	t.Cleanup(func() {
		err := api.Delete(ctx, created.GetID())
		assert.NoError(err)
	})

	// List
	clients, err := api.List(ctx)
	assert.NoError(err)
	assert.NotNil(clients)
	assert.NotEmpty(clients)
}

func certpem(cn string) (certPEM, keyPEM []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, err
	}

	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"Test Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return certPEM, keyPEM, nil
}
