# OpenAPI definition

Source: [Sakura Cloud API portal](https://manual.sakura.ad.jp/api/cloud/portal/),
[OpenAPI YAML](https://manual.sakura.ad.jp/api/cloud/portal/openapis/simple-notification-api.yaml).

The checked-in definition is the unmodified portal YAML for `シンプル通知 API`,
`info.version: "1.0.0"`, retrieved on October 19, 2026.

Regenerate the client from the repository root:

```sh
make -C api/simple-notification ogen
```

The portal definition changes operation and schema names without changing
`info.version`. Settings unions and provider types are now generated separately
for responses, creation requests, and update requests. Request `Tags` no longer
specify `nullable: true`; use an empty slice for an empty tag list.

The routing request examples omit the required `PriorityRank` field, so ogen's
generated example tests skip those examples. The SDK's routing tests explicitly
set `PriorityRank`.
