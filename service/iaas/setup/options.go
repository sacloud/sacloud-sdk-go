// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package setup

import "time"

var (
	//	DefaultNICUpdateWaitDuration デフォルトのNIC更新待ち
	DefaultNICUpdateWaitDuration = 20 * time.Second

	// DefaultMaxRetryCount デフォルトリトライ最大数
	DefaultMaxRetryCount = 3

	// DefaultProvisioningRetryCount リソースごとのプロビジョニングAPI呼び出しのリトライ最大数
	DefaultProvisioningRetryCount = 10

	// DefaultProvisioningWaitInterval リソースごとのプロビジョニングAPI呼び出しのリトライ間隔
	DefaultProvisioningWaitInterval = 5 * time.Second

	// DefaultDeleteRetryCount リソースごとの削除API呼び出しのリトライ最大数
	DefaultDeleteRetryCount = 10

	// DefaultDeleteWaitInterval リソースごとの削除API呼び出しのリトライ間隔
	DefaultDeleteWaitInterval = 10 * time.Second

	// DefaultPollingInterval ポーリング処理の間隔
	DefaultPollingInterval = 5 * time.Second
)

// Options アプライアンス作成時に利用するsetup.RetryableSetupのパラメータ
type Options struct {
	// BootAfterBuild Buildの後に再起動を行うか
	BootAfterBuild bool
	// NICUpdateWaitDuration NIC接続切断操作の後の待ち時間
	NICUpdateWaitDuration time.Duration
	// RetryCount リトライ回数
	RetryCount int
	// ProvisioningRetryCount リトライ回数
	ProvisioningRetryCount int
	// ProvisioningRetryInterval
	ProvisioningRetryInterval time.Duration
	// DeleteRetryCount 削除リトライ回数
	DeleteRetryCount int
	// DeleteRetryInterval 削除リトライ間隔
	DeleteRetryInterval time.Duration
	// sacloud.StateWaiterによるステート待ちの間隔
	PollingInterval time.Duration
}

func (o *Options) Init() {
	if o.NICUpdateWaitDuration == time.Duration(0) {
		o.NICUpdateWaitDuration = DefaultNICUpdateWaitDuration
	}
	if o.RetryCount <= 0 {
		o.RetryCount = DefaultMaxRetryCount
	}
	if o.DeleteRetryCount <= 0 {
		o.DeleteRetryCount = DefaultDeleteRetryCount
	}
	if o.DeleteRetryInterval <= 0 {
		o.DeleteRetryInterval = DefaultDeleteWaitInterval
	}
	if o.ProvisioningRetryCount <= 0 {
		o.ProvisioningRetryCount = DefaultProvisioningRetryCount
	}
	if o.ProvisioningRetryInterval <= 0 {
		o.ProvisioningRetryInterval = DefaultProvisioningWaitInterval
	}
	if o.PollingInterval <= 0 {
		o.PollingInterval = DefaultPollingInterval
	}
}
