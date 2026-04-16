NoSQL OpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/nosql/

現在はv2.0.1を利用しています。

# diff

この実装では不具合やogenへの対応のために定義を変更しています。
OpenAPI定義もしくはAPI側が修正され次第更新します。

```diff
diff --git a/openapi-orig.json b/openapi/openapi.json
index 80bf3e0..e7e7377 100644
--- a/openapi-orig.json
+++ b/openapi/openapi.json
@@ -1163,14 +1163,7 @@
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
@@ -1179,14 +1172,7 @@
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
@@ -1440,9 +1426,11 @@
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
@@ -1708,98 +1696,6 @@
           "Appliance": {
             "$ref": "#/components/schemas/NosqlAppliance"
           },
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
-          },
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
@@ -1809,7 +1705,8 @@
         },
         "required": [
           "Appliance",
-          "ID"
+          "Success",
+          "is_ok"
         ]
       },
       "NosqlUpdateRequest": {

```