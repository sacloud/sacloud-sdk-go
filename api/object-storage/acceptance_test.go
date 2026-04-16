// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	objectstorage "github.com/sacloud/object-storage-api-go"
	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
	"github.com/sacloud/packages-go/envvar"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var siteId = envvar.StringFromEnv("SAKURA_OJS_SITE", "isk01")
var theClient saclient.Client
var accTestFedClient = initFedClient()
var accTestSiteClient = initSiteClient(siteId)

func s3Client(t *testing.T, token, secret string) *minio.Client {
	t.Helper()

	siteOp := objectstorage.NewSiteOp(accTestFedClient)
	site, err := siteOp.Read(context.Background(), siteId)
	if err != nil {
		t.Fatal(err)
	}

	s3Client, err := minio.New(site.S3Endpoint.Value, &minio.Options{
		Creds:        credentials.NewStaticV4(token, secret, ""),
		Region:       site.Region.Value,
		Secure:       true,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	return s3Client
}

func s3ClientFromEnv(t *testing.T) *minio.Client {
	return s3Client(t, os.Getenv("SAKURA_OJS_ACCESS_TOKEN"), os.Getenv("SAKURA_OJS_ACCESS_TOKEN_SECRET"))
}

// TestAccSiteAndStatusAPI サイト/ステータス関連APIの疎通確認
func TestAccSiteAndStatusAPI(t *testing.T) {
	skipIfNoTestAcc(t)
	skipIfNoAPIKey(t)

	ctx := context.Background()
	siteOp := objectstorage.NewSiteOp(accTestFedClient)

	// List
	sites, err := siteOp.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, sites)

	// Read
	site, err := siteOp.Read(ctx, sites[0].ID.Value)
	require.NoError(t, err)
	require.NotEmpty(t, site)

	// Site Status
	statusOp := objectstorage.NewSiteStatusOp(accTestSiteClient)
	status, err := statusOp.Read(ctx)

	require.NoError(t, err)
	require.NotEmpty(t, status)
}

// TestAccBucketHandling バケット周りの一連の操作のテスト
//
// 一連の操作は以下のステップで行う
//   - Step1: バケットを作成
//   - Step2: アクセスキーを用いてAWS SDK経由でバケットが作成されたかを確認
//   - Step3: 各リソースのクリーンアップ
//
// Note: バケット一覧の参照のためにサイトアカウントのアクセスキーが必要
func TestAccBucketHandling(t *testing.T) {
	skipIfNoTestAcc(t)
	skipIfNoAPIKey(t)
	skipIfNoEnv(t, "SAKURA_OJS_ACCESS_TOKEN", "SAKURA_OJS_ACCESS_TOKEN_SECRET")

	ctx := context.Background()
	bucketName := "api-go-acc-" + testutil.Random(28, testutil.CharSetAlpha)

	// Step1: バケット作成
	bucketOp := objectstorage.NewBucketOp(accTestFedClient, accTestSiteClient)
	{
		created, err := bucketOp.Create(ctx, &objectstorage.BucketCreateParams{SiteId: siteId, Bucket: bucketName})
		require.NoError(t, err)
		require.NotEmpty(t, created)
	}

	// Step2: バケットにアクセス
	{
		s3Client := s3ClientFromEnv(t)

		output, err := s3Client.ListBuckets(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, output)

		exist := false
		for _, b := range output {
			if b.Name == bucketName {
				exist = true
				break
			}
		}
		require.True(t, exist, "bucket %q is not exist", bucketName)
	}

	// Step3: クリーンアップ
	require.NoError(t, bucketOp.Delete(ctx, bucketName))
}

// TestAccAccessToBucketWithPermissionKey パーミッションキーによるオブジェクトへのアクセス
func TestAccAccessToBucketObjectWithPermissionKey(t *testing.T) {
	skipIfNoTestAcc(t)
	skipIfNoAPIKey(t)

	ctx := context.Background()
	bucketName := "api-go-acc-" + testutil.Random(28, testutil.CharSetAlpha)

	// Step1: バケット作成
	bucketOp := objectstorage.NewBucketOp(accTestFedClient, accTestSiteClient)
	{
		created, err := bucketOp.Create(ctx, &objectstorage.BucketCreateParams{SiteId: siteId, Bucket: bucketName})
		require.NoError(t, err)
		require.NotEmpty(t, created)
	}

	// Step2: バケットにアクセスできるパーミッション/アクセスキーの作成
	permissionOp := objectstorage.NewPermissionOp(accTestSiteClient)
	var permission *v2.PermissionData
	var key *v2.PermissionKeyData
	{
		created, err := permissionOp.Create(ctx, bucketName, v2.BucketControls{
			v2.BucketControlsItem{
				BucketName: v2.NewOptBucketName(v2.BucketName(bucketName)),
				CanRead:    v2.NewOptCanRead(true),
				CanWrite:   v2.NewOptCanWrite(true),
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, created)
		permission = created

		createdKey, err := permissionOp.CreateAccessKey(ctx, strconv.Itoa(int(permission.ID.Value)))
		require.NoError(t, err)
		require.NotEmpty(t, createdKey)
		key = createdKey
	}

	// Step3: 作成したアクセスキーでバケットにアクセス
	{
		s3Client := s3Client(t, string(key.ID.Value), string(key.Secret.Value))

		objectKey := "foobar"
		objectBodyText := "body of s3://[bucket_name]/foobar"
		objectBody := bytes.NewBufferString(objectBodyText)

		// オブジェクトの作成
		created, err := s3Client.PutObject(ctx, bucketName, objectKey, objectBody, int64(objectBody.Len()), minio.PutObjectOptions{})
		require.NoError(t, err)
		require.NotEmpty(t, created)

		// オブジェクトの読み込み
		read, err := s3Client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
		require.NoError(t, err)
		require.NotEmpty(t, read)

		readText, err := io.ReadAll(read)
		require.NoError(t, err)
		require.Equal(t, objectBodyText, string(readText))

		// オブジェクトの削除
		err = s3Client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
		require.NoError(t, err)
	}

	// Step4: クリーンアップ
	require.NoError(t, permissionOp.DeleteAccessKey(ctx, strconv.Itoa(int(permission.ID.Value)), string(key.ID.Value)))
	require.NoError(t, permissionOp.Delete(ctx, strconv.Itoa(int(permission.ID.Value))))
	require.NoError(t, bucketOp.Delete(ctx, bucketName))
}

// 暗号化やレプリケーションのテスト
func TestAccBucketExtraOperation(t *testing.T) {
	skipIfNoTestAcc(t)
	skipIfNoAPIKey(t)

	ctx := context.Background()
	isk := "isk01"
	tky := "tky01"
	sc1 := initSiteClient(isk)
	sc2 := initSiteClient(tky)
	bucketName1 := "api-go-acc-" + testutil.Random(28, testutil.CharSetAlpha)
	bucketName2 := "api-go-acc-" + testutil.Random(28, testutil.CharSetAlpha)

	b1Op := objectstorage.NewBucketOp(accTestFedClient, sc1)
	{
		created, err := b1Op.Create(ctx, &objectstorage.BucketCreateParams{SiteId: isk, Bucket: bucketName1})
		require.NoError(t, err)
		require.NotEmpty(t, created)
	}
	t.Cleanup(func() {
		if err := b1Op.Delete(ctx, bucketName1); err != nil {
			t.Logf("failed to delete bucket %q: %s", bucketName1, err)
		}
	})

	b2Op := objectstorage.NewBucketOp(accTestFedClient, sc2)
	{
		created, err := b2Op.Create(ctx, &objectstorage.BucketCreateParams{SiteId: tky, Bucket: bucketName2})
		require.NoError(t, err)
		require.NotEmpty(t, created)
	}
	t.Cleanup(func() {
		if err := b2Op.Delete(ctx, bucketName2); err != nil {
			t.Logf("failed to delete bucket %q: %s", bucketName2, err)
		}
	})

	// 暗号化
	{
		beOp := objectstorage.NewBucketExtraOp(sc2, accTestFedClient, bucketName2)
		keyId := os.Getenv("SAKURA_KMS_KEY_ID")
		err := beOp.EnableEncryption(ctx, keyId)
		require.NoError(t, err)

		enc, err := beOp.ReadEncryption(ctx)
		require.NoError(t, err)
		require.Equal(t, keyId, string(enc.KmsKeyID.Value))

		err = beOp.DisableEncryption(ctx)
		require.NoError(t, err)
	}

	// レプリケーション
	{
		beOp := objectstorage.NewBucketExtraOp(sc1, accTestFedClient, bucketName1)
		_, err := beOp.EnableReplication(ctx, bucketName2)
		require.NoError(t, err)

		rep, err := beOp.ReadReplication(ctx)
		require.NoError(t, err)
		require.Equal(t, tky, rep.DestBucket.ClusterID.Value)
		require.Equal(t, bucketName2, rep.DestBucket.Name.Value)

		err = beOp.DisableReplication(ctx)
		require.NoError(t, err)
	}
}

func initFedClient() *objectstorage.FedClient {
	client, err := objectstorage.NewFedClientWithAPIRootURL(&theClient, envvar.StringFromEnv("SAKURA_OJS_ROOT_URL", objectstorage.DefaultAPIRootURL))
	if err != nil {
		panic(err)
	}
	return client
}

func initSiteClient(sid string) *objectstorage.SiteClient {
	client, err := objectstorage.NewSiteClientWithAPIRootURL(&theClient, envvar.StringFromEnv("SAKURA_OJS_ROOT_URL", objectstorage.DefaultAPIRootURL), sid)
	if err != nil {
		panic(err)
	}
	return client
}

// skipIfNoEnv 指定の環境変数のいずれかが空の場合はt.SkipNow()する
func skipIfNoEnv(t *testing.T, envs ...string) {
	var emptyEnvs []string
	for _, env := range envs {
		if os.Getenv(env) == "" {
			emptyEnvs = append(emptyEnvs, env)
		}
	}
	if len(emptyEnvs) > 0 {
		for _, env := range emptyEnvs {
			t.Logf("environment variable %q is not set", env)
		}
		t.SkipNow()
	}
}

func skipIfNoAPIKey(t *testing.T) {
	skipIfNoEnv(t, "SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")
}

func skipIfNoTestAcc(t *testing.T) {
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}
}
