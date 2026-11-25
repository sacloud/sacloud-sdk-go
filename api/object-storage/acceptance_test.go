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
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	objectstorage "github.com/sacloud/sacloud-sdk-go/api/object-storage"
	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
	"github.com/sacloud/sacloud-sdk-go/common/packages/envvar"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
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

	s3Client, err := minio.New(site.S3Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(token, secret, ""),
		Region:       site.Region,
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
	site, err := siteOp.Read(ctx, sites[0].ID)
	require.NoError(t, err)
	require.NotEmpty(t, site)

	siteWithPlansOp := objectstorage.NewSiteWithPlansOp(accTestFedClient, accTestSiteClient)
	plans, err := siteWithPlansOp.ListPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, plans)

	// Site Status
	statusOp := objectstorage.NewSiteStatusOp(accTestSiteClient)
	status, err := statusOp.Read(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, status)

	quota, err := statusOp.ReadQuota(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, quota)
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
	bucketExists := true
	t.Cleanup(func() {
		if bucketExists {
			if err := bucketOp.Delete(ctx, bucketName); err != nil {
				t.Errorf("failed to delete bucket %q: %s", bucketName, err)
			}
		}
	})

	// Bucket API reads
	buckets, err := bucketOp.List(ctx)
	require.NoError(t, err)
	found := false
	for _, bucket := range buckets {
		if string(bucket.Name) == bucketName {
			found = true
			break
		}
	}
	require.True(t, found, "bucket %q is not in the bucket list", bucketName)

	bucketExtraOp := objectstorage.NewBucketExtraOp(accTestSiteClient, accTestFedClient, bucketName)
	penalty, err := bucketExtraOp.ReadPenalty(ctx)
	// This has to be 404, because the brand-new bucket has no penalty yet
	require.Nil(t, penalty)
	require.ErrorContains(t, err, "404")

	usage, err := bucketExtraOp.ReadUsage(ctx)
	require.NoError(t, err)
	require.NotNil(t, usage)

	quota, err := bucketExtraOp.ReadQuota(ctx)
	require.NoError(t, err)
	require.NotNil(t, quota)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	metering, err := objectstorage.NewSiteStatusOp(accTestSiteClient).ReadBucketMetering(ctx, bucketName, monthStart, monthStart)
	require.NoError(t, err)
	require.NotNil(t, metering)

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
	bucketExists = false
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
	bucketExists := true
	t.Cleanup(func() {
		if bucketExists {
			if err := bucketOp.Delete(ctx, bucketName); err != nil {
				t.Errorf("failed to delete bucket %q: %s", bucketName, err)
			}
		}
	})

	// Step2: バケットにアクセスできるパーミッション/アクセスキーの作成
	permissionOp := objectstorage.NewPermissionOp(accTestSiteClient)
	var permission *v2.PermissionData
	var key *v2.PermissionKeyCreateResponseBody
	permissionExists := false
	keyExists := false
	objectExists := false
	var permissionID, objectKey string
	{
		controls := v2.BucketControlsBody{
			v2.BucketControlsBodyItem{
				BucketName: v2.BucketName(bucketName),
				CanRead:    true,
				CanWrite:   true,
			},
		}
		created, err := permissionOp.Create(ctx, bucketName, controls)
		require.NoError(t, err)
		require.NotEmpty(t, created)
		permission = created
		permissionID = strconv.Itoa(int(permission.ID))
		permissionExists = true
		t.Cleanup(func() {
			if permissionExists {
				if err := permissionOp.Delete(ctx, permissionID); err != nil {
					t.Errorf("failed to delete permission %q: %s", permissionID, err)
				}
			}
		})

		read, err := permissionOp.Read(ctx, permissionID)
		require.NoError(t, err)
		require.Equal(t, permission.ID, read.ID)

		permissions, err := permissionOp.List(ctx)
		require.NoError(t, err)
		found := false
		for _, item := range permissions {
			if item.ID == permission.ID {
				found = true
				break
			}
		}
		require.True(t, found, "permission %q is not in the permission list", permissionID)

		updatedName := bucketName + "-updated"
		updated, err := permissionOp.Update(ctx, permissionID, updatedName, controls)
		require.NoError(t, err)
		require.Equal(t, updatedName, string(updated.DisplayName))

		createdKey, err := permissionOp.CreateAccessKey(ctx, permissionID)
		require.NoError(t, err)
		require.NotEmpty(t, createdKey)
		key = createdKey
		keyExists = true
		t.Cleanup(func() {
			if keyExists {
				if err := permissionOp.DeleteAccessKey(ctx, permissionID, string(key.ID)); err != nil {
					t.Errorf("failed to delete access key %q: %s", key.ID, err)
				}
			}
		})

		keys, err := permissionOp.ListAccessKeys(ctx, permissionID)
		require.NoError(t, err)
		require.Len(t, keys, 1)
		require.Equal(t, string(key.ID), string(keys[0].ID))

		readKey, err := permissionOp.ReadAccessKey(ctx, permissionID, string(key.ID))
		require.NoError(t, err)
		require.Equal(t, string(key.ID), string(readKey.ID))
	}

	// Step3: 作成したアクセスキーでバケットにアクセス
	{
		s3Client := s3Client(t, string(key.ID), string(key.Secret))

		objectKey = "foobar"
		objectBodyText := "body of s3://[bucket_name]/foobar"
		objectBody := bytes.NewBufferString(objectBodyText)

		// オブジェクトの作成
		created, err := s3Client.PutObject(ctx, bucketName, objectKey, objectBody, int64(objectBody.Len()), minio.PutObjectOptions{})
		require.NoError(t, err)
		require.NotEmpty(t, created)
		objectExists = true
		t.Cleanup(func() {
			if objectExists {
				if err := s3Client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{}); err != nil {
					t.Errorf("failed to remove object %q: %s", objectKey, err)
				}
			}
		})

		// オブジェクトの読み込み
		read, err := s3Client.GetObject(ctx, bucketName, objectKey, minio.GetObjectOptions{})
		require.NoError(t, err)
		require.NotEmpty(t, read)

		readText, err := io.ReadAll(read)
		require.NoError(t, err)
		require.NoError(t, read.Close())
		require.Equal(t, objectBodyText, string(readText))

		// オブジェクトの削除
		err = s3Client.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
		require.NoError(t, err)
		objectExists = false
	}

	// Step4: クリーンアップ
	require.NoError(t, permissionOp.DeleteAccessKey(ctx, permissionID, string(key.ID)))
	keyExists = false
	require.NoError(t, permissionOp.Delete(ctx, permissionID))
	permissionExists = false
	require.NoError(t, bucketOp.Delete(ctx, bucketName))
	bucketExists = false
}

// 暗号化やレプリケーションのテスト
func TestAccBucketExtraOperation(t *testing.T) {
	skipIfNoTestAcc(t)
	skipIfNoAPIKey(t)
	skipIfNoEnv(t, "SAKURA_KMS_KEY_ID")

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
		encryptionEnabled := true
		t.Cleanup(func() {
			if encryptionEnabled {
				if err := beOp.DisableEncryption(ctx); err != nil {
					t.Errorf("failed to disable bucket encryption: %s", err)
				}
			}
		})

		enc, err := beOp.ReadEncryption(ctx)
		require.NoError(t, err)
		require.Equal(t, keyId, string(enc.KmsKeyID))

		err = beOp.DisableEncryption(ctx)
		require.NoError(t, err)
		encryptionEnabled = false
	}

	// レプリケーション
	{
		beOp := objectstorage.NewBucketExtraOp(sc1, accTestFedClient, bucketName1)
		_, err := beOp.EnableReplication(ctx, bucketName2)
		require.NoError(t, err)
		replicationEnabled := true
		t.Cleanup(func() {
			if replicationEnabled {
				if err := beOp.DisableReplication(ctx); err != nil {
					t.Errorf("failed to disable bucket replication: %s", err)
				}
			}
		})

		rep, err := beOp.ReadReplication(ctx)
		require.NoError(t, err)
		require.Equal(t, tky, rep.DestBucket.ClusterID)
		require.Equal(t, v2.BucketName(bucketName2), rep.DestBucket.Name)

		err = beOp.DisableReplication(ctx)
		require.NoError(t, err)
		replicationEnabled = false
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
