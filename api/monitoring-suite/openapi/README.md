openapi.jsonは https://manual.sakura.ad.jp/api/cloud/monitoring-suite/ からダウンロードできるJSONに以下の変更を加えたものです

```console
$ jq 'del(.paths.[].[].requestBody.content.["application/x-www-form-urlencoded", "multipart/form-data"])' openapi.json
```

以下のパッチがあたっていますが、自動生成に必要なだけで、本質的にAPIが変わっているわけではありません。

```diff
diff --git a/openapi/monitoring-suite-api.json b/openapi/openapi.json
index 7e5a3d3..6ece5a3 100644
--- a/openapi/monitoring-suite-api.json
+++ b/openapi/openapi.json
@@ -4885,7 +4885,22 @@
             "$ref": "#/components/schemas/MapValueNumMatcher"
           }
         ],
-        "title": "FieldMatcher"
+        "title": "FieldMatcher",
+        "discriminator": {
+          "propertyName": "type",
+          "mapping": {
+            "or": "#/components/schemas/OrMatcher",
+            "and": "#/components/schemas/AndMatcher",
+            "string": "#/components/schemas/StrMatcher",
+            "number": "#/components/schemas/NumMatcher",
+            "boolean": "#/components/schemas/BoolMatcher",
+            "enum": "#/components/schemas/EnumMatcher",
+            "map-key-exists": "#/components/schemas/MapKeyExistsMatcher",
+            "map-key-value-matcher": "#/components/schemas/MapKeyValueMatcher",
+            "map-value-string": "#/components/schemas/MapValueStrMatcher",
+            "map-value-number": "#/components/schemas/MapValueNumMatcher"
+          }
+        }
       },
       "FieldModel": {
         "enum": [
```