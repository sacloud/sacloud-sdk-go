AppRun共用型のOpenAPI定義は以下のページで公開されています。

https://manual.sakura.ad.jp/api/cloud/portal/?api=apprun-shared-api

apprun-api-goではogenで生成したコードそのままではエラーが発生するので、そのために以下の変更を生成したコードに施しています。

```
diff --git b/apis/v1/oas_json_gen.go a/apis/v1/oas_json_gen.go
index 87705e8..5990442 100644
--- b/apis/v1/oas_json_gen.go
+++ a/apis/v1/oas_json_gen.go
@@ -9485,6 +9485,12 @@ func (s *HandlerListTrafficsMeta) Decode(d *jx.Decoder) error {
 		return errors.New("invalid: unable to decode HandlerListTrafficsMeta to nil")
 	}
 
+	// OpenAPI上ではnullableだがogenがそれを正しく処理できずにいるため、nullを許容するための特別な処理を追加。
+	if d.Next() == jx.Null {
+		d.Null()
+		return nil
+	}
+
 	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
 		switch string(k) {
 		default:
@@ -12886,6 +12892,12 @@ func (s *HandlerPutTrafficsMeta) Decode(d *jx.Decoder) error {
 		return errors.New("invalid: unable to decode HandlerPutTrafficsMeta to nil")
 	}
 
+	// OpenAPI上ではnullableだがogenがそれを正しく処理できずにいるため、nullを許容するための特別な処理を追加。
+	if d.Next() == jx.Null {
+		d.Null()
+		return nil
+	}
+
 	if err := d.ObjBytes(func(d *jx.Decoder, k []byte) error {
 		switch string(k) {
 		default:
```
