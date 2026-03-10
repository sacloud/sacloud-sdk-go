// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"context"

	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type PublisherAPI interface {
	List(ctx context.Context, count *int, from *int) ([]v1.Publisher, error)
	Read(ctx context.Context, code string) (*v1.Publisher, error)
}

var _ PublisherAPI = (*publisherOp)(nil)

type publisherOp struct {
	client *v1.Client
}

func NewPublisherOp(client *v1.Client) PublisherAPI {
	return &publisherOp{client: client}
}

func (p *publisherOp) List(ctx context.Context, count *int, from *int) (ret []v1.Publisher, err error) {
	res, err := errorFromDecodedResponse("Publisher.List", func() (*v1.PaginatedPublisherList, error) {
		return p.client.PublishersList(ctx, v1.PublishersListParams{
			Count: intoOpt[v1.OptInt](count),
			From:  intoOpt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (p *publisherOp) Read(ctx context.Context, code string) (*v1.Publisher, error) {
	res, err := errorFromDecodedResponse("Publisher.Read", func() (*v1.WrappedPublisher, error) {
		return p.client.PublishersRetrieve(ctx, v1.PublishersRetrieveParams{Code: code})
	})
	return unwrapE[*v1.Publisher](res, err)
}
