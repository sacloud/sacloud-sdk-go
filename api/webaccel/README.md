# webaccel-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/webaccel-api-go.svg)](https://pkg.go.dev/github.com/sacloud/webaccel-api-go)
[![Tests](https://github.com/sacloud/webaccel-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/webaccel-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/webaccel-api-go)](https://goreportcard.com/report/github.com/sacloud/webaccel-api-go)

[ウェブアクセラレータ](https://www.sakura.ad.jp/services/cdn/) の [API](https://manual.sakura.ad.jp/cloud/webaccel/api.html) をGo言語から扱うためのライブラリ

## Overview

従来はiaas-api-go(libsacloud v2)で提供していたAPIライブラリを独立したリポジトリとして切り出したものです。  

#### webaccel-api-goを利用したクライアントコードの例

```go
package example

import (
	"context"
	"log"

	"github.com/sacloud/webaccel-api-go"
)

func Example() {
	// デフォルトではusacloudプロファイルや環境変数が利用される。
	// パラメータを指定することで上書きしたり無効化したりできる
	client := &webaccel.Client{
		//Profile:           "default",
		//AccessToken:       "token",
		//AccessTokenSecret: "secret",
		//DisableProfile:    false,
		//DisableEnv:        false,
	}
	op := webaccel.NewOp(client)

	// サイト一覧
	found, err := op.List(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println(found)

	// 全キャッシュ削除
	deleteAllCacheRequest := &webaccel.DeleteAllCacheRequest{
		Domain: "example.com",
	}
	if err := op.DeleteAllCache(context.Background(), deleteAllCacheRequest); err != nil {
		panic(err)
	}

	// URLごとにキャッシュ削除
	deleteCacheRequest := &webaccel.DeleteCacheRequest{
		URL: []string{
			"https://example.com/url1",
			"https://example.com/url2",
		},
	}
	if _, err := op.DeleteCache(context.Background(), deleteCacheRequest); err != nil {
		panic(err)
	}
}
```

## Installation

Use go get.

    go get github.com/sacloud/webaccel-api-go

Then import the `webaccel` package into your own code.

    import "github.com/sacloud/webaccel-api-go"

## License

`webaccel-api-go` Copyright 2022-2025 The webaccel-api-go authors.

This project is published under [Apache 2.0 License](LICENSE).
