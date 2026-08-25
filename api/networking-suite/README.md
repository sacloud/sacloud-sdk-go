# sacloud-sdk-go/api/networking-suite

Go言語向けのさくらのクラウド ネットワークスイート APIライブラリ

## 概要

sacloud-sdk-go/api/networking-suite はさくらのクラウド ネットワークスイート APIをGo言語から利用するためのAPIライブラリです。

```
package main

import (
	"context"
	"fmt"

	networkingsuite "github.com/sacloud/sacloud-sdk-go/api/networking-suite"
	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

func main() {
	var theClient saclient.Client
	ctx := context.Background()
	client, err := networkingsuite.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	// subnet groups API: Create/Read/Update/Delete
	sgOp := networkingsuite.NewSubnetGroupsOp(client)
	sg, err := sgOp.Create(ctx, v1.CreateSubnetGroup{
		Name:                 "subnet-group 1",
		Description:          "description",
		IPv4AddressRangeCIDR: "10.0.0.0/20",
		Region:               v1.Region{Code: "is1"},
	})
	if err != nil {
		panic(err)
	}

	sg, err := sgOp.Read(ctx, srn.MustParse(sg.SRN))
	if err != nil {
		panic(err)
	}

	// subnet API: Create/Read/Update/Delete
	sOp := networkingsuite.NewSubnetsOp(client)
	s, err := sOp.Create(ctx, v1.CreateSubnet{
		Name:                 "subnet 1",
		Description:          "description",
		IPv4AddressRangeCIDR: "10.0.0.0/24",
		Zone:                 v1.Zone{Code: "is1c"},
		SubnetGroup:          v1.SakuraResourceNameRef{SRN: sg.SRN},
	})
	if err != nil {
		panic(err)
	}
}
```

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## License

Copyright (C) 2026- The sacloud/sacloud-sdk-go authors.
This project is published under [Apache 2.0 License](LICENSE).
