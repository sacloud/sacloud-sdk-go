# sacloud-sdk-go/api/iaas/trace/otel

さくらのクラウド IaaS API 向け OpenTelemetry トレーシング

## 概要

`api/iaas` の API クライアントに OpenTelemetry によるトレーシングを提供します。

`Initialize()` を呼び出すとクライアントのファクトリにフックが登録され、これ以降に生成された API クライアントはすべての API 呼び出しを span として記録します。

span には以下の属性が記録されます。

- API 名やゾーン、リクエストパラメータ(`sacloud.api.arguments.*`)
- レスポンスの内容(`libiaas.api.results.*`)
- エラー発生時のステータスとメッセージ

## 使い方

`Initialize()` を呼び出すだけで有効になります。

```go
package main

import (
	traceotel "github.com/sacloud/sacloud-sdk-go/api/iaas/trace/otel"
)

func main() {
	// トレーシングを有効化
	traceotel.Initialize()

	// これ以降に生成した API クライアントの呼び出しがトレースされる
}
```

具体的なコードは以下を参照してください。

- [examples/main.go](examples/main.go): Jaeger などへ OTLP/gRPC で送信する例
- [example_test.go](example_test.go): `stdouttrace` で出力する例

## ファイル構成

| ファイル | 説明 |
| --- | --- |
| `tracer.go` | 設定(`Option`)と初期化処理 |
| `functions.go` | `Initialize()` などのヘルパー |
| `zz_api_tracer.go` | 各 API 操作のトレーサー実装(自動生成) |

`zz_api_tracer.go` は `api/iaas/internal/tools/gen-otel-tracer` による自動生成ファイルです。直接編集しないでください。再生成する場合は `api/iaas` ディレクトリで `make gen` を実行します。

## License

Copyright (C) 2026 The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](../../../../LICENSE).
