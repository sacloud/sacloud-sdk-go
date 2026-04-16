# Copyright 2026- The sacloud/apprun-dedicated-api-go authors
# SPDX-License-Identifier: Apache-2.0

output "service_principal_id" {
  value = sakura_iam_service_principal.main.id
  description = "ID of the created or updated service principal, which is eligible for invoking apprun. This should be set to the environment variable SAKURA_APPRUN_DEDICATED_SERVICE_PRINCIPAL_ID for the tests to work."
}