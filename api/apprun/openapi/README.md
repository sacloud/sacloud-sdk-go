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
diff --git b/apis/v1/oas_response_decoders_gen.go a/apis/v1/oas_response_decoders_gen.go
index aba03d1..02d6e3a 100644
--- b/apis/v1/oas_response_decoders_gen.go
+++ a/apis/v1/oas_response_decoders_gen.go
@@ -3791,37 +3791,14 @@ func decodePostUserResponse(resp *http.Response) (res PostUserRes, _ error) {
 		}
 		switch {
 		case ct == "application/json":
-			buf, err := io.ReadAll(resp.Body)
-			if err != nil {
-				return res, err
-			}
-			d := jx.DecodeBytes(buf)
-
+			// OpenAPIとは違ってレスポンスが空なので、ここではデコードせずに直接レスポンスを返す。OpenAPIが修正されたら削除
 			var response PostUserConflict
-			if err := func() error {
-				if err := response.Decode(d); err != nil {
-					return err
-				}
-				if err := d.Skip(); err != io.EOF {
-					return errors.New("unexpected trailing data")
-				}
-				return nil
-			}(); err != nil {
-				err = &ogenerrors.DecodeBodyError{
-					ContentType: ct,
-					Body:        buf,
-					Err:         err,
-				}
-				return res, err
-			}
-			// Validate response.
-			if err := func() error {
-				if err := response.Validate(); err != nil {
-					return err
-				}
-				return nil
-			}(); err != nil {
-				return res, errors.Wrap(err, "validate")
+			response.Type = PostUserConflictType("Conflict")
+			response.ModelDefaultError = ModelDefaultError{
+				Error: ModelDefaultErrorError{
+					Code:    409,
+					Message: "Conflict",
+				},
 			}
 			return &response, nil
 		default:

```

また、OpenAPI上ではrequiredになっているものの、実際のレスポンスではOptionalなフィールドに関してはOpenAPI定義を修正しています。

```
diff --git a/openapi/openapi.yaml b/openapi/openapi.yaml
index 7fcd3a4..5e939ca 100644
--- a/openapi/openapi.yaml
+++ b/openapi/openapi.yaml
@@ -2484,7 +2484,7 @@ components:
           - reason
           - message
           - location_type
-          - location
+          # - location
         properties:
           domain:
```