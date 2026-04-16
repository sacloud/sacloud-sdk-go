APIゲートウェイのOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/apigw/

現状はv2.0.1を利用しています。
それに加えてapigw-api-goでは、公開されている定義からogenが未サポートの機能を削除した定義を利用しています。

## OpenAPI定義のdiff

以下の問題に対処するための暫定的な修正:

- 現状ogenが複雑な`anyOf`を処理できないケース
- 現状ogenが`array`に対する`default`を処理できないケース
- CorsConfig/ObjectStorageConfigに対してogenがMergeエラーを出すので、要素が1つの`allOf`を削除
- OpenAPI定義のtypo群 (typeつけわすれ/string型とinteger型の間違い/定義生成時に不具合によるおかしな指定)

```diff
diff --git a/openapi/orig-openapi.json b/openapi/openapi.json
index 222cce1..fc2974f 100644
--- a/openapi/orig-openapi.json
+++ b/openapi/openapi.json
@@ -98,6 +98,7 @@
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "services": {
                           "type": "array",
@@ -282,6 +283,7 @@
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "route": {
                           "$ref": "#/components/schemas/RouteDetail"
@@ -334,6 +336,7 @@
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "routes": {
                           "type": "array",
@@ -385,11 +388,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "route": {
                           "$ref": "#/components/schemas/RouteDetail"
@@ -560,11 +565,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "routeAuthorization": {
                           "$ref": "#/components/schemas/RouteAuthorizationDetailResponse"
@@ -652,11 +659,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "requestTransformation": {
                           "$ref": "#/components/schemas/RequestTransformation"
@@ -747,11 +756,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "responseTransformation": {
                           "$ref": "#/components/schemas/ResponseTransformation"
@@ -1488,8 +1499,6 @@
                   ],
                   "properties": {
                     "apigw": {
-                      "type": "object",
-                      "properties": null,
                       "$ref": "#/components/schemas/DomainDTO"
                     }
                   }
@@ -1663,8 +1672,6 @@
                   ],
                   "properties": {
                     "apigw": {
-                      "type": "object",
-                      "properties": null,
                       "$ref": "#/components/schemas/CertificateDTO"
                     }
                   }
@@ -1782,6 +1789,7 @@
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "plans": {
                           "type": "array",
@@ -2058,11 +2066,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "oidcs": {
                           "type": "array",
@@ -2111,11 +2121,13 @@
             "content": {
               "application/json": {
                 "schema": {
+                  "type": "object",
                   "required": [
                     "apigw"
                   ],
                   "properties": {
                     "apigw": {
+                      "type": "object",
                       "properties": {
                         "oidc": {
                           "$ref": "#/components/schemas/OidcDetail"
@@ -2384,8 +2396,6 @@
       "OidcSummary": {
         "description": "OIDC認証要約情報",
         "type": "object",
-        "allOf": [
-          {
             "properties": {
               "id": {
                 "type": "string",
@@ -2400,8 +2410,6 @@
                 "readOnly": true
               }
             }
-          }
-        ]
       },
       "Service": {
         "description": "Service情報",
@@ -2605,30 +2613,24 @@
       "ServiceSummary": {
         "description": "Service要約情報",
         "type": "object",
-        "allOf": [
-          {
-            "properties": {
-              "id": {
-                "type": "string",
-                "format": "uuid",
-                "readOnly": true,
-                "description": "Entityを識別するためのID",
-                "example": "b69f7bfe-1d8f-48bb-8b85-b47359366e48"
-              },
-              "name": {
-                "$ref": "#/components/schemas/Name",
-                "description": "Service名<br>Service名は半角英数字およびアンダースコアのみを許可",
-                "example": "serviceName"
-              }
-            }
+        "properties": {
+          "id": {
+            "type": "string",
+            "format": "uuid",
+            "readOnly": true,
+            "description": "Entityを識別するためのID",
+            "example": "b69f7bfe-1d8f-48bb-8b85-b47359366e48"
+          },
+          "name": {
+            "$ref": "#/components/schemas/Name",
+            "description": "Service名<br>Service名は半角英数字およびアンダースコアのみを許可",
+             "example": "serviceName"
           }
-        ]
+        }
       },
       "CorsConfig": {
         "description": "CORS設定",
         "type": "object",
-        "allOf": [
-          {
             "properties": {
               "credentials": {
                 "type": "boolean",
@@ -2670,17 +2672,6 @@
                   "CONNECT",
                   "TRACE"
                 ],
-                "default": [
-                  "GET",
-                  "POST",
-                  "PUT",
-                  "DELETE",
-                  "PATCH",
-                  "OPTIONS",
-                  "HEAD",
-                  "CONNECT",
-                  "TRACE"
-                ],
                 "description": "CORS許可メソッド<br>未指定の場合は全メソッドを許可"
               },
               "accessControlAllowOrigins": {
@@ -2700,14 +2691,10 @@
                 "example": false
               }
             }
-          }
-        ]
       },
       "ObjectStorageConfig": {
         "description": "Object Storage設定",
         "type": "object",
-        "allOf": [
-          {
             "properties": {
               "bucketName": {
                 "type": "string",
@@ -2763,14 +2750,10 @@
               "secretAccessKey",
               "useDocumentIndex"
             ]
-          }
-        ]
       },
       "IpRestrictionConfig": {
         "type": "object",
         "description": "IP制限",
-        "allOf": [
-          {
             "properties": {
               "protocols": {
                 "type": "string",
@@ -2810,8 +2793,6 @@
               "restrictedBy",
               "ips"
             ]
-          }
-        ]
       },
       "Route": {
         "type": "object",
@@ -2883,17 +2864,6 @@
                   "GET",
                   "POST"
                 ],
-                "default": [
-                  "GET",
-                  "POST",
-                  "PUT",
-                  "DELETE",
-                  "PATCH",
-                  "OPTIONS",
-                  "HEAD",
-                  "CONNECT",
-                  "TRACE"
-                ],
                 "description": "RouteにアクセスするためのHTTPメソッド<br>未指定の場合は全メソッドを許可<br>オブジェクトストレージ形式のサービスに紐づく場合は、GET・HEAD・OPTIONSメソッドのみを許可"
               },
               "httpsRedirectStatusCode": {
@@ -3249,32 +3219,6 @@
                 ]
               }
             }
-          },
-          {
-            "anyOf": [
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "rsa"
-                ]
-              },
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "ecdsa"
-                ]
-              },
-              {
-                "type": "object",
-                "required": [
-                  "name",
-                  "rsa",
-                  "ecdsa"
-                ]
-              }
-            ]
           }
         ]
       },
@@ -3451,9 +3395,6 @@
         "properties": {
           "isACLEnabled": {
             "type": "boolean",
-            "enum": [
-              true
-            ],
             "description": "認可設定が有効かどうか"
           },
           "groups": {
@@ -3908,8 +3849,7 @@
               },
               "price": {
                 "type": "string",
-                "minimum": 0,
-                "example": 1000,
+                "example": "1000",
                 "description": "価格"
               },
               "description": {
@@ -4017,8 +3957,6 @@
       "SubscriptionPlanResponse": {
         "description": "料金プラン",
         "type": "object",
-        "allOf": [
-          {
             "properties": {
               "planID": {
                 "type": "string",
@@ -4033,8 +3971,7 @@
               },
               "price": {
                 "type": "string",
-                "minimum": 0,
-                "example": 1000,
+                "example": "1000",
                 "description": "価格"
               },
               "maxServices": {
@@ -4066,10 +4003,9 @@
                 "$ref": "#/components/schemas/Overage"
               }
             }
-          }
-        ]
       },
       "SubscriptionList": {
+        "type": "object",
         "properties": {
           "subscriptions": {
             "type": "array",

```

## 生成されたコードのdiff

以下の問題に対処するための暫定的な修正:

- ogenが現状`writeOnly`を認識しない不具合に対するrequiredチェックの無効化
- ecdsaフィールドのレスポンスの`{}`が`[]`に変換されてしまう問題

```diff
ddiff --git a/apis/v1/oas_json_gen.go b/apis/v1/oas_json_gen.go
index a2aa88f..af4e80a 100644
--- a/apis/v1/oas_json_gen.go
+++ b/apis/v1/oas_json_gen.go
@@ -2574,7 +2574,7 @@ func (s *BasicAuth) Decode(d *jx.Decoder) error {
        // Validate required fields.
        var failures []validate.FieldError
        for i, mask := range [1]uint8{
-               0b00011000,
+               0b00001000,
        } {
                if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
                        // Mask only required fields and check equality to mask using XOR.
@@ -2920,7 +2920,13 @@ func (s *CertificateDetails) Decode(d *jx.Decoder) error {
                }
                return nil
        }); err != nil {
-               return errors.Wrap(err, "decode CertificateDetails")
+               // ecdsaを指定していない場合には{}が返ってくるはずだが、現状APIの裏側で変換する処理が入ってしまい[]が返ってくるので、
+               // 修正されるまでそれを無視する
+               if errArr := d.Arr(func(d *jx.Decoder) error { return nil }); errArr == nil {
+                       return nil
+               } else {
+                       return errors.Wrap(err, "decode CertificateDetails")
+               }
        }

        return nil
@@ -10340,7 +10340,7 @@ func (s *HmacAuth) Decode(d *jx.Decoder) error {
        // Validate required fields.
        var failures []validate.FieldError
        for i, mask := range [1]uint8{
-               0b00011000,
+               0b00001000,
        } {
                if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
                        // Mask only required fields and check equality to mask using XOR.
@@ -10779,7 +10779,7 @@ func (s *Jwt) Decode(d *jx.Decoder) error {
        // Validate required fields.
        var failures []validate.FieldError
        for i, mask := range [1]uint8{
-               0b00111000,
+               0b00101000,
        } {
                if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
                        // Mask only required fields and check equality to mask using XOR.
@@ -11060,7 +11060,7 @@ func (s *ObjectStorageConfig) Decode(d *jx.Decoder) error {
        // Validate required fields.
        var failures []validate.FieldError
        for i, mask := range [1]uint8{
-               0b01111101,
+               0b01001101,
        } {
                if result := (requiredBitSet[i] & mask) ^ mask; result != 0 {
                        // Mask only required fields and check equality to mask using XOR.
```