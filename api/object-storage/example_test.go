//go:build example

// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage_test

import (
	"context"
	"fmt"
	"strconv"

	objectstorage "github.com/sacloud/object-storage-api-go"
	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
	"github.com/sacloud/saclient-go"
)

// Example_siteAPI サイト(クラスタ)APIの利用例
func Example_siteAPI() {
	var theClient saclient.Client
	client, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}

	// サイト一覧の取得
	siteOp := objectstorage.NewSiteOp(client)
	sites, err := siteOp.List(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(sites[0].DisplayName.Value)
	// output:
	// 石狩第1サイト
}

// Example_bucketAPI バケットAPIの利用例
func Example_bucketAPI() {
	var theClient saclient.Client
	fedClient, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// サイトIDが必要なため取得しておく
	siteOp := objectstorage.NewSiteOp(fedClient)
	sites, err := siteOp.List(ctx)
	if err != nil {
		panic(err)
	}
	siteId := sites[0].ID.Value
	siteClient, err := objectstorage.NewSiteClient(&theClient, siteId)
	if err != nil {
		panic(err)
	}

	// バケットの作成
	bucketName := "your-bucket-from-api-go"
	bucketAPI := objectstorage.NewBucketOp(fedClient, siteClient)
	bucket, err := bucketAPI.Create(ctx, &objectstorage.BucketCreateParams{
		Bucket: bucketName,
		SiteId: siteId,
	})
	if err != nil {
		panic(err)
	}

	// バケットの削除
	defer func() {
		if err := bucketAPI.Delete(ctx, bucketName); err != nil {
			panic(err)
		}
	}()

	fmt.Println(bucket.Name.Value)
	// output:
	// your-bucket-from-api-go
}

// Example_accountAPI アカウントAPIの利用例
func Example_accountAPI() {
	var theClient saclient.Client
	fedClient, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// サイトIDが必要なため取得しておく
	siteOp := objectstorage.NewSiteOp(fedClient)
	sites, err := siteOp.List(ctx)
	if err != nil {
		panic(err)
	}
	siteId := sites[0].ID.Value
	siteClient, err := objectstorage.NewSiteClient(&theClient, siteId)
	if err != nil {
		panic(err)
	}

	// アカウントの参照
	accountOp := objectstorage.NewAccountOp(siteClient)
	account, err := accountOp.Read(ctx)
	if err != nil {
		panic(err)
	}

	// アクセスキーの作成
	accessKey, err := accountOp.CreateAccessKey(ctx)
	if err != nil {
		panic(err)
	}

	// アクセスキーの一覧
	accessKeys, err := accountOp.ListAccessKeys(ctx)
	if err != nil {
		panic(err)
	}
	if len(accessKeys) == 0 {
		panic("ListAccessKeys failed")
	}

	// アクセスキーの削除
	defer func() {
		if err := accountOp.DeleteAccessKey(ctx, string(accessKey.ID.Value)); err != nil {
			panic(err)
		}
	}()

	fmt.Printf("AccountCode: %t, Secret: %t", len(account.Code.Value) > 0, len(accessKey.Secret.Value) > 0)
	// output:
	// AccountCode: true, Secret: true
}

// Example_permissionAPI パーミッションAPIの利用例
func Example_permissionAPI() {
	var theClient saclient.Client
	fedClient, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// サイトIDが必要なため取得しておく
	siteOp := objectstorage.NewSiteOp(fedClient)
	sites, err := siteOp.List(ctx)
	if err != nil {
		panic(err)
	}
	siteId := sites[0].ID.Value
	siteClient, err := objectstorage.NewSiteClient(&theClient, siteId)
	if err != nil {
		panic(err)
	}

	// パーミッションの作成
	permissionOp := objectstorage.NewPermissionOp(siteClient)
	permission, err := permissionOp.Create(ctx, "foobar", v2.BucketControls{
		v2.BucketControlsItem{
			BucketName: v2.NewOptBucketName("bucket1"),
			CanRead:    v2.NewOptCanRead(true),
			CanWrite:   v2.NewOptCanWrite(true),
		},
	})
	if err != nil {
		panic(err)
	}

	// アクセスキーの作成
	accessKey, err := permissionOp.CreateAccessKey(ctx, strconv.Itoa(int(permission.ID.Value)))
	if err != nil {
		panic(err)
	}

	// パーミッション/アクセスキーの削除
	defer func() {
		if err := permissionOp.DeleteAccessKey(ctx, strconv.Itoa(int(permission.ID.Value)), string(accessKey.ID.Value)); err != nil {
			panic(err)
		}
		if err := permissionOp.Delete(ctx, strconv.Itoa(int(permission.ID.Value))); err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Permission: %s, Secret: %t", permission.DisplayName.Value, len(accessKey.Secret.Value) > 0)
	// output:
	// Permission: foobar, Secret: true
}

// Example_siteStatusAPI サイトステータスAPIの利用例
func Example_siteStatusAPI() {
	var theClient saclient.Client
	fedClient, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// サイトIDが必要なため取得しておく
	siteOp := objectstorage.NewSiteOp(fedClient)
	sites, err := siteOp.List(ctx)
	if err != nil {
		panic(err)
	}
	siteId := sites[0].ID.Value
	siteClient, err := objectstorage.NewSiteClient(&theClient, siteId)
	if err != nil {
		panic(err)
	}

	// サイトステータスの参照
	statusOp := objectstorage.NewSiteStatusOp(siteClient)
	status, err := statusOp.Read(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(status.StatusCode.Value.Status.Value)
	// output:
	// ok
}
