# Copyright 2026- The sacloud/apprun-dedicated-api-go authors
# SPDX-License-Identifier: Apache-2.0

terraform {
  required_providers {
    sakura = {
      source  = "sacloud/sakura"
      version = "3.5.0" # needs IAM, ">= 3.5.0" is must.
    }
  }
}

provider "sakura" {
  # all params are set via ~/.usacloud
  profile = terraform.workspace
}
