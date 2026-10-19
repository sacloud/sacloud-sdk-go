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

`PriorityRank` is marked both `readOnly: true` and `required`. Under OpenAPI 3.0,
this makes it required in responses, not requests, so its omission from the
routing request examples is valid. The current ogen-generated validators still
require it in requests, causing the corresponding generated example tests to
skip. The handwritten acceptance tests explicitly set it to satisfy those
validators; this does not establish that the API requires it in requests.
