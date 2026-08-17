# sacloud-sdk-go/api/webaccel

[ウェブアクセラレータ](https://www.sakura.ad.jp/services/cdn/) の [API](https://manual.sakura.ad.jp/cloud/webaccel/api.html) をGo言語から扱うためのライブラリ

## Overview

sacloud-sdk-go/api/webaccelはさくらのクラウド ウェブアクセラレータ APIをGo言語から利用するためのAPIライブラリです。

#### sacloud-sdk-go/api/webaccelを利用したクライアントコードの例

```go
package example

import (
	"context"
	"log"

	"github.com/sacloud/sacloud-sdk-go/api/webaccel"
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

    go get github.com/sacloud/sacloud-sdk-go/api/webaccel

Then import the `webaccel` package into your own code.

    import "github.com/sacloud/sacloud-sdk-go/api/webaccel"

## License

Copyright (C) 2025- The sacloud/sacloud-sdk-go Authors.

This project is published under [Apache 2.0 License](LICENSE).
