# sacloud/simple-notification-api-go
[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/simple-notification-api-go.svg)](https://pkg.go.dev/github.com/sacloud/simple-notification-api-go)
[![Tests](https://github.com/sacloud/simple-notification-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/simple-notification-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/simple-notification-api-go)](https://goreportcard.com/report/github.com/sacloud/simple-notification-api-go)

さくらのクラウド シンプル通知 Go言語向け APIライブラリ

マニュアル: https://manual.sakura.ad.jp/cloud/appliance/simplenotification/


## 概要
sacloud/simple-notification-api-goはさくらのクラウド シンプル通知 APIをGo言語から利用するためのAPIライブラリです。

## 利用イメージ

### 通知先＆通知グループ登録
```
import (
	"context"
	"os"

	"github.com/sacloud/saclient-go"
	simplenotification "github.com/sacloud/simple-notification-api-go"
	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
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
	mailAddress := os.Getenv("SAKURA_DESTINATION_EMAIL_ADDRESS")

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
					Type:  v1.CommonServiceItemDestinationSettingsTypeWebhook,
					Value: mailAddress,
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
```
import (
	"context"
	"fmt"
	"os"

	"github.com/sacloud/saclient-go"
	simplenotification "github.com/sacloud/simple-notification-api-go"
	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
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
	messageReq := v1.SendNotificationMessageRequest{
		Message: "Hello, Simple Notification API!", //送信するメッセージ
	}
	resp, err := groupAPI.SendMessage(ctx, groupID, messageReq)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent successfully", resp)

}
```
⚠️ v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## License

`simple-notification-api-go` Copyright (C) 2026- The sacloud/simple-notification-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).