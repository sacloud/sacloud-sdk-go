オブジェクトストレージのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/objectstorage/

現状はv1.5.0を利用しています。それに加えてobject-storage-api-goでは、公開されている定義の不具合を解消するため下記の変更を加えています。

```diff
--- openapi.json	2026-02-02 14:00:41
+++ openapi/openapi.json	2026-02-03 11:29:02
@@ -2710,6 +2710,7 @@
           "fee": {
             "type": "object",
             "description": "Fee information (tax included)",
+            "nullable": true,
             "properties": {
               "for_month": {
                 "type": "integer",
@@ -2827,7 +2828,8 @@
                 "type": "object",
                 "properties": {
                   "val": {
-                    "type": "integer",
+                    "type": "number",
+                    "format": "float",
                     "example": 10002495
                   },
                   "quota": {
@@ -2981,6 +2983,11 @@
         "type": "integer",
         "format": "int64",
         "example": 123
+      },
+      "PermissionKeyID": {
+        "description": "Permission Key ID",
+        "type": "string",
+        "example": "AKIAIOSFODNN7EXAMPLE"
       },
       "ResourceID": {
         "description": "Resource ID",
@@ -2991,7 +2998,6 @@
       "Code": {
         "description": "Code",
         "type": "string",
-        "pattern": "^[\\w@]+$",
         "example": "abc01234@foo@isk01"
       },
       "CreatedAt": {
@@ -3003,7 +3009,7 @@
       "BucketName": {
         "description": "バケット名",
         "type": "string",
-        "pattern": "^[a-zA-Z][a-zA-Z0-9]{2,}",
+        "pattern": "^[a-zA-Z][a-zA-Z0-9\\-]{2,}",
         "example": "Abcdefg1235-"
       },
       "CanRead": {
@@ -3279,7 +3285,7 @@
                 "id": {
                   "allOf": [
                     {
-                      "$ref": "#/components/schemas/PermissionID"
+                      "$ref": "#/components/schemas/PermissionKeyID"
                     }
                   ]
                 },
@@ -3313,7 +3319,7 @@
               "id": {
                 "allOf": [
                   {
-                    "$ref": "#/components/schemas/PermissionID"
+                    "$ref": "#/components/schemas/PermissionKeyID"
                   }
                 ]
               },

```