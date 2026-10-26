NoSQL OpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/nosql/

現在はv2.1.1を利用しています。

# diff

この実装では不具合やogenへの対応のために定義を変更しています。
OpenAPI定義もしくはAPI側が修正され次第更新します。

- patternにanyOfがあることで `Feature "sum types with same names" is not implemented yet.`が発生する問題
- データベースの削除時のレスポンスにおいて、レスポンスがおかしくnullableじゃないDiskのEncryptionKeyフィールドに`null`が返ってくる問題
- `NosqlCreateResponse`に実際は存在しないフィールド群があり、そのうちのIDがrequiredに指定されている問題

```diff
--- openapi-orig.json	2026-10-23 17:48:27
+++ openapi/openapi.json	2026-10-26 09:40:04
@@ -1200,14 +1200,7 @@
                 "type": "string",
                 "nullable": true,
                 "description": "データベースバージョン  \n**新規作成時必須**\n",
-                "anyOf": [
-                  {
-                    "pattern": "^\\d+\\.\\d+\\.\\d+$"
-                  },
-                  {
-                    "pattern": "^$"
-                  }
-                ],
+                "pattern": "^\\d+\\.\\d+\\.\\d+$",
                 "example": "4.1.9"
               },
               "DefaultUser": {
@@ -1216,14 +1209,7 @@
                 "description": "デフォルトユーザ名  \n**新規作成時必須**\n",
                 "example": "defaultuser01",
                 "maxLength": 20,
-                "anyOf": [
-                  {
-                    "pattern": "^[a-z][a-z0-9_]{3,19}$"
-                  },
-                  {
-                    "pattern": "^$"
-                  }
-                ]
+                "pattern": "^[a-z][a-z0-9_]{3,19}$"
               },
               "DiskSize": {
                 "type": "integer",
@@ -1477,9 +1463,11 @@
                   "EncryptionKey": {
                     "type": "object",
                     "description": "暗号化キー情報",
+                    "nullable": true,
                     "properties": {
                       "KMSKeyID": {
                         "type": "string",
+                        "nullable": true,
                         "description": "KMSキーID",
                         "example": "113700349294"
                       }
@@ -1744,99 +1732,7 @@
         "properties": {
           "Appliance": {
             "$ref": "#/components/schemas/NosqlAppliance"
-          },
-          "Class": {
-            "type": "string",
-            "description": "クラス",
-            "example": "nosql"
-          },
-          "Name": {
-            "type": "string",
-            "description": "NoSQLの名前",
-            "example": "CassandraName",
-            "minLength": 1,
-            "maxLength": 64
           },
-          "Description": {
-            "type": "string",
-            "description": "NoSQLの説明",
-            "example": "説明",
-            "minLength": 0,
-            "maxLength": 512
-          },
-          "Plan": {
-            "$ref": "#/components/schemas/Plan"
-          },
-          "Settings": {
-            "$ref": "#/components/schemas/NosqlSettings"
-          },
-          "Remark": {
-            "$ref": "#/components/schemas/NosqlRemark"
-          },
-          "ID": {
-            "type": "string",
-            "example": "113600097295"
-          },
-          "Account": {
-            "type": "object",
-            "properties": {
-              "ID": {
-                "type": "string"
-              }
-            }
-          },
-          "Tags": {
-            "$ref": "#/components/schemas/Tags"
-          },
-          "Availability": {
-            "$ref": "#/components/schemas/Availability"
-          },
-          "ServerCount": {
-            "type": "integer",
-            "example": 1
-          },
-          "HiddenRemark": {
-            "type": "object",
-            "properties": {
-              "PlanSpec": {
-                "type": "object",
-                "properties": {
-                  "Note": {
-                    "type": "object",
-                    "properties": {
-                      "ID": {
-                        "type": "string"
-                      }
-                    }
-                  },
-                  "ServiceClass": {
-                    "type": "string",
-                    "example": "cloud/nosql/plan/1"
-                  }
-                }
-              },
-              "Encrypted": {
-                "type": "object",
-                "properties": {
-                  "Algorithm": {
-                    "type": "string"
-                  },
-                  "IV": {
-                    "type": "string"
-                  },
-                  "md5": {
-                    "type": "string"
-                  },
-                  "Associative": {
-                    "type": "boolean"
-                  },
-                  "Data": {
-                    "type": "string"
-                  }
-                }
-              }
-            }
-          },
           "Success": {
             "$ref": "#/components/schemas/Success"
           },
@@ -1846,7 +1742,6 @@
         },
         "required": [
           "Appliance",
-          "ID"
         ]
       },
       "NosqlUpdateRequest": {

```