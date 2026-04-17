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
├── internal/              # Shared internal packages
│   ├── api-client/        # Common API client implementation
│   ├── go-http/           # HTTP communication layer
│   ├── packages/          # General utility packages
│   ├── saclient/          # Authentication & configuration
│   └── services/          # Service layer common implementation
└── makefiles/             # Shared Makefile recipes
```

### Layer Descriptions

| Layer | Purpose | Example Package |
|-------|---------|-----------------|
| `api/*` | Low-level REST API clients | `github.com/sacloud/iaas-api-go` |
| `service/*` | High-level abstractions over API clients | `github.com/sacloud/iaas-service-go` |
| `internal/*` | Shared implementation details | Internal use only |

