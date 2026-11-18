オブジェクトストレージのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/objectstorage/

現状はv1.7.0を利用しています。それに加えてobject-storage-api-goでは、公開されている定義の不具合を解消するため下記の変更を加えています。

```diff
--- /openapi.json   2026-11-13 10:13:09.128183350 +0900
+++ api/object-storage/openapi/openapi.json     2026-11-17 17:16:38.584969981 +0900
@@ -2160,7 +2160,6 @@
             "name": "from",
             "schema": {
               "type": "string",
-              "format": "date",
               "example": "2022-03"
             },
             "description": "範囲の開始月",
@@ -2171,7 +2170,6 @@
             "name": "to",
             "schema": {
               "type": "string",
-              "format": "date",
               "example": "2022-08"
             },
             "description": "範囲の終了月",
@@ -2628,6 +2626,7 @@
         "type": "string",
         "description": "Plan type enum",
         "enum": [
+          "",
           "standard",
           "archive"
         ],
@@ -2912,6 +2911,7 @@
           },
           "fee": {
             "type": "object",
+            "nullable": true,
             "required": [
               "for_month",
               "monthly",
@@ -2942,6 +2942,7 @@
       },
       "ServiceClassPath": {
         "type": "string",
+        "minLength": 0,
         "description": "Plan identifier (service class path). Attributes such as capacity are resolved by the server based on this value.",
         "example": "objectstorage/arc02/bucket/2t"
       },
@@ -3250,13 +3251,13 @@
           },
           "from": {
             "type": "string",
-            "format": "date",
+            "format": "date-time",
             "example": "2022-08-01",
             "description": "The first day of metering range"
           },
           "to": {
             "type": "string",
-            "format": "date",
+            "format": "date-time",
             "example": "2022-08-31",
             "description": "The last day of metering range"
           }
```