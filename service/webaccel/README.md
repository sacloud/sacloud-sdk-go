# webaccel-service-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/webaccel-service-go.svg)](https://pkg.go.dev/github.com/sacloud/webaccel-service-go)
[![Tests](https://github.com/sacloud/webaccel-service-go/workflows/Tests/badge.svg)](https://github.com/sacloud/webaccel-service-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/webaccel-service-go)](https://goreportcard.com/report/github.com/sacloud/webaccel-service-go)

[ウェブアクセラレータ](https://www.sakura.ad.jp/services/cdn/)向け高レベルAPIライブラリ

## 概要

ウェブアクセラレータのAPIをラップし、CRUD+L+Action操作を統一的な手順で行えるインターフェースを提供します。

:warning: webaccel-service-goは現在開発中です。

関連プロジェクト: 
- [sacloud/webaccel-api-go](https://github.com/sacloud/webaccel-api-go)
- [sacloud/services](https://github.com/sacloud/services)

インターフェースの例:
```go
// サイト操作の例
func (s *Service) Find(req *FindRequest) ([]*webaccel.Site, error)
func (s *Service) FindWithContext(ctx context.Context, req *FindRequest) ([]*webaccel.Site, error)

func (s *Service) Read(req *ReadRequest) (*webaccel.Site, error)
func (s *Service) ReadWithContext(ctx context.Context, req *ReadRequest) (*webaccel.Site, error)

func (s *Service) Update(req *UpdateRequest) (*webaccel.Site, error)
func (s *Service) UpdateWithContext(ctx context.Context, req *UpdateRequest) (*webaccel.Site, error)
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

    go get github.com/sacloud/webaccel-service-go

Then import the `webaccel` package into your own code.

    import "github.com/sacloud/webaccel-service-go"

## License

`webaccel-service-go` Copyright 2022-2025 The sacloud/webaccel-service-go authors.

This project is published under [Apache 2.0 License](LICENSE).

