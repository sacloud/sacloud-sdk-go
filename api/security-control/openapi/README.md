現状は非公開のOpenAPI定義を利用してクライアントを生成しています。

また、生成されたコードではうまく動かないところがあるので、下記の修正をしています。

## var wrapperDotRetryAfterValの型の修正

- oas_schemas_gen.go

`TooManyRequestsRetryAfter` を定義。

```
// TooManyRequestsRetryAfter represents sum type.
type TooManyRequestsRetryAfter struct {
	Type                TooManyRequestsRetryAfterType // switch on this field
	ClientDelaySeconds  ClientDelaySeconds
	ClientRetryDateTime ClientRetryDateTime
}

// TooManyRequestsRetryAfterType is oneOf type of TooManyRequestsRetryAfter.
type TooManyRequestsRetryAfterType string

// Possible values for TooManyRequestsRetryAfterType.
const (
	ClientDelaySecondsTooManyRequestsRetryAfter  TooManyRequestsRetryAfterType = "ClientDelaySeconds"
	ClientRetryDateTimeTooManyRequestsRetryAfter TooManyRequestsRetryAfterType = "ClientRetryDateTime"
)

// IsClientDelaySeconds reports whether TooManyRequestsRetryAfter is ClientDelaySeconds.
func (s TooManyRequestsRetryAfter) IsClientDelaySeconds() bool {
	return s.Type == ClientDelaySecondsTooManyRequestsRetryAfter
}

// IsClientRetryDateTime reports whether TooManyRequestsRetryAfter is ClientRetryDateTime.
func (s TooManyRequestsRetryAfter) IsClientRetryDateTime() bool {
	return s.Type == ClientRetryDateTimeTooManyRequestsRetryAfter
}

// SetClientDelaySeconds sets TooManyRequestsRetryAfter to ClientDelaySeconds.
func (s *TooManyRequestsRetryAfter) SetClientDelaySeconds(v ClientDelaySeconds) {
	s.Type = ClientDelaySecondsTooManyRequestsRetryAfter
	s.ClientDelaySeconds = v
}

// GetClientDelaySeconds returns ClientDelaySeconds and true boolean if TooManyRequestsRetryAfter is ClientDelaySeconds.
func (s TooManyRequestsRetryAfter) GetClientDelaySeconds() (v ClientDelaySeconds, ok bool) {
	if !s.IsClientDelaySeconds() {
		return v, false
	}
	return s.ClientDelaySeconds, true
}

// NewClientDelaySecondsTooManyRequestsRetryAfter returns new TooManyRequestsRetryAfter from ClientDelaySeconds.
func NewClientDelaySecondsTooManyRequestsRetryAfter(v ClientDelaySeconds) TooManyRequestsRetryAfter {
	var s TooManyRequestsRetryAfter
	s.SetClientDelaySeconds(v)
	return s
}

// SetClientRetryDateTime sets ProjectActivationReadTooManyRequestsRetryAfter to ClientRetryDateTime.
func (s *TooManyRequestsRetryAfter) SetClientRetryDateTime(v ClientRetryDateTime) {
	s.Type = ClientRetryDateTimeTooManyRequestsRetryAfter
	s.ClientRetryDateTime = v
}

// GetClientRetryDateTime returns ClientRetryDateTime and true boolean if ProjectActivationReadTooManyRequestsRetryAfter is ClientRetryDateTime.
func (s TooManyRequestsRetryAfter) GetClientRetryDateTime() (v ClientRetryDateTime, ok bool) {
	if !s.IsClientRetryDateTime() {
		return v, false
	}
	return s.ClientRetryDateTime, true
}

// NewClientRetryDateTimeTooManyRequestsRetryAfter returns new TooManyRequestsRetryAfter from ClientRetryDateTime.
func NewClientRetryDateTimeTooManyRequestsRetryAfter(v ClientRetryDateTime) TooManyRequestsRetryAfter {
	var s TooManyRequestsRetryAfter
	s.SetClientRetryDateTime(v)
	return s
}

// NewOptTooManyRequestsRetryAfter returns new OptTooManyRequestsRetryAfter with value set to v.
func NewOptTooManyRequestsRetryAfter(v TooManyRequestsRetryAfter) OptTooManyRequestsRetryAfter {
	return OptTooManyRequestsRetryAfter{
		Value: v,
		Set:   true,
	}
}

// OptTooManyRequestsRetryAfter is optional TooManyRequestsRetryAfter.
type OptTooManyRequestsRetryAfter struct {
	Value TooManyRequestsRetryAfter
	Set   bool
}

// IsSet returns true if OptTooManyRequestsRetryAfter was set.
func (o OptTooManyRequestsRetryAfter) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptTooManyRequestsRetryAfter) Reset() {
	var v TooManyRequestsRetryAfter
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptTooManyRequestsRetryAfter) SetTo(v TooManyRequestsRetryAfter) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptTooManyRequestsRetryAfter) Get() (v TooManyRequestsRetryAfter, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptTooManyRequestsRetryAfter) Or(d TooManyRequestsRetryAfter) TooManyRequestsRetryAfter {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}
```

`TooManyRequestsHeaders`の定義を変更

```
diff --git a/apis/v1/oas_schemas_gen.go b/apis/v1/oas_schemas_gen.go
index 562dd94..cdd5665 100644
--- a/apis/v1/oas_schemas_gen.go
+++ b/apis/v1/oas_schemas_gen.go
@@ -5129,12 +5129,12 @@ func (s *TooManyRequestsDetail) UnmarshalText(data []byte) error {
 
 // TooManyRequestsHeaders wraps TooManyRequests with response headers.
 type TooManyRequestsHeaders struct {
-       RetryAfter OptProjectActivationReadTooManyRequestsRetryAfter
+       RetryAfter OptTooManyRequestsRetryAfter
        Response   TooManyRequests
 }
 
 // GetRetryAfter returns the value of RetryAfter.
-func (s *TooManyRequestsHeaders) GetRetryAfter() OptProjectActivationReadTooManyRequestsRetryAfter {
+func (s *TooManyRequestsHeaders) GetRetryAfter() OptTooManyRequestsRetryAfter {
        return s.RetryAfter
 }
 
@@ -5144,7 +5144,7 @@ func (s *TooManyRequestsHeaders) GetResponse() TooManyRequests {
 }
 
 // SetRetryAfter sets the value of RetryAfter.
-func (s *TooManyRequestsHeaders) SetRetryAfter(val OptProjectActivationReadTooManyRequestsRetryAfter) {
+func (s *TooManyRequestsHeaders) SetRetryAfter(val OptTooManyRequestsRetryAfter) {
        s.RetryAfter = val
 }

```

- oas_response_decoders_gen.go

`var wrapperDotRetryAfterVal AutomatedActionsCreateTooManyRequestsRetryAfter` 等の定義を `var wrapperDotRetryAfterVal TooManyRequestsRetryAfter` に変更。

```
$ sed -i '' 's/var wrapperDotRetryAfterVal .*/var wrapperDotRetryAfterVal TooManyRequestsRetryAfter/' apis/v1/oas_response_decoders_gen.go
```

- apis以下の全体のdiff

```
diff -u apis/v2/oas_response_decoders_gen.go apis/v1/oas_response_decoders_gen.go
--- apis/v2/oas_response_decoders_gen.go	2026-01-27 17:01:51
+++ apis/v1/oas_response_decoders_gen.go	2026-01-27 16:58:01
@@ -348,7 +348,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsCreateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -478,7 +478,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsCreateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -853,7 +853,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsDeleteTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -983,7 +983,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsDeleteTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -1831,7 +1831,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -1961,7 +1961,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -2533,7 +2533,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -2663,7 +2663,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal AutomatedActionsUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -3235,7 +3235,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesListTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -3365,7 +3365,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesListTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -3859,7 +3859,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -3989,7 +3989,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -4561,7 +4561,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -4691,7 +4691,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal EvaluationRulesUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -5254,7 +5254,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationCreateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -5384,7 +5384,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationCreateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -5869,7 +5869,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -5999,7 +5999,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationReadTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -6562,7 +6562,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
@@ -6692,7 +6692,7 @@
 				if err := func() error {
 					if err := h.HasParam(cfg); err == nil {
 						if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
-							var wrapperDotRetryAfterVal ProjectActivationUpdateTooManyRequestsRetryAfter
+							var wrapperDotRetryAfterVal TooManyRequestsRetryAfter
 							if err := func() error {
 								var failures []error
 								// Try to decode as ClientDelaySeconds
diff -u apis/v2/oas_schemas_gen.go apis/v1/oas_schemas_gen.go
--- apis/v2/oas_schemas_gen.go	2026-01-27 17:01:51
+++ apis/v1/oas_schemas_gen.go	2026-01-27 16:58:01
@@ -5129,12 +5129,12 @@
 
 // TooManyRequestsHeaders wraps TooManyRequests with response headers.
 type TooManyRequestsHeaders struct {
-	RetryAfter OptProjectActivationReadTooManyRequestsRetryAfter
+	RetryAfter OptTooManyRequestsRetryAfter
 	Response   TooManyRequests
 }
 
 // GetRetryAfter returns the value of RetryAfter.
-func (s *TooManyRequestsHeaders) GetRetryAfter() OptProjectActivationReadTooManyRequestsRetryAfter {
+func (s *TooManyRequestsHeaders) GetRetryAfter() OptTooManyRequestsRetryAfter {
 	return s.RetryAfter
 }
 
@@ -5144,7 +5144,7 @@
 }
 
 // SetRetryAfter sets the value of RetryAfter.
-func (s *TooManyRequestsHeaders) SetRetryAfter(val OptProjectActivationReadTooManyRequestsRetryAfter) {
+func (s *TooManyRequestsHeaders) SetRetryAfter(val OptTooManyRequestsRetryAfter) {
 	s.RetryAfter = val
 }
 
@@ -5385,3 +5385,117 @@
 func (*UnexpectedErrorStatusCode) projectActivationCreateRes() {}
 func (*UnexpectedErrorStatusCode) projectActivationReadRes()   {}
 func (*UnexpectedErrorStatusCode) projectActivationUpdateRes() {}
+
+// TooManyRequestsRetryAfter represents sum type.
+type TooManyRequestsRetryAfter struct {
+	Type                TooManyRequestsRetryAfterType // switch on this field
+	ClientDelaySeconds  ClientDelaySeconds
+	ClientRetryDateTime ClientRetryDateTime
+}
+
+// TooManyRequestsRetryAfterType is oneOf type of TooManyRequestsRetryAfter.
+type TooManyRequestsRetryAfterType string
+
+// Possible values for TooManyRequestsRetryAfterType.
+const (
+	ClientDelaySecondsTooManyRequestsRetryAfter  TooManyRequestsRetryAfterType = "ClientDelaySeconds"
+	ClientRetryDateTimeTooManyRequestsRetryAfter TooManyRequestsRetryAfterType = "ClientRetryDateTime"
+)
+
+// IsClientDelaySeconds reports whether TooManyRequestsRetryAfter is ClientDelaySeconds.
+func (s TooManyRequestsRetryAfter) IsClientDelaySeconds() bool {
+	return s.Type == ClientDelaySecondsTooManyRequestsRetryAfter
+}
+
+// IsClientRetryDateTime reports whether TooManyRequestsRetryAfter is ClientRetryDateTime.
+func (s TooManyRequestsRetryAfter) IsClientRetryDateTime() bool {
+	return s.Type == ClientRetryDateTimeTooManyRequestsRetryAfter
+}
+
+// SetClientDelaySeconds sets TooManyRequestsRetryAfter to ClientDelaySeconds.
+func (s *TooManyRequestsRetryAfter) SetClientDelaySeconds(v ClientDelaySeconds) {
+	s.Type = ClientDelaySecondsTooManyRequestsRetryAfter
+	s.ClientDelaySeconds = v
+}
+
+// GetClientDelaySeconds returns ClientDelaySeconds and true boolean if TooManyRequestsRetryAfter is ClientDelaySeconds.
+func (s TooManyRequestsRetryAfter) GetClientDelaySeconds() (v ClientDelaySeconds, ok bool) {
+	if !s.IsClientDelaySeconds() {
+		return v, false
+	}
+	return s.ClientDelaySeconds, true
+}
+
+// NewClientDelaySecondsTooManyRequestsRetryAfter returns new TooManyRequestsRetryAfter from ClientDelaySeconds.
+func NewClientDelaySecondsTooManyRequestsRetryAfter(v ClientDelaySeconds) TooManyRequestsRetryAfter {
+	var s TooManyRequestsRetryAfter
+	s.SetClientDelaySeconds(v)
+	return s
+}
+
+// SetClientRetryDateTime sets ProjectActivationReadTooManyRequestsRetryAfter to ClientRetryDateTime.
+func (s *TooManyRequestsRetryAfter) SetClientRetryDateTime(v ClientRetryDateTime) {
+	s.Type = ClientRetryDateTimeTooManyRequestsRetryAfter
+	s.ClientRetryDateTime = v
+}
+
+// GetClientRetryDateTime returns ClientRetryDateTime and true boolean if ProjectActivationReadTooManyRequestsRetryAfter is ClientRetryDateTime.
+func (s TooManyRequestsRetryAfter) GetClientRetryDateTime() (v ClientRetryDateTime, ok bool) {
+	if !s.IsClientRetryDateTime() {
+		return v, false
+	}
+	return s.ClientRetryDateTime, true
+}
+
+// NewClientRetryDateTimeTooManyRequestsRetryAfter returns new TooManyRequestsRetryAfter from ClientRetryDateTime.
+func NewClientRetryDateTimeTooManyRequestsRetryAfter(v ClientRetryDateTime) TooManyRequestsRetryAfter {
+	var s TooManyRequestsRetryAfter
+	s.SetClientRetryDateTime(v)
+	return s
+}
+
+// NewOptTooManyRequestsRetryAfter returns new OptTooManyRequestsRetryAfter with value set to v.
+func NewOptTooManyRequestsRetryAfter(v TooManyRequestsRetryAfter) OptTooManyRequestsRetryAfter {
+	return OptTooManyRequestsRetryAfter{
+		Value: v,
+		Set:   true,
+	}
+}
+
+// OptTooManyRequestsRetryAfter is optional TooManyRequestsRetryAfter.
+type OptTooManyRequestsRetryAfter struct {
+	Value TooManyRequestsRetryAfter
+	Set   bool
+}
+
+// IsSet returns true if OptTooManyRequestsRetryAfter was set.
+func (o OptTooManyRequestsRetryAfter) IsSet() bool { return o.Set }
+
+// Reset unsets value.
+func (o *OptTooManyRequestsRetryAfter) Reset() {
+	var v TooManyRequestsRetryAfter
+	o.Value = v
+	o.Set = false
+}
+
+// SetTo sets value to v.
+func (o *OptTooManyRequestsRetryAfter) SetTo(v TooManyRequestsRetryAfter) {
+	o.Set = true
+	o.Value = v
+}
+
+// Get returns value and boolean that denotes whether value was set.
+func (o OptTooManyRequestsRetryAfter) Get() (v TooManyRequestsRetryAfter, ok bool) {
+	if !o.Set {
+		return v, false
+	}
+	return o.Value, true
+}
+
+// Or returns value if set, or given parameter if does not.
+func (o OptTooManyRequestsRetryAfter) Or(d TooManyRequestsRetryAfter) TooManyRequestsRetryAfter {
+	if v, ok := o.Get(); ok {
+		return v
+	}
+	return d
+}

```