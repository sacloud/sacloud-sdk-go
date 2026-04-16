# sacloud/secretmanager-api-go

Go言語向けのさくらのクラウド シークレットマネージャ APIライブラリ

シークレットマネージャ ドキュメント: https://manual.sakura.ad.jp/cloud/appliance/secretsmanager/index.html

## 概要

sacloud/secretmanager-api-goはさくらのクラウド KMS APIをGo言語から利用するためのAPIライブラリです。

```go
package main

import (
    "context"
    "fmt"

	"github.com/sacloud/saclient-go"
    sm "github.com/sacloud/secretmanager-api-go"
    v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

var theClient saclient.Client

func main() {
	client, err := sm.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	keyId := os.Getenv("SAKURA_KMS_KEY_ID") // コンパネやkms-api-goなどで取得
	vaultOp := sm.NewVaultOp(client)

	vault, err := vaultOp.Create(ctx, v1.CreateVault{
		Name:        "app1_vault",
		Description: v1.NewOptString("vault for app1"),
		KmsKeyID:    keyId,
		Tags:        []string{"app1"},
	})
	if err != nil {
		panic(err)
	}

	secOp := sm.NewSecretOp(client, vault.ID)

	resCreate, err := secOp.Create(ctx, v1.CreateSecret{
		Name:  "secret1",
		Value: "Secret Value 1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("version: " + strconv.Itoa(resCreate.LatestVersion))

	resList, err := secOp.List(ctx)
	if err != nil {
		panic(err)
	}

	for _, sec := range resList {
		fmt.Println("name: " + sec.Name + ", version: " + strconv.Itoa(sec.LatestVersion))
	}

	resUn, err := secOp.Unveil(ctx, v1.Unveil{
		Name: "secret1",
		//Version: v1.NewOptNilInt(1), // Versionを指定して取得も可能
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("value: " + resUn.Value)
}
```

[vaults_test.go](./vaults_test.go) / [secrets_test.go](./secrets_test.go) も参照。

### クライアントに設定を渡す

`saclient-go`側で設定を読み込み、渡します。

```go
// API keysをコードから指定する例
import (
	"os"
	// ...

	"github.com/sacloud/saclient-go"
	sm "github.com/sacloud/secretmanager-api-go"
)

var theClient saclient.Client

func main() {
	theClient.FlagSet().Parse(os.Args[1:])
	theClient.SetEnviron(os.Environ())
	theClient.SetWith(saclient.WithFavouringBearerAuthentication()) // etc.
	client, err := sm.NewClient(&theClient)
	// ...
}
```

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## ogenによるコード生成

以下のコマンドを実行

```
$ go get -tool github.com/ogen-go/ogen/cmd/ogen@latest
$ go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi-fixed.json
```

## License

`secretmanager-api-go` Copyright (C) 2025- The sacloud/secretmanager-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).