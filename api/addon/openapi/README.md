`openapi.json` は https://manual.sakura.ad.jp/api/cloud/portal/openapis/addon-api.json から取得したJSONファイルに以下の変更を加えたものです。

```diff
diff --git a/api/addon/openapi/openapi.json b/api/addon/openapi/openapi.json
index 2d49b4c7..b32350ef 100644
--- a/api/addon/openapi/openapi.json
+++ b/api/addon/openapi/openapi.json
@@ -2862,6 +2862,11 @@
           },
           "data": {
             "$ref": "#/components/schemas/ResourceGroupData"
+          },
+          "url": {
+            "type": "string",
+            "format": "uri",
+            "readOnly": true
           }
         },
         "additionalProperties": false
```