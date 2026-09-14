// Copyright 2025- The sacloud/cloudhsm-api-go authors
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

package cloudhsm

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/cloudhsm/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type CloudHSMAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.CloudHSM, error)
	Create(ctx context.Context, request CloudHSMCreateParams) (*v1.CreateCloudHSM, error)
	Read(ctx context.Context, id string) (*v1.CloudHSM, error)
	Update(ctx context.Context, id string, params CloudHSMUpdateParams) (*v1.CloudHSM, error)
	Delete(ctx context.Context, id string) error
}

var _ CloudHSMAPI = (*CloudHSMOp)(nil)

type CloudHSMOp struct {
	client *v1.Client
}

func NewCloudHSMOp(client *v1.Client) CloudHSMAPI {
	return &CloudHSMOp{client: client}
}

func (op *CloudHSMOp) List(ctx context.Context, count, from *int) ([]v1.CloudHSM, error) {
	resp, err := op.client.ListCloudHSMs(
		ctx,
		v1.ListCloudHSMsParams{
			Count: into.Opt[v1.OptInt](count),
			From:  into.Opt[v1.OptInt](from),
		},
	)
	if err != nil {
		return nil, NewAPIError("CloudHSM.List", 0, err)
	}
	return resp.GetCloudHSMs(), nil
}

type CloudHSMCreateParams struct {
	Name               string
	Description        *string
	Tags               []string
	IPv4NetworkAddress string
	IPv4PrefixLength   int
}

func (op *CloudHSMOp) Create(ctx context.Context, p CloudHSMCreateParams) (*v1.CreateCloudHSM, error) {
	resp, err := op.client.CreateCloudHSM(
		ctx,
		&v1.WrappedCreateCloudHSMRequest{
			CloudHSM: v1.CreateCloudHSMRequest{
				Name:               p.Name,
				Description:        into.Opt[v1.OptString](p.Description),
				Tags:               into.OptNilArray[v1.OptNilStringArray](&p.Tags),
				IPv4NetworkAddress: p.IPv4NetworkAddress,
				IPv4PrefixLength:   p.IPv4PrefixLength,
			},
		},
	)

	if err == nil {
		ret := resp.GetCloudHSM()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSM.Create", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("CloudHSM.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("CloudHSM.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *CloudHSMOp) Read(ctx context.Context, id string) (*v1.CloudHSM, error) {
	resp, err := op.client.ReadCloudHSM(
		ctx,
		v1.ReadCloudHSMParams{
			ResourceID: id,
		},
	)

	if err == nil {
		ret := resp.GetCloudHSM()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSM.Read", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("CloudHSM.Read", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return nil, NewAPIError("CloudHSM.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

type CloudHSMUpdateParams struct {
	Name               string
	Description        *string
	Tags               []string
	IPv4NetworkAddress string
	IPv4PrefixLength   int
}

func (op *CloudHSMOp) Update(ctx context.Context, id string, p CloudHSMUpdateParams) (*v1.CloudHSM, error) {
	resp, err := op.client.UpdateCloudHSM(
		ctx,
		&v1.WrappedCloudHSMRequest{
			CloudHSM: v1.CloudHSMRequest{
				Name:               p.Name,
				Description:        into.Opt[v1.OptString](p.Description),
				Tags:               into.OptNilArray[v1.OptNilStringArray](&p.Tags),
				IPv4NetworkAddress: p.IPv4NetworkAddress,
				IPv4PrefixLength:   p.IPv4PrefixLength,
			},
		},
		v1.UpdateCloudHSMParams{
			ResourceID: id,
		},
	)

	if err == nil {
		ret := resp.GetCloudHSM()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("CloudHSM.Update", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("CloudHSM.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("CloudHSM.Update", 0, err)
	}
}

func (op *CloudHSMOp) Delete(ctx context.Context, id string) error {
	err := op.client.DeleteCloudHSM(
		ctx,
		v1.DeleteCloudHSMParams{
			ResourceID: id,
		},
	)

	if err == nil {
		return nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return NewAPIError("CloudHSM.Delete", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return NewAPIError("CloudHSM.Delete", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return NewAPIError("CloudHSM.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}
