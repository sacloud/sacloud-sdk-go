// Copyright 2022-2025 The webaccel-api-go authors
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

package webaccel_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/packages-go/pointer"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/webaccel-api-go"
	"github.com/stretchr/testify/require"
)

func checkEnv(t *testing.T, requireEnvs ...string) {
	if !testutil.IsAccTest() {
		t.Skip("environment variables required: TESTACC")
	}
	testutil.PreCheckEnvsFunc("SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET")(t)
	testutil.PreCheckEnvsFunc(requireEnvs...)(t)
}

func testClient() webaccel.API {
	return webaccel.NewOp(&webaccel.Client{
		Options: &client.Options{},
	})
}

// TestScenario_Op_Create_Enable_Disable_Delete
// サイトの状態により実行結果が変化するメソッドのテストシナリオ
// 実行順序: Create -> Enable -> Disable -> Delete
func TestScenario_Op_Create_Enable_Disable_Delete(t *testing.T) {
	checkEnv(t)

	client := testClient()
	name := testutil.RandomName("webaccel-api-go-test-", 8, testutil.CharSetAlpha)
	created, err := client.Create(context.Background(), &webaccel.CreateSiteRequest{
		Name:            name,
		DomainType:      "subdomain",
		OriginType:      webaccel.OriginTypesWebServer,
		Origin:          "docs.usacloud.jp",
		OriginProtocol:  webaccel.OriginProtocolsHttps,
		VarySupport:     webaccel.VarySupportEnabled,
		DefaultCacheTTL: pointer.NewInt(3600),
		NormalizeAE:     webaccel.NormalizeAEBrGz,
	})

	require.NoError(t, err)
	require.Equal(t, created.Name, name)
	require.Equal(t, created.VarySupport, webaccel.VarySupportEnabled)
	require.Equal(t, created.NormalizeAE, webaccel.NormalizeAEBrGz)
	require.Equal(t, created.DefaultCacheTTL, 3600)
	require.NotEmpty(t, created.ID)

	site, err := client.UpdateStatus(context.Background(), created.ID, &webaccel.UpdateSiteStatusRequest{
		Status: "enabled",
	})
	require.NoError(t, err)
	require.Equal(t, site.Status, "enabled")
	site, err = client.UpdateStatus(context.Background(), created.ID, &webaccel.UpdateSiteStatusRequest{
		Status: "disabled",
	})
	require.NoError(t, err)
	require.Equal(t, site.Status, "disabled")

	deleted, err := client.Delete(context.Background(), created.ID)

	require.NoError(t, err)
	require.Equal(t, deleted.ID, created.ID)
}

func TestOp_List(t *testing.T) {
	checkEnv(t)

	client := testClient()
	found, err := client.List(context.Background())
	require.NoError(t, err)

	if found.Count == 0 {
		t.Skip("webaccel doesn't have any sites")
	}

	site := found.Sites[0]
	require.NotEmpty(t, site.ID)
	require.NotEmpty(t, site.Name)
	require.NotEmpty(t, site.DomainType)
	require.NotEmpty(t, site.Domain)
	require.NotEmpty(t, site.Subdomain)
	require.NotEmpty(t, site.ASCIIDomain)
	require.NotEmpty(t, site.Origin)
	require.NotEmpty(t, site.Status)
	require.NotEmpty(t, site.CreatedAt)
}

func TestOp_Read(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_SITE_ID")

	client := testClient()
	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
	read, err := client.Read(context.Background(), siteId)

	require.NoError(t, err)
	require.Equal(t, read.ID, siteId)
}

func TestOp_Update(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_SITE_ID")

	client := testClient()
	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
	name := testutil.RandomName("webaccel-api-go-test-", 8, testutil.CharSetAlpha)
	updated, err := client.Update(context.Background(), siteId, &webaccel.UpdateSiteRequest{
		Name:              name,
		VarySupport:       webaccel.VarySupportDisabled,
		CORSRules:         &[]*webaccel.CORSRule{},
		OnetimeURLSecrets: &[]string{},
		DefaultCacheTTL:   pointer.NewInt(0),
		NormalizeAE:       webaccel.NormalizeAEGz,
	})

	require.NoError(t, err)
	require.Equal(t, updated.Name, name)
	require.Equal(t, updated.VarySupport, webaccel.VarySupportDisabled)
	require.Empty(t, updated.CORSRules)
	require.Empty(t, updated.OnetimeURLSecrets)
	require.Equal(t, updated.DefaultCacheTTL, 0)
}

func TestOp_OriginGuardToken(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_SITE_ID")

	var (
		siteId           string
		token, nextToken *webaccel.OriginGuardTokenResponse
		err              error
	)
	client := testClient()

	t.Run("create origin guard token", func(t *testing.T) {
		siteId = os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
		token, err = client.CreateOriginGuardToken(context.Background(), siteId)
		require.NoError(t, err)
		require.NotEqual(t, len(token.OriginGuardToken), 0)
		require.Equal(t, len(token.NextOriginGuardToken), 0)
	})
	t.Run("read origin guard token", func(t *testing.T) {
		res, err := client.Read(context.Background(), siteId)
		require.NoError(t, err)
		require.NotEqual(t, len(res.OriginGuardToken), 0)
		require.Equal(t, res.OriginGuardToken, token.OriginGuardToken)
	})

	t.Run("create next origin guard token", func(t *testing.T) {
		nextToken, err = client.CreateNextOriginGuardToken(context.Background(), siteId)
		require.NoError(t, err)
		require.NotEqual(t, len(nextToken.OriginGuardToken), 0)
		require.NotEqual(t, len(nextToken.NextOriginGuardToken), 0)
		require.Equal(t, nextToken.OriginGuardToken, token.OriginGuardToken)
	})
	t.Run("migrate to the next origin guard token", func(t *testing.T) {
		token, err = client.CreateOriginGuardToken(context.Background(), siteId)
		require.NoError(t, err)
		require.NotEqual(t, len(token.OriginGuardToken), 0)
		require.Equal(t, len(token.NextOriginGuardToken), 0)
		require.Equal(t, token.OriginGuardToken, nextToken.NextOriginGuardToken)
	})

	t.Run("delete next origin guard token", func(t *testing.T) {
		err = client.DeleteNextOriginGuardToken(context.Background(), siteId)
		require.NoError(t, err)
	})

	t.Run("delete origin guard token", func(t *testing.T) {
		err = client.DeleteOriginGuardToken(context.Background(), siteId)
		require.NoError(t, err)
	})
}

// NOTE: Let's Encrypt 自動更新が無効化されたサイトを用意しないと失敗する
func TestOp_AutoCertUpdate(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_SITE_ID")
	client := testClient()
	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")

	t.Run("create automatic certificate update configuration", func(t *testing.T) {
		err := client.CreateAutoCertUpdate(context.Background(), siteId)
		require.NoError(t, err)
		time.Sleep(time.Second)
	})
	t.Run("delete automatic certificate update configuration", func(t *testing.T) {
		err := client.DeleteAutoCertUpdate(context.Background(), siteId)
		require.NoError(t, err)
	})
}

func TestWebAccelOp_ACL(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_SITE_ID")

	client := testClient()
	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
	ctx := context.Background()

	t.Run("create ACL", func(t *testing.T) {
		acl := "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all"
		result, err := client.UpsertACL(ctx, siteId, acl)

		require.NoError(t, err)
		require.Equal(t, acl, result.ACL)
	})
	t.Run("read ACL", func(t *testing.T) {
		acl := "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all"
		result, err := client.ReadACL(ctx, siteId)

		require.NoError(t, err)
		require.Equal(t, acl, result.ACL)
	})
	t.Run("update ACL", func(t *testing.T) {
		acl := "allow 192.0.2.5/25\nallow 198.51.100.0\ndeny all"
		result, err := client.UpsertACL(ctx, siteId, acl)

		require.NoError(t, err)
		require.Equal(t, acl, result.ACL)
	})
	t.Run("delete ACL", func(t *testing.T) {
		if err := client.DeleteACL(ctx, siteId); err != nil {
			t.Fatal("got unexpected error", err)
		}

		result, err := client.ReadACL(ctx, siteId)
		require.NoError(t, err)
		require.Empty(t, result.ACL)
	})
}

func TestWebAccelOp_Cert(t *testing.T) {
	envKeys := []string{
		"SAKURACLOUD_WEBACCEL_SITE_ID",
		"SAKURACLOUD_WEBACCEL_CERT",
		"SAKURACLOUD_WEBACCEL_KEY",
		"SAKURACLOUD_WEBACCEL_CERT_UPD",
		"SAKURACLOUD_WEBACCEL_KEY_UPD",
	}
	checkEnv(t, envKeys...)

	client := testClient()
	ctx := context.Background()
	id := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
	crt := os.Getenv("SAKURACLOUD_WEBACCEL_CERT")
	key := os.Getenv("SAKURACLOUD_WEBACCEL_KEY")
	crtUpd := os.Getenv("SAKURACLOUD_WEBACCEL_CERT_UPD")
	keyUpd := os.Getenv("SAKURACLOUD_WEBACCEL_KEY_UPD")

	// create certs
	_, err := client.CreateCertificate(ctx, id, &webaccel.CreateOrUpdateCertificateRequest{
		CertificateChain: crt,
		Key:              key,
	})
	require.NoError(t, err)

	// update certs
	certs, err := client.UpdateCertificate(ctx, id, &webaccel.CreateOrUpdateCertificateRequest{
		CertificateChain: crtUpd,
		Key:              keyUpd,
	})
	require.NoError(t, err)

	// read cert
	read, err := client.ReadCertificate(ctx, id)
	require.NoError(t, err)

	require.Equal(t, certs, read)

	// delete certs
	err = client.DeleteCertificate(ctx, id)
	require.NoError(t, err)

	// read again
	read, err = client.ReadCertificate(ctx, id)
	require.NoError(t, err)
	require.Empty(t, read.Current)
	require.NotEmpty(t, read.Old)
}

func TestOp_DeleteAllCache(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURACLOUD_WEBACCEL_DOMAIN")(t)

	client := testClient()

	// delete cache
	err := client.DeleteAllCache(context.Background(), &webaccel.DeleteAllCacheRequest{
		Domain: os.Getenv("SAKURACLOUD_WEBACCEL_DOMAIN"),
	})
	require.NoError(t, err)
}

func TestOp_DeleteCache(t *testing.T) {
	checkEnv(t, "SAKURACLOUD_WEBACCEL_URLS")

	client := testClient()
	result, err := client.DeleteCache(context.Background(), &webaccel.DeleteCacheRequest{
		URL: strings.Split(os.Getenv("SAKURACLOUD_WEBACCEL_URLS"), ","),
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)
}

// NOTE: 当月中にトラフィックのないアカウントまたは月初のテストで失敗する
func TestOp_MonthlyUsage(t *testing.T) {
	checkEnv(t)

	client := testClient()
	results, err := client.MonthlyUsage(context.Background(), "")

	require.NoError(t, err)
	require.NotEmpty(t, results.Year)
	require.NotEmpty(t, results.Month)
	require.NotEmpty(t, results.MonthlyUsages)
}

func TestOp_LogUploadConfig(t *testing.T) {
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_WEBACCEL_SITE_ID",
		"SAKURASTORAGE_BUCKET_NAME",
		"SAKURASTORAGE_ACCESS_KEY",
		"SAKURASTORAGE_ACCESS_SECRET",
	)(t)

	client := testClient()

	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")
	region := "jp-north-1"
	endpoint := "https://s3.isk01.sakurastorage.jp"
	bucketName := os.Getenv("SAKURASTORAGE_BUCKET_NAME")
	accessKey := os.Getenv("SAKURASTORAGE_ACCESS_KEY")
	accessSecret := os.Getenv("SAKURASTORAGE_ACCESS_SECRET")
	status := "enabled"

	var (
		applied, read *webaccel.LogUploadConfig
		err           error
	)

	require.NotEmpty(t, bucketName)
	require.NotEmpty(t, accessKey)
	require.NotEmpty(t, accessSecret)

	t.Run("apply the log upload configuration", func(t *testing.T) {
		applied, err = client.ApplyLogUploadConfig(context.Background(), siteId, &webaccel.LogUploadConfig{
			Region:          region,
			Endpoint:        endpoint,
			Bucket:          bucketName,
			AccessKeyID:     accessKey,
			SecretAccessKey: accessSecret,
			Status:          status,
		})

		require.NoError(t, err)
		require.Equal(t, applied.Region, region)
		require.Equal(t, applied.Endpoint, endpoint)
		require.Equal(t, applied.Bucket, bucketName)
		require.Equal(t, applied.AccessKeyID, accessKey)
		require.Equal(t, applied.SecretAccessKey, accessSecret)
		require.Equal(t, applied.Status, status)
	})

	t.Run("read a log upload configuration", func(t *testing.T) {
		read, err = client.ReadLogUploadConfig(context.Background(), siteId)

		require.NoError(t, err)
		require.Equal(t, read.Region, applied.Region)
		require.Equal(t, read.Endpoint, applied.Endpoint)
		require.Equal(t, read.Bucket, applied.Bucket)
		require.NotEmpty(t, read.AccessKeyID)
		require.NotEmpty(t, read.SecretAccessKey)
		require.Equal(t, read.Status, applied.Status)
	})

	t.Run("delete the log upload configuration", func(t *testing.T) {
		err = client.DeleteLogUploadConfig(context.Background(), siteId)
		require.NoError(t, err)
	})
}
