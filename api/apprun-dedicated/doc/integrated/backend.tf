# Copyright 2026- The sacloud/apprun-dedicated-api-go authors
# SPDX-License-Identifier: Apache-2.0

terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
}