// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package certificate_test

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/certificate"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api CertificateAPI) {
	clid := v1.ClusterID(uuid.New())
	assert = require.New(t)
	client, err := apprun_test.NewTestClient(v, s...)
	assert.NoError(err)
	api = NewCertificateOp(client, clid)
	return assert, api
}

func TestNewCertificateOp(t *testing.T) {
	assert, api := setup(t, nil, http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.ListCertificateResponse
	expected.SetFake()
	expected.SetCertificates(make([]v1.ReadCertificate, 2))
	expected.Certificates[0] = apprun_test.FakeCertificate()
	expected.Certificates[1] = apprun_test.FakeCertificate()
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 0, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Equal(expected.GetCertificates(), actual)
	assert.Equal(expected.GetNextCursor().Or(v1.CertificateID(uuid.Nil)), *cursor)
}

func TestList_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, cursor, err := api.List(t.Context(), 0, nil)

	assert.Error(err)
	assert.Nil(actual)
	assert.Nil(cursor)
	assert.False(saclient.IsNotFoundError(err))
}

func TestCreate(t *testing.T) {
	var expected v1.CreateCertificateResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	c, p, err := apprun_test.OreSign()
	assert.NoError(err)

	actual, err := api.Create(t.Context(), CreateParams{
		Name:           testutil.RandomName("cert", 8, testutil.CharSetAlphaNum),
		CertificatePEM: string(c),
		PrivateKeyPEM:  string(p),
	})

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetCertificate(), *actual)
}

func TestCreate_failed(t *testing.T) {
	expected := apprun_test.Fake400Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	c, p, err := apprun_test.OreSign()
	assert.NoError(err)

	actual, err := api.Create(t.Context(), CreateParams{
		Name:           testutil.RandomName("cert", 8, testutil.CharSetAlphaNum),
		CertificatePEM: string(c),
		PrivateKeyPEM:  string(p),
	})

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetCertificateResponse
	cid := v1.CertificateID(uuid.New())
	expected.SetFake()
	expected.Certificate = apprun_test.FakeCertificate()
	expected.Certificate.SetCertificateID(cid)
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), cid)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(cid, actual.GetCertificateID())
}

func TestRead_failed(t *testing.T) {
	expected := apprun_test.Fake404Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	certID := v1.CertificateID(uuid.New())
	actual, err := api.Read(t.Context(), certID)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestUpdate(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)
	cid := v1.CertificateID(uuid.New())
	c, p, err := apprun_test.OreSign()
	assert.NoError(err)

	err = api.Update(t.Context(), cid, UpdateParams{
		Name:           testutil.RandomName("cert", 8, testutil.CharSetAlphaNum),
		CertificatePEM: string(c),
		PrivateKeyPEM:  string(p),
	})

	assert.NoError(err)
}

func TestUpdate_failed(t *testing.T) {
	expected := apprun_test.Fake400Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))
	cid := v1.CertificateID(uuid.New())
	c, p, err := apprun_test.OreSign()
	assert.NoError(err)

	err = api.Update(t.Context(), cid, UpdateParams{
		Name:           testutil.RandomName("cert", 8, testutil.CharSetAlphaNum),
		CertificatePEM: string(c),
		PrivateKeyPEM:  string(p),
	})

	assert.Error(err)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	certID := v1.CertificateID(uuid.New())
	err := api.Delete(t.Context(), certID)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	expected := apprun_test.Fake404Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	certID := v1.CertificateID(uuid.New())
	err := api.Delete(t.Context(), certID)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}
