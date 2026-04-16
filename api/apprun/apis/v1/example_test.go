// Copyright 2021-2024 The sacloud/apprun-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1_test

import (
	"context"
	"fmt"
	"io"
	"os"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

var serverURL = "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api"

func Example_getUser() {
	token := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	secret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	client, err := v1.NewClientWithResponses(serverURL, func(c *v1.Client) error {
		c.RequestEditors = []v1.RequestEditorFn{
			v1.AppRunAuthInterceptor(token, secret),
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.GetUser(context.Background())
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
}
