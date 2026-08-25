AppRun共用型のOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/portal/?api=apprun-shared-api

現在はv1.5.0を利用しています。api/apprunではいくつかの変更をOpenAPIや生成したコードに施しています。

- OpenAPIにtypeが足りない

```
diff --git a/api/apprun/openapi/openapi.yaml b/api/apprun/openapi/openapi.yaml
index e75eb20d..2a1b7645 100644
--- a/api/apprun/openapi/openapi.yaml
+++ b/api/apprun/openapi/openapi.yaml
@@ -1947,6 +1947,7 @@ components:
           type: array
           maxItems: 10
           items:
+            type: object
             required:
               - from_ip
               - from_ip_prefix_length
@@ -1975,6 +1976,7 @@ components:
           type: array
           maxItems: 10
           items:
+            type: object
             required:
               - from_ip
               - from_ip_prefix_length
@@ -2438,6 +2440,7 @@ components:
           type: array
           maxItems: 10
           items:
+            type: object
             required:
               - from_ip
               - from_ip_prefix_length
```

- jsonのデコードでnullがうまく扱えない

```
diff --git a/api/apprun/apis/v1/oas_json_gen.go b/api/apprun/apis/v1/oas_json_gen.go
index a70f4781..24ee1cf0 100644
--- a/api/apprun/apis/v1/oas_json_gen.go
+++ b/api/apprun/apis/v1/oas_json_gen.go
@@ -5951,6 +5951,12 @@ func (s *HandlerListTrafficMeta) Decode(d *jx.Decoder) error {
                return errors.New("invalid: unable to decode HandlerListTrafficMeta to nil")
        }
 
+       // OpenAPI上ではnullableだがogenがそれを正しく処理できずにいるため、nullを許容するための特別な処理を追加。
+       if d.Next() == jx.Null {
+               d.Null()
+               return nil
+       }
+
        if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
                switch string(k) {
                default:
@@ -10714,6 +10720,12 @@ func (s *HandlerUpdateTrafficMeta) Decode(d *jx.Decoder) error {
                return errors.New("invalid: unable to decode HandlerUpdateTrafficMeta to nil")
        }
 
+       // OpenAPI上ではnullableだがogenがそれを正しく処理できずにいるため、nullを許容するための特別な処理を追加。
+       if d.Next() == jx.Null {
+               d.Null()
+               return nil
+       }
+
        if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
                switch string(k) {
                default:
```
