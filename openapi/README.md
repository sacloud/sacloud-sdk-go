openapi.jsonは https://manual.sakura.ad.jp/api/cloud/monitoring-suite/ からダウンロードできるJSONに以下の変更を加えたものです

```console
$ jq 'del(.paths.[].[].requestBody.content.["application/x-www-form-urlencoded", "multipart/form-data"])' openapi.json
```

以下のパッチがあたっていますが、自動生成に必要なだけで、本質的にAPIが変わっているわけではありません。

```diff
diff --git a/openapi/openapi.json b/openapi/openapi.json
index cbb2485..af7c72c 100644
--- a/openapi/openapi.json
+++ b/openapi/openapi.json
@@ -3936,7 +3936,19 @@
             "$ref": "#/components/schemas/MapKeyValueMatcher"
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
+            "enum": "#/components/schemas/EnumMatcher",
+            "map-key-exists": "#/components/schemas/MapKeyExistsMatcher",
+            "map-key-value-matcher": "#/components/schemas/MapKeyValueMatcher"
+          }
+        }
       },
       "FieldModel": {
         "enum": [
```