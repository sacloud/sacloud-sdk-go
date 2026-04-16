# Copyright 2026- The sacloud/apprun-dedicated-api-go authors
# SPDX-License-Identifier: Apache-2.0

data "sakura_iam_service_principal" "name" {
  id = var.service_principal_id
}

data "sakura_iam_project" "main" {
  id = data.sakura_iam_service_principal.name.project_id
}

resource "random_id" "main" {
  byte_length = 8
  prefix      = "apprun-dedicated-integration-test-"
  keepers = {
    project_id = data.sakura_iam_project.main.id
  }
}

resource "sakura_iam_service_principal" "main" {
  project_id  = data.sakura_iam_project.main.id
  name        = random_id.main.b64_url
  description = "Service Principal for Apprun Dedicated API integration tests"
}

resource "sakura_iam_policy" "main" {
  target    = "project"
  target_id = data.sakura_iam_project.main.id

  bindings = [
    {
      role = {
        type = "preset"
        id   = "resource-creator"
      }
      principals = [
        {
          type = "service-principal"
          id   = sakura_iam_service_principal.main.id
        }
      ]
    }
  ]
}