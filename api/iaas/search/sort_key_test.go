// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortKey_MarshalJSON(t *testing.T) {
	testCases := []struct {
		input  SortKey
		expect string
	}{
		{input: SortKeyAsc("Name"), expect: `"Name"`},
		{input: SortKeyDesc("Name"), expect: `"-Name"`},
		{input: SortKey{}, expect: `""`},
	}
	for _, tc := range testCases {
		v, err := json.Marshal(tc.input)
		require.NoError(t, err)
		require.Equal(t, tc.expect, string(v))
	}
}

func TestSortKey_UnmarshalJSON(t *testing.T) {
	t.Run("string form", func(t *testing.T) {
		testCases := []struct {
			input  string
			expect SortKey
		}{
			{input: `"Name"`, expect: SortKeyAsc("Name")},
			{input: `"-Name"`, expect: SortKeyDesc("Name")},
			{input: `""`, expect: SortKey{}},
		}
		for _, tc := range testCases {
			var k SortKey
			err := json.Unmarshal([]byte(tc.input), &k)
			require.NoError(t, err)
			require.Equal(t, tc.expect, k)
		}
	})

	t.Run("object form", func(t *testing.T) {
		testCases := []struct {
			input  string
			expect SortKey
		}{
			{input: `{"Key":"Name","Order":0}`, expect: SortKeyAsc("Name")},
			{input: `{"Key":"Name","Order":1}`, expect: SortKeyDesc("Name")},
			{input: `{"Key":"Name"}`, expect: SortKeyAsc("Name")},
		}
		for _, tc := range testCases {
			var k SortKey
			err := json.Unmarshal([]byte(tc.input), &k)
			require.NoError(t, err)
			require.Equal(t, tc.expect, k)
		}
	})

	t.Run("null", func(t *testing.T) {
		// null は値を変更しない
		k := SortKeyAsc("Name")
		err := json.Unmarshal([]byte("null"), &k)
		require.NoError(t, err)
		require.Equal(t, SortKeyAsc("Name"), k)
	})

	t.Run("invalid", func(t *testing.T) {
		for _, data := range []string{`1`, `true`, `[]`} {
			var k SortKey
			err := json.Unmarshal([]byte(data), &k)
			require.Error(t, err)
			require.ErrorContains(t, err, "SortKey")
		}
	})
}

func TestSortKeys_UnmarshalJSON(t *testing.T) {
	// skr の --request 経由でのデコードを想定した、FindCondition と同様の構造へのデコード
	var v struct {
		Sort SortKeys `json:"Sort"`
	}
	err := json.Unmarshal([]byte(`{"Sort":["Name","-CreatedAt"]}`), &v)
	require.NoError(t, err)
	require.Equal(t, SortKeys{SortKeyAsc("Name"), SortKeyDesc("CreatedAt")}, v.Sort)

	// Marshal結果とのラウンドトリップ
	v2, err := json.Marshal(v.Sort)
	require.NoError(t, err)
	require.Equal(t, `["Name","-CreatedAt"]`, string(v2))

	var got SortKeys
	err = json.Unmarshal(v2, &got)
	require.NoError(t, err)
	require.Equal(t, v.Sort, got)
}
