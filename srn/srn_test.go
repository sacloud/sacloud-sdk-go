// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package srn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSRN(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want SRN
	}{
		{
			name: "numeric number",
			in:   "srnv1:is1y:sakura.iaas.switch:113702228717",
			want: SRN{Version: 1, Location: "is1y", Resource: "sakura.iaas.switch", ID: "113702228717"},
		},
		{
			name: "sqids like number",
			in:   "srnv1:sakura:sakura.cloudhsm.license:5ReQGzN",
			want: SRN{Version: 1, Location: "sakura", Resource: "sakura.cloudhsm.license", ID: "5ReQGzN"},
		},
		{
			name: "uuid number",
			in:   "srnv1:is1b:sakura.apprun.dedicated.auto-scaling-group:019e738e-0915-7a7c-9bd2-ef4284076927",
			want: SRN{Version: 1, Location: "is1b", Resource: "sakura.apprun.dedicated.auto-scaling-group", ID: "019e738e-0915-7a7c-9bd2-ef4284076927"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)

			actual, err := Parse(tt.in)
			assert.NoError(err)
			assert.Equal(tt.want, actual)
			assert.Equal(tt.in, actual.String())
			assert.True(IsSRN(tt.in))
		})
	}
}

func TestParseSRN_Invalid(t *testing.T) {
	tests := []string{
		"",
		"srnv1:is1y:sakura.iaas.switch",
		"srnvx:is1y:sakura.iaas.switch:113702228717",
		"srn1:is1y:sakura.iaas.switch:113702228717",
		"srnv1::sakura.iaas.switch:113702228717",
		"srnv1:is1y::113702228717",
	}

	for _, in := range tests {
		t.Run(in, func(t *testing.T) {
			_, err := Parse(in)
			require.Error(t, err)
			require.False(t, IsSRN(in))
		})
	}
}

func TestSRN_String(t *testing.T) {
	tests := []struct {
		name string
		srn  SRN
		want string
	}{
		{
			name: "numeric id",
			srn: SRN{
				Version:  1,
				Location: "is1y",
				Resource: "sakura.iaas.switch",
				ID:       "113702228717",
			},
			want: "srnv1:is1y:sakura.iaas.switch:113702228717",
		},
		{
			name: "sqids id",
			srn: SRN{
				Version:  1,
				Location: "sakura",
				Resource: "sakura.cloudhsm.license",
				ID:       "5ReQGzN",
			},
			want: "srnv1:sakura:sakura.cloudhsm.license:5ReQGzN",
		},
		{
			name: "uuid id",
			srn: SRN{
				Version:  1,
				Location: "is1b",
				Resource: "sakura.apprun.dedicated.auto-scaling-group",
				ID:       "019e738e-0915-7a7c-9bd2-ef4284076927",
			},
			want: "srnv1:is1b:sakura.apprun.dedicated.auto-scaling-group:019e738e-0915-7a7c-9bd2-ef4284076927",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.srn.String())
		})
	}
}
