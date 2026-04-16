// Copyright 2026- The sacloud/simple-notification-api-go Authors
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

package simplenotification

import (
	"context"
	"fmt"

	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
)

// contextKey is a custom type for context keys in this package
type contextKey string

const providerClassKey contextKey = "Provider.Class"

func setContextProviderClass(ctx context.Context, providerClass v1.CommonServiceItemProviderClass) context.Context {
	return context.WithValue(ctx, providerClassKey, providerClass)
}

func getContextProviderClass(ctx context.Context) (v1.CommonServiceItemProviderClass, error) {
	v := ctx.Value(providerClassKey)
	if v == nil {
		return "", fmt.Errorf("Provider.Class not found in context")
	}
	s, ok := v.(v1.CommonServiceItemProviderClass)
	if !ok {
		return "", fmt.Errorf("Provider.Class is not a string")
	}
	return s, nil
}
