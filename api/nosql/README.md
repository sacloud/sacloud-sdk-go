# sacloud/nosql-api-go

Go言語向けのさくらのクラウド NoSQL APIライブラリ

NoSQL ドキュメント: https://manual.sakura.ad.jp/cloud/appliance/nosql/index.html

## 概要

sacloud/nosql-api-goはさくらのクラウド NoSQL APIをGo言語から利用するためのAPIライブラリです。

```go
package main

import (
	"context"
	"fmt"
	"net/netip"

	nosql "github.com/sacloud/nosql-api-go"
	v1 "github.com/sacloud/nosql-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

func main() {
	theClient := saclient.Client{}
	client, err := nosql.NewClient(&theClient)
	if err != nil {
		panic(err)
	}
    
    ctx := context.Background()
	databaseOp := nosql.NewDatabaseOp(client)
	// 以下はテスト用のセットアップ。アドレスなどは実際の環境向けに編集する
	resCreated, err := databaseOp.Create(ctx, nosql.Plan100GB, v1.NosqlCreateRequestAppliance{
		Name:        "sdk-test-db",
		Description: v1.NewOptString("This is a test database"),
		Tags:        v1.NewOptNilTags([]string{"nosql"}),
		Settings: v1.NosqlSettings{
			Backup: v1.NewOptNilNosqlSettingsBackup(v1.NosqlSettingsBackup{
				Connect:   "nfs://192.168.0.31/export",
				DayOfWeek: v1.NewOptNilNosqlSettingsBackupDayOfWeekItemArray([]v1.NosqlSettingsBackupDayOfWeekItem{"sun"}),
				Time:      v1.NewOptNilString("00:00"),
				Rotate:    2,
			}),
			SourceNetwork:    []string{},
			Password:         v1.NewOptPassword("sdktest-12345"),
			ReserveIPAddress: v1.NewOptIPv4(netip.MustParseAddr("192.168.0.10")),
			Repair: v1.NewOptNilNosqlSettingsRepair(v1.NosqlSettingsRepair{
				Incremental: v1.NewOptNosqlSettingsRepairIncremental(v1.NosqlSettingsRepairIncremental{
					DaysOfWeek: []v1.NosqlSettingsRepairIncrementalDaysOfWeekItem{"sun"},
					Time:       "01:00",
				}),
				Full: v1.NewOptNosqlSettingsRepairFull(v1.NosqlSettingsRepairFull{
					DayOfWeek: v1.NosqlSettingsRepairFullDayOfWeek("sun"),
					Time:      "02:00",
					Interval:  v1.NosqlSettingsRepairFullInterval(7),
				}),
			}),
		},
		Remark: v1.NosqlRemark{
			// NosqlRemarkNosqlの設定は現状固定値なので、DefaultUser以外は変更しない
			Nosql: v1.NosqlRemarkNosql{
				DatabaseEngine:  v1.NewOptNilNosqlRemarkNosqlDatabaseEngine("Cassandra"),
				DatabaseVersion: v1.NewOptNilString("4.1.10"),
				DefaultUser:     v1.NewOptNilString("sdktest"),
				Port:            v1.NewOptNilInt(9042),
				Storage:         v1.NewOptNilNosqlRemarkNosqlStorage("SSD"),
				Zone:            "tk1b",
			},
			Servers: []v1.NosqlRemarkServersItem{
				{UserIPAddress: netip.MustParseAddr("192.168.0.4")},
				{UserIPAddress: netip.MustParseAddr("192.168.0.5")},
				{UserIPAddress: netip.MustParseAddr("192.168.0.6")},
			},
			Network: v1.NosqlRemarkNetwork{
				DefaultRoute:   "192.168.0.1",
				NetworkMaskLen: 24,
			},
		},
		UserInterfaces: []v1.NosqlCreateRequestApplianceUserInterfacesItem{
			{
				// 実際のSwitchのリソースIDを指定する
				Switch:         v1.NosqlCreateRequestApplianceUserInterfacesItemSwitch{ID: "111111111111"},
				UserIPAddress1: netip.MustParseAddr("192.168.0.4"),
				UserIPAddress2: v1.NewOptIPv4(netip.MustParseAddr("192.168.0.5")),
				UserIPAddress3: v1.NewOptIPv4(netip.MustParseAddr("192.168.0.6")),
				UserSubnet: v1.NosqlCreateRequestApplianceUserInterfacesItemUserSubnet{
					DefaultRoute:   "192.168.0.1",
					NetworkMaskLen: 24,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Database: %+v\n", resCreated)

	res, err := databaseOp.Read(ctx, resCreated.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read Database: %+v\n", res)

	instanceOp := nosql.NewInstanceOp(client, res.ID, "tk1b")
	version, err := instanceOp.GetVersion(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version: %+v\n", version)

	// Backup API
	// backupOp := nosql.NewBackupOp(client, res.ID)
}
```

NOTE: DatabaseAPIにあるNoSQL更新APIは設定のバリデーションのみで、実際に更新するには追加で反映APIを呼ぶ必要があります: `Update` → `ApplyChanges`.

### ノードの追加

```
func main()
{
    // 省略
	instanceOp := nosql.NewInstanceOpWithZone(client, primaryNodeId, "tk1b")
	resAdded, err := instanceOp.AddNodes(ctx, nosql.Plan100GB, v1.NosqlCreateRequestAppliance{
		Name: "sdk-test-db-add",
		Settings: v1.NosqlSettings{
			ReserveIPAddress: v1.NewOptIPv4(netip.MustParseAddr("192.168.0.11")),
		},
		Remark: v1.NosqlRemark{
			Nosql: v1.NosqlRemarkNosql{
				Zone: "tk1b",
			},
			Servers: []v1.NosqlRemarkServersItem{
				{UserIPAddress: netip.MustParseAddr("192.168.0.7")},
				{UserIPAddress: netip.MustParseAddr("192.168.0.8")},
			},
			Network: v1.NosqlRemarkNetwork{
				DefaultRoute:   "192.168.0.1",
				NetworkMaskLen: 24,
			},
		},
		UserInterfaces: []v1.NosqlCreateRequestApplianceUserInterfacesItem{
			{
				Switch:         v1.NosqlCreateRequestApplianceUserInterfacesItemSwitch{ID: "111111111111"},
				UserIPAddress1: netip.MustParseAddr("192.168.0.7"),
				UserIPAddress2: v1.NewOptIPv4(netip.MustParseAddr("192.168.0.8")),
				UserSubnet: v1.NosqlCreateRequestApplianceUserInterfacesItemUserSubnet{
					DefaultRoute:   "192.168.0.1",
					NetworkMaskLen: 24,
				},
			},
		},
	})
}
```

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## ogenによるコード生成

以下のコマンドを実行

```
$ go get -tool github.com/ogen-go/ogen/cmd/ogen@latest
$ go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.json
```

## License

`nosql-api-go` Copyright (C) 2025- The sacloud/nosql-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).