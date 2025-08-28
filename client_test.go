package monitoringsuite

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewClientWithApiUrl_OtherZone(t *testing.T) {
	c, err := NewClientWithApiUrl("https://secure.sakura.ad.jp/cloud/zone/tk1b/api/monitoring/1.0/")
	require.NoError(t, err)
	require.NotNil(t, c)
}
