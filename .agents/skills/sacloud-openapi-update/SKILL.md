---
name: sacloud-openapi-update
description: Update a Sakura Cloud API package from the OpenAPI definition published in the Sakura Cloud API portal. Use when asked to update, synchronize, regenerate, or follow an OpenAPI specification, including SimpleMQ.
compatibility: GitHub Copilot CLI and OpenCode
---

# Sakura Cloud OpenAPI SDK update

Use this skill to update an `api/<service>` package when its Sakura Cloud
OpenAPI definition changes. The checked-in definition and the generated `ogen`
client are source artifacts: update them together, and do not manually edit
`oas_*_gen.go` files.

## 1. Identify the authoritative specification

1. Open <https://manual.sakura.ad.jp/api/cloud/portal/>.
2. Find the API card for the requested service and use its **OpenAPI YAML** or
   **JSON** link. Do not derive the file name from the service name: portal API
   IDs are not always the package directory name.
3. Record the direct `openapis/<api-id>.yaml` or `.json` URL, the `info.title`,
   and the `info.version`. Download the definition into a temporary directory
   first with `curl --fail --location --show-error`; this prevents a failed
   request from replacing the checked-in definition.
4. Read `api/<service>/openapi/README.md`, `api/<service>/Makefile`, and any
   `patch/` directory before replacing a definition. The README documents
   service-specific sources or generator limitations. Preserve deliberate,
   checked-in compatibility edits.

For SimpleMQ, the two package definitions map to the following portal URLs:

| Package file | API purpose | Source URL |
| --- | --- | --- |
| `api/simplemq/openapi/message.yaml` | Queue message send/receive API | `https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-api.yaml` |
| `api/simplemq/openapi/queue.yaml` | Sakura Cloud queue management API | `https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-sacloud-api.yaml` |

Example of a safe SimpleMQ download and replacement:

```bash
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

curl --fail --location --show-error \
  --output "$tmpdir/message.yaml" \
  https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-api.yaml
curl --fail --location --show-error \
  --output "$tmpdir/queue.yaml" \
  https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-sacloud-api.yaml

cmp -s "$tmpdir/message.yaml" api/simplemq/openapi/message.yaml ||
  cp "$tmpdir/message.yaml" api/simplemq/openapi/message.yaml
cmp -s "$tmpdir/queue.yaml" api/simplemq/openapi/queue.yaml ||
  cp "$tmpdir/queue.yaml" api/simplemq/openapi/queue.yaml
```

Inspect the diff before generation. Verify that it is a definition update, not
an HTML error page or an unexpected unrelated API.

## 2. Regenerate the API client

1. Treat the package `Makefile` as the authoritative generation command.
   Most generated API packages use `go tool ogen` and keep its configuration in
   `ogen-config.yaml` or `ogen-config.yml`.
2. Run the package's generation target from the repository root:

   ```bash
   make -C api/<service> gen
   ```

   For SimpleMQ, this generates both `apis/v1/queue` and `apis/v1/message`,
   then applies `patch/01_list_filter.patch`. That patch is required because
   the queue-list filter cannot currently be represented in its OpenAPI
   definition. Do not omit it or hand-edit the generated output afterward.
3. If `ogen` reports an unsupported or invalid construct, first compare it
   with existing service-specific workarounds in `openapi/README.md` and
   `patch/`. Keep an auditable workaround as either a documented derived
   definition or a patch applied by `make gen`; never leave an unrepeatable
   edit in generated code.
4. Review the generated diff for changed operations, request fields, response
   variants, requiredness, defaults, authentication schemes, and endpoints.
   Generated type or operation changes may require corresponding updates to
   the handwritten package facade, such as `client.go`, `queue.go`, and
   `message.go` in SimpleMQ.

## 3. Verify the update

Run the narrow package checks after reconciling handwritten code with the
generated interfaces:

```bash
make -C api/<service> test
go -C api/<service> build ./...
git diff --check
```

For SimpleMQ:

```bash
make -C api/simplemq test
go -C api/simplemq build ./...
git diff --check
```

`make test` runs the package Go tests with the race detector. Compilation and
generated encode/decode tests do not prove authentication, endpoint behavior,
or response compatibility, so run real-account integration tests when
confirming an OpenAPI update.

Do not run real-account integration tests yourself: the required credentials
are commonly unavailable to the agent. After completing the implementation
and the local checks, ask the user to run the package's existing
acceptance-test target with a dedicated test account:

```bash
SAKURA_ACCESS_TOKEN=... \
SAKURA_ACCESS_TOKEN_SECRET=... \
make -C api/<service> testacc
```

For SimpleMQ, tell the user that this test creates, updates, lists, reads,
and deletes a queue, then sends, receives, extends, and deletes a message.
It requires `SAKURA_ACCESS_TOKEN` and `SAKURA_ACCESS_TOKEN_SECRET` for the
dedicated test account. If the test fails during cleanup, ask the user to
confirm that no `SDK-Test-Queue` resource remains.

Before finishing, confirm all of the following:

- Each downloaded portal definition is tracked in `api/<service>/openapi/`.
- `make gen` succeeds from a clean generated-client directory and reproduces
  every required patch.
- Handwritten public interfaces and error handling cover newly generated
  operations or response statuses when the API change exposes them.
- The diff contains only the definition, generated output, intentional
  handwritten compatibility updates, and directly related tests or docs.

## Final step: request real-account verification

In the final response, after reporting the implementation and local-check
results, ask the user to run the following command with a dedicated test
account to verify the actual API behavior:

```bash
make -C api/<service> testacc
```

State the required environment variables for the target package without
requesting or exposing their values. Explain that skipped integration tests
do not confirm actual API behavior. For SimpleMQ, also ask the user to
confirm that no `SDK-Test-Queue` resource remains when cleanup fails.
