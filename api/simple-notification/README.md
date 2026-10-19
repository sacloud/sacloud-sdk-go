# sacloud-sdk-go/api/simple-notification

さくらのクラウド シンプル通知 Go言語向け APIライブラリ

マニュアル: https://manual.sakura.ad.jp/cloud/appliance/simplenotification/


## 概要
sacloud-sdk-go/api/simple-notificationはさくらのクラウド シンプル通知 APIをGo言語から利用するためのAPIライブラリです。

OpenAPI定義の取得元とクライアントの再生成手順は
[openapi/README.md](openapi/README.md)を参照してください。

## 利用イメージ

### 通知先＆通知グループ登録
```go
import (
	"context"
	"os"

	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	simplenotification "github.com/sacloud/sacloud-sdk-go/api/simple-notification"
	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
)

func main() {
	// デフォルトでusacloud互換プロファイル or 環境変数(SAKURA_ACCESS_TOKEN{_SECRET})が利用される
	var theClient saclient.Client
	ctx := context.Background()

	simplenotificationClient, err := simplenotification.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	//使用する通知先のメールアドレスは環境変数から取得
	mailAddress := os.Getenv("SAKURA_SIMPLE_NOTIFICATION_MAILADDRESS")

	//　通知先を作成
	destAPI := simplenotification.NewDestinationOp(simplenotificationClient)
	destReq := v1.PostCommonServiceItemRequest{
		CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
			Name:        "my-destination",
			Description: "my destination description",
			Tags:        []string{"tag1", "tag2"},
			Settings: v1.PostCommonServiceItemRequestCommonServiceItemSettings{
				Type: v1.CommonServiceItemDestinationSettingsPostCommonServiceItemRequestCommonServiceItemSettings,
				CommonServiceItemDestinationSettings: v1.CommonServiceItemDestinationSettings{
					Type:  v1.CommonServiceItemDestinationSettingsTypeEmail,
					Value: mailAddress,
				},
			},
			Icon: v1.NilCommonServiceItemIcon{
				Null: false,
				Value: v1.CommonServiceItemIcon{
					ID: "112901627732", //Debian icon ID
				},
			},
		},
	}
	destResp, err := destAPI.Create(ctx, destReq)
	if err != nil {
		panic(err)
	}
	println("Created Destination ID:", destResp.CommonServiceItem.ID)

	//通知グループを作成
	groupAPI := simplenotification.NewGroupOp(simplenotificationClient)
	groupReq := v1.PostCommonServiceItemRequest{
		CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
			Name:        "my-group",
			Description: "my group description",
			Tags:        []string{"tag1", "tag2"},
			Settings: v1.PostCommonServiceItemRequestCommonServiceItemSettings{
				Type: v1.CommonServiceItemGroupSettingsPostCommonServiceItemRequestCommonServiceItemSettings,
				CommonServiceItemGroupSettings: v1.CommonServiceItemGroupSettings{
					Destinations: []string{destResp.CommonServiceItem.ID}, //先ほど作成した通知先をグループに追加
				},
			},
		},
	}
	groupResp, err := groupAPI.Create(ctx, groupReq)
	if err != nil {
		panic(err)
	}
	println("Created Group ID:", groupResp.CommonServiceItem.ID)
}
```

### 通知メッセージの送信(通知先がメールの場合は事前に本登録手続きが必要)
```go
import (
	"context"
	"fmt"
	"os"

	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	simplenotification "github.com/sacloud/sacloud-sdk-go/api/simple-notification"
	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
)

func main() {
	// デフォルトでusacloud互換プロファイル or 環境変数(SAKURA_ACCESS_TOKEN{_SECRET})が利用される
	var theClient saclient.Client
	ctx := context.Background()

	simplenotificationClient, err := simplenotification.NewClient(&theClient)
	if err != nil {
		panic(err)
	}
	groupAPI := simplenotification.NewGroupOp(simplenotificationClient)

	//使用する送信先のグループIDは環境変数から取得
	groupID := os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUPID")

	// 通知メッセージを送信
	// なお、通知先がメールの場合は、事前にメールされる本登録手続きが必要
	messageReq := v1.SendNotificationRequest{
		Message: "Hello, Simple Notification API!", //送信するメッセージ
	}
	resp, err := groupAPI.SendMessage(ctx, groupID, messageReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent successfully", resp)
}
```
> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## 実アカウントでの確認

専用のテストアカウントを使用してください。通知先グループのテストでは
実際にメッセージを送信するため、受信を確認できるテスト用の通知先が必要です。

### 事前準備

1. コントロールパネルでテスト用の通知先を作成します。メールの場合は、
   届いたメールから本登録を済ませ、その通知先IDを控えてください。
2. その通知先を登録したテスト用の通知先グループを作成し、グループIDを控えます。
3. `RoutingAPI.ListSource` が返す `Sources` から、使用する通知元の `id` と
   `name` を確認します。通知先IDやグループIDとは別のIDです。
   2026年10月19日の確認時には、`id: "2"` はモニタリングスイートでした。
   実行するアカウントで取得した一覧を使用してください。

通知先と通知先グループの操作は、公式マニュアルの
[通知先](https://manual.sakura.ad.jp/cloud/appliance/simplenotification/destination.html)、
[通知先グループ](https://manual.sakura.ad.jp/cloud/appliance/simplenotification/destination_group.html)を参照してください。

テスト内でも通知先・グループを作成しますが、各テストの終了時に削除され、
他のテストには引き継がれません。以下の環境変数には、事前に作成した
リソースのIDを指定してください。事前に用意したリソースはテストでは削除しません。

### 環境変数と実行方法

| 環境変数 | 用途 |
|---|---|
| `SAKURA_ACCESS_TOKEN` / `SAKURA_ACCESS_TOKEN_SECRET` | APIの認証情報 |
| `SAKURA_RATE_LIMIT` | APIリクエストのレート制限。例: `5`（`0` はエラー） |
| `SAKURA_DESTINATION_TEST_MAILADDRESS` | 作成するメール通知先のアドレス |
| `SAKURA_DESTINATION_TEST_ID` | 通知先グループに追加する既存の通知先ID |
| `SAKURA_ROUTING_TEST_LABEL` / `SAKURA_ROUTING_TEST_LABELVAL` | 通知ルーティングの照合ラベル名・値 |
| `SAKURA_ROUTING_TEST_SOURCEID` | `RoutingAPI.ListSource` の `Sources[].id` |
| `SAKURA_ROUTING_TEST_TARGETGROUPID` | 通知ルーティングの送信先グループID |

認証情報を環境変数に設定したうえで、次の `<...>` を実際の値に置き換えて実行します。
ラベルの例は、このテストが作成するルーティングの `MatchLabels` に設定されます。

```sh
export SAKURA_DESTINATION_TEST_MAILADDRESS='<テスト用メールアドレス>'
export SAKURA_DESTINATION_TEST_ID='<本登録済みの通知先ID>'
export SAKURA_ROUTING_TEST_TARGETGROUPID='<事前に作成した通知先グループID>'
export SAKURA_ROUTING_TEST_SOURCEID='<通知元一覧で確認したID>'
export SAKURA_ROUTING_TEST_LABEL='environment'
export SAKURA_ROUTING_TEST_LABELVAL='development'

SAKURA_RATE_LIMIT=5 make -C api/simple-notification testacc TESTARGS='-count=1'
```

`-count=1` はGoのテストキャッシュを無効にし、実際にAPIを呼び出すための指定です。
ラベル名は英数字・`_`・`-` の64文字以内、値は半角の表示可能文字と
スペースの64文字以内です。このテストではルーティングの設定操作を確認し、
通知元からルーティングを経由した配送までは確認しません。

通知先・通知先グループ・通知ルーティングの作成、取得、更新、削除に加え、
通知送信や通知履歴の取得などを行います。
必要な環境変数がないテストはスキップされるため、スキップした場合は
実際のAPI動作を確認できていません。
通知履歴が空の場合、履歴の個別取得もスキップされます。

送信APIの成功は配送完了を意味しません。通知先への到着や通知履歴も確認してください。
実行時には、受付成功後に履歴へ `notification group not found` が記録された例もあります。
現在のテストは送信受付後にグループを削除するため、配送前に削除された可能性も
考慮してください。
削除やクリーンアップに失敗した場合は、テストが作成したリソースが残っていないか
コントロールパネルで確認してください。

## License

Copyright (C) 2026- The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).