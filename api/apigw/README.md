# sacloud/apigw-api-go

Go言語向けのさくらのクラウド APIゲートウェイ APIライブラリ

APIゲートウェイ ドキュメント: https://manual.sakura.ad.jp/cloud/appliance/api-gateway/index.html

## 概要

sacloud/apigw-api-goはさくらのクラウド APIゲートウェイ APIをGo言語から利用するためのAPIライブラリです。

```go
import (
	"context"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
)

func main() {
	client, err := apigw.NewClient()
	ctx := context.Background()

	serviceOp := apigw.NewServiceOp(client)
	service, err := serviceOp.Create(ctx, &v1.ServiceDetail{
		Name:     "test-service",
		Host:     "example.sakura.ad.jp",
		Port:     v1.NewOptInt(80),
		Protocol: "http",
	})
	if err != nil {
		// エラー処理
	}
	defer func() { _ = serviceOp.Delete(ctx, service.ID.Value) }()

	routeOp := apigw.NewRouteOp(client, service.ID.Value)
	route, err := routeOp.Create(ctx, &v1.RouteDetail{
		Name:      v1.NewOptName("test-route"),
		Methods:   []v1.HTTPMethod{v1.HTTPMethodGET, v1.HTTPMethodPOST},
		Hosts:     []string{service.RouteHost.Value},
		Protocols: v1.NewOptRouteDetailProtocols(v1.RouteDetailProtocolsHTTPHTTPS),
		Tags:      []string{"Test"},
	})
	if err != nil {
		// エラー処理
	}
	defer func() { _ = routeOp.Delete(ctx, route.ID.Value) }()

	// サブスクリプションに関する操作
	subscriptionOp := apigw.NewSubscriptionOp(client)
	// ユーザに関する操作
	userOp := apigw.NewUserOp(client)
	// ユーザに関する追加設定。所属グループや認証の設定
	userExtraOp := apigw.NewUserExtraOp(client, user.ID.Value)
	// グループに関する操作
	groupOp := apigw.NewGroupOp(client)
	// Routeに関する追加設定。認可やリクエスト・レスポンス変換
	routeExtraOp := apigw.NewRouteExtraOp(client, service.ID.Value, route.ID.Value)
	// ドメインに関する操作
	domainOp := apigw.NewDomainOp(client)
	// 証明書に関する操作
	certOp := apigw.NewCertificateOp(client)
}
```

各 `xxx_test.go` も参照。

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## ogenによるコード生成

以下のコマンドを実行

```
$ go get -tool github.com/ogen-go/ogen/cmd/ogen@latest
$ go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.json
```

## TODO

- OIDC機能の実装

## License

`apigw-api-go` Copyright (C) 2025- The sacloud/apigw-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).