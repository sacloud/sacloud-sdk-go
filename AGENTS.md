# Agent Instructions for the Sacloud SDK for Go

This document describes the project architecture and guidelines for AI agents working with the Sacloud SDK for Go project.

## About this project

Sacloud SDK for Go is a monorepo containing Go SDK libraries for [さくらのクラウド](https://cloud.sakura.ad.jp/) services. It provides both low-level API clients and high-level service libraries for interacting with Sakura Cloud's various services.

## Architectural Overview

The project follows a layered architecture with clear separation of concerns:

```
sacloud-sdk-go/
├── api/                   # Low-level API client libraries
│   ├── addon/             # Add-on API (Data Lake, WAF, etc.)
│   ├── apigw/             # API Gateway
│   ├── apprun/            # AppRun
│   ├── apprun-dedicated/  # AppRun Dedicated
│   ├── cloudhsm/          # CloudHSM
│   ├── dedicated-storage/ # Dedicated Storage
│   ├── eventbus/          # EventBus
│   ├── iaas/              # IaaS (Servers, Disks, Networks, etc.)
│   ├── iam/               # IAM (Authentication & Authorization)
│   ├── kms/               # Key Management Service
│   ├── monitoring-suite/  # Monitoring
│   ├── nosql/             # NoSQL
│   ├── object-storage/    # Object Storage
│   ├── secretmanager/     # Secret Manager
│   ├── security-control/  # Security Control
│   ├── service-endpoint-gateway/
│   ├── simple-notification/
│   ├── simplemq/          # SimpleMQ
│   ├── webaccel/          # Web Accelerator
│   └── workflows/         # Workflows
├── service/               # High-level service libraries
│   ├── iaas/              # IaaS high-level API
│   └── webaccel/          # Web Accelerator high-level API
├── common/                # Shared public packages
│   ├── packages/          # General utility packages
│   └── saclient/          # Authentication & configuration
├── internal/              # Shared internal packages
│   ├── api-client/        # Common API client implementation
│   ├── go-http/           # HTTP communication layer
│   └── services/          # Service layer common implementation
├── srn/                   # SRN (Sakura Resource Name) helpers
└── makefiles/             # Shared Makefile recipes
```

### Layer Descriptions

| Layer | Purpose | Example Package |
|-------|---------|-----------------|
| `api/*` | Low-level REST API clients | `github.com/sacloud/sacloud-sdk-go/api/iaas` |
| `service/*` | High-level abstractions over API clients | `github.com/sacloud/sacloud-sdk-go/service/iaas` |
| `common/*` | Shared public packages (utilities, auth/config) | `github.com/sacloud/sacloud-sdk-go/common/packages` |
| `internal/*` | Shared implementation details | Internal use only |
| `srn/` | Sakura Resource Name parsing and helpers | `github.com/sacloud/sacloud-sdk-go/srn` |

### Nested Modules

Most packages in this repository share the root module (`github.com/sacloud/sacloud-sdk-go`). A small number of optional sub-modules are defined for dependency isolation. For example:

- `api/iaas/trace/otel` — OpenTelemetry tracing integration for the IaaS API client (`github.com/sacloud/sacloud-sdk-go/api/iaas/trace/otel`).

## License Header (Required for all files)

- This repo is transitioning from verbose license headers to a concise format.
- For new files, use the concise header below, but do **not** copy `YYYY-` literally: replace it with the current year when creating the file.
- If there are already files alongside the new file, follow their existing year style and module name.

```go
// Copyright YYYY- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0
```

