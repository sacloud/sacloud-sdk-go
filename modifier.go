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
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
)

const (
	commonServiceItemPath      = "commonserviceitem"
	commonServiceItemKey       = "CommonServiceItem"
	commonServiceItemListKey   = "CommonServiceItems"
	commonServiceItemIconKey   = "Icon"
	commonServiceItemIconIDKey = "ID"
)

func modifiyMiddleware() saclient.Middleware {
	return func(req *http.Request, pull func() (saclient.Middleware, bool)) (*http.Response, error) {
		if err := requestModifier(req); err != nil {
			return nil, err
		}
		cont, ok := pull()
		if !ok {
			return nil, errors.New("middleware not found error")
		}
		resp, err := cont(req, pull)
		if err != nil {
			return nil, err
		}
		if err := responseModifier(req, resp); err != nil {
			return nil, err
		}
		return resp, nil
	}
}
func requestModifier(req *http.Request) error {
	listpath := strings.TrimSuffix(req.URL.Path, "/")

	// commonserviceitem list API , Provider.Class query param setting
	if path.Base(listpath) == commonServiceItemPath &&
		req.Method == http.MethodGet {
		providerTarget, err := getContextProviderClass(req.Context())
		if err != nil {
			return err
		}
		if err := setJSONOnlyQuery(req, providerTarget); err != nil {
			return err
		}
	}
	return nil
}

func responseModifier(req *http.Request, resp *http.Response) error {
	// commonserviceitem　list and get check.
	if !strings.Contains(req.URL.String(), commonServiceItemPath) {
		return nil
	}
	body := resp.Body
	defer func() {
		_ = body.Close()
	}()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	// JSON decode
	var data map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		// Unmarshal errro, restore body and return
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return nil
	}
	// replace null to ""
	converted := replaceIconNullWithCommonServiceItem(data)
	// JSON encode
	newBody, err := json.Marshal(converted)
	if err != nil {
		return err
	}
	// replace body
	resp.Body = io.NopCloser(bytes.NewReader(newBody))
	resp.ContentLength = int64(len(newBody))
	resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
	return nil
}

type filterQuery struct {
	Filter map[string]string `json:"Filter"`
}

func setJSONOnlyQuery(req *http.Request, providerClass v1.CommonServiceItemProviderClass) error {
	q := filterQuery{
		Filter: map[string]string{
			"Provider.Class": string(providerClass),
		},
	}
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	req.URL.RawQuery = string(b)
	return nil
}

func replaceIconNullWithCommonServiceItem(items map[string]interface{}) map[string]interface{} {
	//  if Icon value is null , replace with empty object
	replaceIcon := func(data map[string]interface{}) {
		if v, ok := data[commonServiceItemIconKey]; ok && v == nil {
			data[commonServiceItemIconKey] = map[string]interface{}{}
			data[commonServiceItemIconKey].(map[string]interface{})[commonServiceItemIconIDKey] = ""
		}
	}

	// case : List
	if itemsList, ok := items[commonServiceItemListKey].([]interface{}); ok {
		for i, item := range itemsList {
			if data, ok := item.(map[string]interface{}); ok {
				replaceIcon(data)
				itemsList[i] = data
			}
		}
		items[commonServiceItemListKey] = itemsList
		return items
	}

	// case : default
	if data, ok := items[commonServiceItemKey].(map[string]interface{}); ok {
		replaceIcon(data)
		items[commonServiceItemKey] = data
		return items
	}
	return items
}
