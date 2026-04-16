// Copyright 2022-2025 The webaccel-api-go authors
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

package webaccel

// OriginType
const (
	OriginTypesWebServer     = "0" // ウェブサーバ
	OriginTypesObjectStorage = "1" // オブジェクトストレージ
)

// RequestProtocol
const (
	RequestProtocolsHttpAndHttps    = "0" // http/https
	RequestProtocolsHttpsOnly       = "1" // httpsのみ
	RequestProtocolsRedirectToHttps = "2" // httpsに」リダイレクト
)

// OriginProtocol
const (
	OriginProtocolsHttp  = "http"
	OriginProtocolsHttps = "https"
)

// VarySupport
const (
	VarySupportDisabled = "0" // 無効
	VarySupportEnabled  = "1" // 有効
)

// DocIndex
const (
	DocIndexDisabled = "0" // 無効
	DocIndexEnabled  = "1" // 有効
)

// NormalizeAE
const (
	NormalizeAEGz   = "1" // gzipに正規化
	NormalizeAEBrGz = "3" // brとgzipの組に正規化
)

// 各種パラメタの文字列表現。
// Note: APIリクエストそのものには指定できない。
const (
	gunzipCompressionNickname              = "gzip"
	brotliCompressionNickname              = "br+gzip"
	httpOrHttpsRequestProtocolNickname     = "http+https"
	httpsOnlyRequestProtocolNickname       = "https"
	httpsRedirectedRequestProtocolNickname = "https-redirect"
)

var (
	// NormalizeAENicknameStrings
	// NormalizeAEパラメタの文字列表現。APIリクエストには直接指定できないことに注意。
	// MapNormalizeAENicknameToValue を用いてAPIリクエストに指定する値に変換できる。
	NormalizeAENicknameStrings = []string{
		gunzipCompressionNickname,
		brotliCompressionNickname,
	}
	// RequestProtocolStrings
	// RequestProtocolパラメタの文字列表現。APIリクエストには直接指定できないことに注意。
	// MapRequestProtocolNicknameToValue を用いてAPIリクエストに指定する値に変換できる。
	RequestProtocolStrings = []string{
		httpOrHttpsRequestProtocolNickname,
		httpsOnlyRequestProtocolNickname,
		httpsRedirectedRequestProtocolNickname,
	}
	// OriginProtocolStrings
	// OriginProtocolパラメタの文字列表現。いずれかの値をAPIリクエストに直接指定できる。
	OriginProtocolStrings = []string{
		OriginProtocolsHttp,
		OriginProtocolsHttps,
	}
)
