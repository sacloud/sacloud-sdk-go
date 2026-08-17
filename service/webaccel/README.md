# sacloud-sdk-go/service/webaccel

[ウェブアクセラレータ](https://www.sakura.ad.jp/services/cdn/)向け高レベルAPIライブラリ

## 概要

ウェブアクセラレータのAPIをラップし、CRUD+L+Action操作を統一的な手順で行えるインターフェースを提供します。

> [!WARNING]
> sacloud-sdk-go/service/webaccelは現在開発中です。

関連プロジェクト:
- [sacloud-sdk-go/api/webaccel](../../api/webaccel)
- [sacloud-sdk-go/internal/services](../../internal/services)

インターフェースの例:
```go
import (
	"context"

	"github.com/sacloud/sacloud-sdk-go/api/webaccel"
	"github.com/sacloud/sacloud-sdk-go/service/webaccel/site"
)

// サイト操作の例
func (s *site.Service) Find(req *site.FindRequest) ([]*webaccel.Site, error)
func (s *site.Service) FindWithContext(ctx context.Context, req *site.FindRequest) ([]*webaccel.Site, error)

func (s *site.Service) Read(req *site.ReadRequest) (*webaccel.Site, error)
func (s *site.Service) ReadWithContext(ctx context.Context, req *site.ReadRequest) (*webaccel.Site, error)

func (s *site.Service) Update(req *site.UpdateRequest) (*webaccel.Site, error)
func (s *site.Service) UpdateWithContext(ctx context.Context, req *site.UpdateRequest) (*webaccel.Site, error)
```

以下のリソースに対応しています。

```console
.
├── cache
├── site
│   └── certificate
└── usage
```

## Installation

Use go get.

    go get github.com/sacloud/sacloud-sdk-go/service/webaccel

Then import the `webaccel` package into your own code.

    import "github.com/sacloud/sacloud-sdk-go/service/webaccel"

## License

Copyright (C) 2022-2025 The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).

