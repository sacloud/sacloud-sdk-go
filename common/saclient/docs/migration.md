# [`sacloud/api-client-go`](https://github.com/sacloud/api-client-go) からの移行について

過去に `api-client-go` を利用していたプログラムは、以下のようにして `saclient-go` へ移行できます

# 一旦ライブラリを切り替える

新機能などに対応せず、ライブラリを切り替えるのみに注力する場合です。一旦これをやると良いでしょう。これはおもに2つの対応が必要で、

- まずクライアントを構築している部分を変更する必要があり
- 次にリクエスト呼び出し部分を変更する必要がある

## `NewFactory` を移行する

既存コードがこのようになっている場合があります:

```golang
type MyClient struct {
    Profile string
    // ...
    factory *client.Factory
}

func (c *MyClient) init() error {
    var opts []*client.Options
    o, err := client.OptionsFromProfile(c.Profile)
    if err != nil {
        return err
    }
    opts = append(opts, o)
    // ...
    c.factory = client.NewFactory(opts...)
}

func (c *MyClient) Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error) {
    if err := c.init(); err != nil {
        return nil, err
    }

    req, err := c.newRequest(ctx, method, uri, body)
    if err != nil {
        return nil, err
    }

    // API call
    resp, err := c.factory.NewHttpRequestDoer().Do(req)
    if err != nil {
        return nil, err
    }
    // ...
}
```

このようなコードであれば、

- まず構造体で`Factory`のかわりに`saclient.ClientAPI`を持たせる
- 次に`client.NewFactory`のかわりに`saclient.NewFactory`を呼ぶ
- `ClientAPI`で直接`Do`

といった変更で対応可能です。

```diff
--- before
+++ after
@@ -0,0 +0,0 @@
 type MyClient struct {
     Profile string
     // ...
-    factory *client.Factory
+    factory saclient.ClientAPI
 }
 
 func (c *MyClient) init() error {
     var opts []*client.Options
     o, err := client.OptionsFromProfile(c.Profile)
     if err != nil {
         return err
     }
     opts = append(opts, o)
     // ...
-    c.factory = client.NewFactory(opts...)
+    c.factory = saclient.NewFactory(opts...)
 }
 
 func (c *MyClient) Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error) {
     if err := c.init(); err != nil {
         return nil, err
     }
 
     req, err := c.newRequest(ctx, method, uri, body)
     if err != nil {
         return nil, err
     }
 
     // API call
-    resp, err := c.factory.NewHttpRequestDoer().Do(req)
+    resp, err := c.factory.Do(req)
     if err != nil {
         return nil, err
     }
     // ...
 }
```

## `NewClient` を移行する

既存コードがこのようになっている場合があります:

```golang
func NewClientWithApiUrl(apiUrl string, params ...client.ClientParam) (*v1.Client, error) {
    c, err := client.NewClient(apiUrl, params...)
    if err != nil {
        return nil, NewError("NewClientWithApiUrl", err)
    }

    v1Client, err := v1.NewClient(c.ServerURL(), v1.WithClient(c.NewHttpRequestDoer()))
    if err != nil {
        return nil, NewError("NewClientWithApiUrl", err)
    }

    return v1Client, nil
}
```

このようなコードであれば、

- `client.NewClient`のかわりに`saclient.NewClient`を呼ぶ
- `NewHttpRequestDoer`は不要

という変更で対応可能です。

```diff
--- before
+++ after
@@ -0,0 +0,0 @@
 func NewClientWithApiUrl(apiUrl string, params ...client.ClientParam) (*v1.Client, error) {
-    c, err := client.NewClient(apiUrl, params...)
+    c, err := saclient.NewClient(apiUrl, params...)
     if err != nil {
         return nil, NewError("NewClientWithApiUrl", err)
     }
 
-    v1Client, err := v1.NewClient(c.ServerURL(), v1.WithClient(c.NewHttpRequestDoer()))
+    v1Client, err := v1.NewClient(c.ServerURL(), v1.WithClient(c))
     if err != nil {
         return nil, NewError("NewClientWithApiUrl", err)
     }
 
     return v1Client, nil
 }
```

# 新機能に対応する

`saclient-go`では`api-client-go`と比べて新しい機能が追加されていますが、これらに対応するには上記のような書き換えでは不可です。

外部(プロセスの`main`とか)で生成した`saclient.Client`を受け取れるような新規のAPIを作成してください。たとえば上記に加えて、

```diff
--- before
+++ after
@@ -0,0 +0,0 @@
-func NewClientWithApiUrl(apiUrl string, params ...client.ClientParam) (*v1.Client, error) {
-    c, err := saclient.NewClient(apiUrl, params...)
-    if err != nil {
-        return nil, NewError("NewClientWithApiUrl", err)
-    }
+func NewClientWithSaclient(c saclient.ClientAPI) (*v1.Client, error) { 
     v1Client, err := v1.NewClient(c.ServerURL(), v1.WithClient(c))
     if err != nil {
         return nil, NewError("NewClientWithApiUrl", err)
     }
 
     return v1Client, nil
 }
```

使う側は

```golang
func main() {
    var client saclient.Client
    client.FlagSet().Parse(os.Args[1:])
    // ... etc
    v1, err := NewClientWithSaclient(&client)
    if err != nil {
        // ...
    }
}
```
のようになるでしょう。