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
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
)

type CloudHSMAPI interface {
	List(ctx context.Context) ([]v1.CloudHSM, error)
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

func (op *CloudHSMOp) List(ctx context.Context) ([]v1.CloudHSM, error) {
	resp, err := op.client.CloudhsmCloudhsmsList(ctx)
	if err != nil {
		return nil, NewAPIError("CloudHSM.List", 0, err)
	}
	return resp.CloudHSMs, nil
}

type CloudHSMCreateParams struct {
	Name               string
	Description        *string
	Tags               []string
	Ipv4NetworkAddress string
	Ipv4PrefixLength   int
}

func (op *CloudHSMOp) Create(ctx context.Context, p CloudHSMCreateParams) (*v1.CreateCloudHSM, error) {
	if p.Tags == nil {
		p.Tags = []string{}
	}
	resp, err := op.client.CloudhsmCloudhsmsCreate(
		ctx,
		&v1.WrappedCreateCloudHSM{
			CloudHSM: v1.CreateCloudHSM{
				Name:               p.Name,
				Description:        intoOpt[v1.OptString](p.Description),
				Tags:               p.Tags,
				Availability:       v1.AvailabilityEnumAvailable,
				ServiceClass:       v1.ServiceClassEnumCloudCloudhsmPartition,
				Ipv4NetworkAddress: p.Ipv4NetworkAddress,
				Ipv4PrefixLength:   p.Ipv4PrefixLength,
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
	resp, err := op.client.CloudhsmCloudhsmsRetrieve(
		ctx,
		v1.CloudhsmCloudhsmsRetrieveParams{
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
	Ipv4NetworkAddress string
	Ipv4PrefixLength   int
}

func (op *CloudHSMOp) Update(ctx context.Context, id string, p CloudHSMUpdateParams) (*v1.CloudHSM, error) {
	if p.Tags == nil {
		p.Tags = []string{}
	}

	resp, err := op.client.CloudhsmCloudhsmsUpdate(
		ctx,
		&v1.WrappedCloudHSM{
			CloudHSM: v1.CloudHSM{
				ServiceClass:       v1.ServiceClassEnumCloudCloudhsmPartition,
				Availability:       v1.AvailabilityEnumAvailable,
				Name:               p.Name,
				Description:        intoOpt[v1.OptString](p.Description),
				Tags:               p.Tags,
				Ipv4NetworkAddress: p.Ipv4NetworkAddress,
				Ipv4PrefixLength:   p.Ipv4PrefixLength,
			},
		},
		v1.CloudhsmCloudhsmsUpdateParams{
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
	err := op.client.CloudhsmCloudhsmsDestroy(
		ctx,
		v1.CloudhsmCloudhsmsDestroyParams{
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
