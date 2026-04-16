シークレットマネージャのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/security-encryption/#tag/secretmanager-vault

secretmanager-api-goではここで公開されている定義からCloudHSM / KSM向けの定義を削除したのを利用しています。

## OpenAPI定義のdiff

同梱されているOpenAPI定義は以下の問題を回避するために一部修正されています

- Unveilのtypo
- ogenがreadOnly/writeOnlyをうまくサポートできない問題

```diff
--- openapi.json	2025-06-30 10:19:21
+++ openapi-fixed.json	2025-06-30 10:21:49
@@ -223,7 +223,7 @@
             "content": {
               "application/json": {
                 "schema": {
-                  "$ref": "#/components/schemas/WrappedCreateSecret"
+                  "$ref": "#/components/schemas/WrappedSecret"
                 }
               }
             },
@@ -283,7 +283,7 @@
           "content": {
             "application/json": {
               "schema": {
-                "$ref": "#/components/schemas/WrappedUnvail"
+                "$ref": "#/components/schemas/WrappedUnveil"
               }
             }
           },
@@ -294,7 +294,7 @@
             "content": {
               "application/json": {
                 "schema": {
-                  "$ref": "#/components/schemas/WrappedUnvail"
+                  "$ref": "#/components/schemas/WrappedUnveil"
                 }
               }
             },
@@ -509,7 +509,7 @@
         "type": "string",
         "description": "* `cloud/cloudhsm/partition` - Type-L7"
       },
-      "Unvail": {
+      "Unveil": {
         "type": "object",
         "properties": {
           "Name": {
@@ -584,22 +584,22 @@
           "Name"
         ]
       },
-      "WrappedCreateKey": {
+      "WrappedCreateSecret": {
         "type": "object",
         "properties": {
-          "Key": {
-            "$ref": "#/components/schemas/CreateKey"
+          "Secret": {
+            "$ref": "#/components/schemas/CreateSecret"
           }
         },
         "required": [
-          "Key"
+          "Secret"
         ]
       },
-      "WrappedCreateSecret": {
+      "WrappedSecret": {
         "type": "object",
         "properties": {
           "Secret": {
-            "$ref": "#/components/schemas/CreateSecret"
+            "$ref": "#/components/schemas/Secret"
           }
         },
         "required": [
@@ -628,22 +628,11 @@
           "Vault"
         ]
       },
-      "WrappedKey": {
+      "WrappedUnveil": {
         "type": "object",
         "properties": {
-          "Key": {
-            "$ref": "#/components/schemas/Key"
-          }
-        },
-        "required": [
-          "Key"
-        ]
-      },
-      "WrappedUnvail": {
-        "type": "object",
-        "properties": {
           "Secret": {
-            "$ref": "#/components/schemas/Unvail"
+            "$ref": "#/components/schemas/Unveil"
           }
         },
         "required": [

```