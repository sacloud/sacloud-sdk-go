# object-storage-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/object-storage-api-go.svg)](https://pkg.go.dev/github.com/sacloud/object-storage-api-go)
[![Tests](https://github.com/sacloud/object-storage-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/object-storage-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/object-storage-api-go)](https://goreportcard.com/report/github.com/sacloud/object-storage-api-go)

Go言語向けのさくらのクラウド オブジェクトストレージAPIライブラリ

オブジェクトストレージAPIドキュメント: [https://manual.sakura.ad.jp/cloud/objectstorage/api/api-json.html](https://manual.sakura.ad.jp/cloud/objectstorage/api/api-json.html)

## 概要

sacloud/object-storage-api-goはさくらのクラウド オブジェクトストレージAPIをGo言語から利用するためのAPIライブラリです。  

- 概要/設計/実装方針: [docs/overview.md](https://github.com/sacloud/object-storage-api-go/blob/main/docs/design/overview.md)

利用イメージ:

```go
import (
    "context"
    "os"
	
    objectstorage "github.com/sacloud/object-storage-api-go"
    "github.com/sacloud/saclient-go"
)

func main() {
	// デフォルトでusacloud互換プロファイル or 環境変数(SAKURA_ACCESS_TOKEN{_SECRET}が利用される
    var theClient saclient.Client
	ctx := context.Background()
	// サイトに依存しない処理にはFedClientを利用
	fedClient, err := objectstorage.NewFedClient(&theClient)
	if err != nil {
		panic(err)
	}

	// サイト一覧を取得
	siteOp := objectstorage.NewSiteOp(fedClient)
	sites, err := siteOp.List(ctx)
	if err != nil {
		panic(err)
	}
	siteId := sites[0].ID.Value

	// サイトに依存する処理にはSiteClientを利用
	siteClient, err := objectstorage.NewSiteClient(&theClient, siteId)
	if err != nil {
		panic(err)
	}

	// バケットの作成
	bucketName := "your-bucket-name"
	bucketOp := objectstorage.NewBucketOp(fedClient, siteClient)
	bucket, err := bucketOp.Create(ctx, &objectstorage.BucketCreateParams{
		Bucket: bucketName,
		SiteId: siteId,
	})

	// バケットの削除
	defer func() {
		if err := bucketOp.Delete(ctx, bucketName); err != nil {
			panic(err)
		}
	}()

	fmt.Println(bucket.Name.Value)
}
```

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

### 関連プロジェクト

- [sacloud/sacloud-go](https://github.com/sacloud/sacloud-go): sacloud/object-storage-api-goを用いた高レベルAPIライブラリ

## License

`sacloud/object-storage-api-go` Copyright (C) 2022-2026 [The sacloud/object-storage-api-go Authors](AUTHORS).

This project is published under [Apache 2.0 License](LICENSE).
