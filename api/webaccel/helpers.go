package webaccel

import (
	"fmt"
)

// MapRequestProtocolNicknameToValue maps non-origin request protocol nickname
// (e.g. http+https) to a string const which is acceptable as a parameter of Site.
func MapRequestProtocolNicknameToValue(protocolNickname string) (string, error) {
	switch protocolNickname {
	case httpOrHttpsRequestProtocolNickname:
		return RequestProtocolsHttpAndHttps, nil
	case httpsOnlyRequestProtocolNickname:
		return RequestProtocolsHttpsOnly, nil
	case httpsRedirectedRequestProtocolNickname:
		return RequestProtocolsRedirectToHttps, nil
	default:
		return "", fmt.Errorf("invalid request protocol: %s", protocolNickname)
	}
}

// MapRequestProtocolToNickname maps the site's (non-origin) request protocol value
// (e.g. `1`) into human-readable string format such as `https`.
func MapRequestProtocolToNickname(site *Site) (string, error) {
	switch site.RequestProtocol {
	case RequestProtocolsHttpAndHttps:
		return httpOrHttpsRequestProtocolNickname, nil
	case RequestProtocolsHttpsOnly:
		return httpsOnlyRequestProtocolNickname, nil
	case RequestProtocolsRedirectToHttps:
		return httpsRedirectedRequestProtocolNickname, nil
	default:
		return "", fmt.Errorf("invalid request protocol: %s", site.RequestProtocol)
	}
}

// MapNormalizeAENicknameToValue maps compression type nickname (e.g. br+gzip) to
// string const which is acceptable as a parameter of Site.
func MapNormalizeAENicknameToValue(compressionNickname string) (string, error) {
	switch compressionNickname {
	case gunzipCompressionNickname:
		return NormalizeAEGz, nil
	case brotliCompressionNickname:
		return NormalizeAEBrGz, nil
	}
	return "", fmt.Errorf("invalid normalize_ae parameter: '%s'", compressionNickname)
}

// MapNormalizeAEValueToNickname maps the site's accept-encoding normalization parameter value
// into human-readable  nickname such as `gzip`.
func MapNormalizeAEValueToNickname(site *Site) (string, error) {
	if site.NormalizeAE != "" {
		if site.NormalizeAE == NormalizeAEBrGz {
			return brotliCompressionNickname, nil
		} else if site.NormalizeAE == NormalizeAEGz {
			return gunzipCompressionNickname, nil
		}
		return "", fmt.Errorf("invalid normalize_ae parameter: '%s'", site.NormalizeAE)
	}
	//NOTE: APIが返却するデフォルト値は""。
	// このフィールドでで "gzip" と "" が持つ効果は同一であるため、
	// "gzip" として正規化する
	return gunzipCompressionNickname, nil
}
