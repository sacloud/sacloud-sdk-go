// Copyright 2022-2025 The sacloud/iaas-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otel_test

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/ostype"
	traceotel "github.com/sacloud/iaas-api-go/trace/otel"
	"github.com/sacloud/iaas-api-go/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Example ローカルのJaegerを利用する例
//
// あらかじめJaegerを起動しておくこと
//
//	$ docker run -d --name jaeger -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one:latest
func Example() {
	tp, err := tracerProvider()
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main")
	ctx, span := tr.Start(ctx, "foo")
	defer span.End()

	// サンプルAPIリクエスト
	op(ctx)

	// Jaeger UI( http://localhost:16686/search など)を開くとトレースが確認できるはず
}

func tracerProvider() (*tracesdk.TracerProvider, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
	}
	// Create the OTLP/gRPC exporter
	exp, err := otlptracegrpc.New(context.Background())
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),

		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("iaas-api-go"),
			attribute.String("version", iaas.Version),
		)),
	)
	return tp, nil
}

func op(ctx context.Context) {
	// set factory func
	traceotel.Initialize()

	caller := iaas.NewClient(
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
		os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
	)
	archiveOp := iaas.NewArchiveOp(caller)

	// normal operation
	archiveOp.Find(ctx, "is1a", &iaas.FindCondition{ // nolint
		Count:  1,
		From:   0,
		Filter: ostype.ArchiveCriteria[ostype.Ubuntu],
	})

	// invalid operation(not foundエラーになるはず)
	archiveOp.Read(ctx, "is1a", types.ID(1)) // nolint
}
