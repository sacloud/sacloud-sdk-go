# sacloud-sdk-go/api/kms

Go言語向けのさくらのクラウド KMS APIライブラリ

KMS ドキュメント: https://manual.sakura.ad.jp/cloud/appliance/kms/index.html

## 概要

sacloud-sdk-go/api/kmsはさくらのクラウド KMS APIをGo言語から利用するためのAPIライブラリです。

```go
package main

import (
    "context"
    "fmt"

    "github.com/sacloud/sacloud-sdk-go/common/saclient"
    "github.com/sacloud/sacloud-sdk-go/api/kms"
    v1 "github.com/sacloud/sacloud-sdk-go/api/kms/apis/v1"
)

func main() {
	var theClient saclient.Client
	client, err := kms.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	keyOp := kms.NewKeyOp(client)
	// 自動生成のケース
	res, err := keyOp.Create(ctx, v1.CreateKey{
		Name:        "App key",
		Description: v1.NewOptString("key gen from go client"),
		KeyOrigin:   v1.KeyOriginEnumGenerated,
		Tags:        []string{"App1", "Key1"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Name)

	cipher, err := keyOp.Encrypt(ctx, res.ID, []byte("hello world!"), v1.KeyEncryptAlgoEnumAes256Gcm)
	plain, err := keyOp.Decrypt(ctx, res.ID, cipher)
	// plain is "hello world!"

	// Read / Update / Delete / Rotate / ChangeStatus and more...
}
```

[keys_test.go](./keys_test.go) も参照。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## ogenによるコード生成

以下のコマンドを実行

```
$ go get -tool github.com/ogen-go/ogen/cmd/ogen@latest
$ go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.json
```

## License

Copyright (C) 2025- The sacloud/sacloud-sdk-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。