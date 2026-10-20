# OpenAPI current-difference check

`sources.json` maps OpenAPI definitions in this repository to the authoritative
YAML files linked from the [Sakura Cloud API portal](https://manual.sakura.ad.jp/api/cloud/portal/).
The `OpenAPI drift` GitHub Actions workflow is scheduled for Fridays at 09:00
JST (00:00 UTC) and can also be started manually.

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
process exits successfully when the comparison and report output complete, even
when differences are found. A non-zero status indicates an execution error, such
as a download, parsing, or report-writing failure. KMS, CloudHSM, and
SecretManager share a combined portal specification, so path and
component-set comparisons are skipped for those packages; matching component
definitions are still compared.

Treat the output as a review aid, not proof that a local definition is stale.
Some local definitions intentionally include compatibility fixes, and
formatting or documentation differences can also appear.

## Example output

The default output looks like this (excerpt from October 19, 2026; omitted
rows and details are marked with `...`). Versions and counts depend on the
current repository and portal definitions.

```text
# OpenAPI comparison summary

Compared APIs: 18; version differences: 10; APIs with structural differences: 12.

| API component | Repository version | Portal version | Status | Structural differences |
|---|---:|---:|---|---:|
| IAM | 2.1.0 | 2.1.0 | :white_check_mark: | 0 |
| Monitoring Suite | 1.1.0 | 1.3.0 | :x: | 169 |
| SimpleMQ message | 1.1.2 | 1.1.4 | :x: | 1 |
| SimpleMQ queue | 1.1.2 | 1.1.2 | :white_check_mark: | 0 |
...
| AppRun shared | 1.5.0 | 1.5.0 | :x: | 3 |
...

# OpenAPI structural differences

## Monitoring Suite
- components.schemas.Type: present locally, missing from portal
- components.schemas.Type1: present locally, missing from portal
...
```

In GitHub Actions, the summary is rendered as a Markdown table with status
icons. A component can have matching versions but still show `:x:` because
its structure differs, as with AppRun shared above.
