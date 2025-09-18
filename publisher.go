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
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type PublisherAPI interface {
	List(ctx context.Context, count *int64, from *int64) ([]v1.Publisher, error)
	Read(ctx context.Context, code string) (*v1.Publisher, error)
}

var _ PublisherAPI = (*publisherOp)(nil)

type publisherOp struct {
	client *v1.Client
}

func NewPublisherOp(client *v1.Client) PublisherAPI {
	return &publisherOp{client: client}
}

func (p *publisherOp) List(ctx context.Context, count *int64, from *int64) ([]v1.Publisher, error) {
	params := v1.PublishersListParams{
		Count: intoOpt[v1.OptInt64](count),
		From:  intoOpt[v1.OptInt64](from),
	}

	result, err := p.client.PublishersList(ctx, params)
	if err != nil {
		return nil, NewAPIError("Publisher.List", http.StatusInternalServerError, errors.Wrap(err, "internal server error"))
	} else {
		return result.Results, err
	}
}

func (p *publisherOp) Read(ctx context.Context, code string) (*v1.Publisher, error) {
	params := v1.PublishersRetrieveParams{Code: code}

	result, err := p.client.PublishersRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusNotFound:
			return nil, NewAPIError("Publisher.Read", e.StatusCode, errors.Wrap(err, "publisher not found"))
		default:
			return nil, NewAPIError("Publisher.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("Publisher.Read", 0, err)
	} else {
		pub := &v1.Publisher{
			Code:        result.GetCode(),
			Description: result.GetDescription(),
			Variants:    result.GetVariants(),
		}
		return pub, nil
	}
}
