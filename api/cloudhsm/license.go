// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

type LicenseAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.CloudHSMSoftwareLicense, error)
	Create(ctx context.Context, request CloudHSMSoftwareLicenseCreateParams) (*v1.CreateCloudHSMSoftwareLicense, error)
	Read(ctx context.Context, id string) (*v1.CloudHSMSoftwareLicense, error)
	Update(ctx context.Context, id string, params CloudHSMSoftwareLicenseUpdateParams) (*v1.CloudHSMSoftwareLicense, error)
	Delete(ctx context.Context, id string) error
}

var _ LicenseAPI = (*LicenseOp)(nil)

type LicenseOp struct {
	client *v1.Client
}

func NewLicenseOp(client *v1.Client) LicenseAPI {
	return &LicenseOp{client: client}
}

func (op *LicenseOp) List(ctx context.Context, count, from *int) ([]v1.CloudHSMSoftwareLicense, error) {
	resp, err := op.client.ListCloudHSMLicenses(
		ctx,
		v1.ListCloudHSMLicensesParams{
			Count: into.Opt[v1.OptInt](count),
			From:  into.Opt[v1.OptInt](from),
		},
	)
	if err != nil {
		return nil, NewAPIError("License.List", 0, err)
	}
	return resp.Licenses, nil
}

type CloudHSMSoftwareLicenseCreateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *LicenseOp) Create(ctx context.Context, p CloudHSMSoftwareLicenseCreateParams) (*v1.CreateCloudHSMSoftwareLicense, error) {
	resp, err := op.client.CreateCloudHSMLicense(
		ctx,
		&v1.WrappedCreateCloudHSMSoftwareLicenseRequest{
			License: v1.CreateCloudHSMSoftwareLicenseRequest{
				ServiceClass: v1.CreateCloudHSMSoftwareLicenseRequestServiceClassCloudCloudhsmLicenseL7,
				Name:         p.Name,
				Description:  into.Opt[v1.OptString](p.Description),
				Tags:         into.OptNilArray[v1.OptNilStringArray](&p.Tags),
			},
		},
	)

	if err == nil {
		ret := resp.GetLicense()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("License.Create", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("License.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("License.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *LicenseOp) Read(ctx context.Context, id string) (*v1.CloudHSMSoftwareLicense, error) {
	resp, err := op.client.ReadCloudHSMLicense(
		ctx,
		v1.ReadCloudHSMLicenseParams{
			ResourceID: id,
		},
	)

	if err == nil {
		ret := resp.GetLicense()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("License.Read", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("License.Read", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return nil, NewAPIError("License.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

type CloudHSMSoftwareLicenseUpdateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *LicenseOp) Update(ctx context.Context, id string, p CloudHSMSoftwareLicenseUpdateParams) (*v1.CloudHSMSoftwareLicense, error) {
	resp, err := op.client.UpdateCloudHSMLicense(
		ctx,
		&v1.WrappedCloudHSMSoftwareLicenseRequest{
			License: v1.CloudHSMSoftwareLicenseRequest{
				Name:        p.Name,
				Description: into.Opt[v1.OptString](p.Description),
				Tags:        into.OptNilArray[v1.OptNilStringArray](&p.Tags),
			},
		},
		v1.UpdateCloudHSMLicenseParams{
			ResourceID: id,
		},
	)

	if err == nil {
		ret := resp.GetLicense()
		return &ret, nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("License.Update", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return nil, NewAPIError("License.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return nil, NewAPIError("License.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *LicenseOp) Delete(ctx context.Context, id string) error {
	err := op.client.DeleteCloudHSMLicense(
		ctx,
		v1.DeleteCloudHSMLicenseParams{
			ResourceID: id,
		},
	)

	if err == nil {
		return nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return NewAPIError("License.Delete", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return NewAPIError("License.Delete", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return NewAPIError("License.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}
