# sacloud/eventbus-api-go

Go言語向けのさくらのクラウド EventBus APIライブラリ

EventBus APIドキュメント: https://manual.sakura.ad.jp/cloud/appliance/eventbus/index.html

## 概要

sacloud/eventbus-api-goはさくらのクラウド EventBus APIをGo言語から利用するためのAPIライブラリです。

```go
package main

import (
    "context"
    "fmt"
    "strconv"
    "time"

    eventbus "github.com/sacloud/eventbus-api-go"
    v1 "github.com/sacloud/eventbus-api-go/apis/v1"
)

func main() {
    client, err := eventbus.NewClient()
    if err != nil {
        panic(err)
    }
    ctx := context.Background()
    pcOp := eventbus.NewProcessConfigurationOp(client)
    schedOp := eventbus.NewScheduleOp(client)

    // テスト用の実行設定の生成 (1111111111はリソースIDの例)
    pc, err := pcOp.Create(ctx, v1.ProcessConfigurationRequestSettings{
        Name: "実行設定1", Description: "アプリ向け実行設定",
        Settings: eventbus.CreateSimpleNotificationSettings("1111111111", "Hello"),
    })
    if err != nil {
        panic(err)
    }
    pcId := strconv.FormatInt(pc.ID, 10)

    res, err := schedOp.Create(ctx, v1.ScheduleRequestSettings{
        Name: "スケジュール1", Description: "アプリ向けスケジュール",
        Settings: v1.ScheduleSettings{
            ProcessConfigurationID: pcId,
            RecurringStep:          10,
            RecurringUnit:          "min",
            StartsAt:               time.Now().UnixMilli(),
        },
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(res.Name)
}
```

[process_configurations_test.go](./process_configurations_test.go) / [schedules_test.go](./schedules_test.go) も参照。

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## License

`eventbus-api-go` Copyright (C) 2025- The sacloud/eventbus-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).
