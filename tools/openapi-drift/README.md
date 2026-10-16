# OpenAPI current-difference check

`sources.json` maps OpenAPI definitions in this repository to the authoritative
YAML files linked from the [Sakura Cloud API portal](https://manual.sakura.ad.jp/api/cloud/portal/).
The `OpenAPI drift` GitHub Actions workflow runs every Monday and can also be
started manually.

Run the comparison from the repository root:

```sh
go run ./tools/openapi-drift
```

It starts with a summary table showing each component's repository and portal
`info.version`, difference status, and number of structural differences. The
status is `:x:` when versions or structures differ, and `:white_check_mark:`
otherwise. GitHub Actions also displays the report in the workflow run's job
summary, without requiring you to open step logs. Detailed
differences in component definitions, schemas, and paths follow the table. The
process exits with a non-zero status when differences are found. KMS, CloudHSM,
and SecretManager share a combined portal specification, so path and
component-set comparisons are skipped for those packages; matching component
definitions are still compared.

Treat the output as a review aid, not proof that a local definition is stale.
Some local definitions intentionally include compatibility fixes, and
formatting or documentation differences can also appear.
