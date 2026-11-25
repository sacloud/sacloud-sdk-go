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
	resp, err := errorFromDecodedResponse("License.List", func() (*v1.PaginatedCloudHSMSoftwareLicenseList, error) {
		return op.client.ListCloudHSMLicenses(ctx, v1.ListCloudHSMLicensesParams{
			Count: into.Opt[v1.OptInt](count),
			From:  into.Opt[v1.OptInt](from),
		})
	})
	if err != nil {
		return nil, err
	}

	return resp.GetLicenses(), nil
}

type CloudHSMSoftwareLicenseCreateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *LicenseOp) Create(ctx context.Context, p CloudHSMSoftwareLicenseCreateParams) (*v1.CreateCloudHSMSoftwareLicense, error) {
	resp, err := errorFromDecodedResponse("License.Create", func() (*v1.WrappedCreateCloudHSMSoftwareLicense, error) {
		return op.client.CreateCloudHSMLicense(ctx, &v1.WrappedCreateCloudHSMSoftwareLicenseRequest{
			License: v1.CreateCloudHSMSoftwareLicenseRequest{
				ServiceClass: v1.CreateCloudHSMSoftwareLicenseRequestServiceClassCloudCloudhsmLicenseL7,
				Name:         p.Name,
				Description:  into.Opt[v1.OptString](p.Description),
				Tags:         into.OptNonnilArray[v1.OptNilStringArray](&p.Tags),
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetLicense()), nil
}

func (op *LicenseOp) Read(ctx context.Context, id string) (*v1.CloudHSMSoftwareLicense, error) {
	resp, err := errorFromDecodedResponse("License.Read", func() (*v1.WrappedCloudHSMSoftwareLicense, error) {
		return op.client.ReadCloudHSMLicense(ctx, v1.ReadCloudHSMLicenseParams{
			ResourceID: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetLicense()), nil
}

type CloudHSMSoftwareLicenseUpdateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *LicenseOp) Update(ctx context.Context, id string, p CloudHSMSoftwareLicenseUpdateParams) (*v1.CloudHSMSoftwareLicense, error) {
	resp, err := errorFromDecodedResponse("License.Update", func() (*v1.WrappedCloudHSMSoftwareLicense, error) {
		return op.client.UpdateCloudHSMLicense(ctx, &v1.WrappedCloudHSMSoftwareLicenseRequest{
			License: v1.CloudHSMSoftwareLicenseRequest{
				Name:        p.Name,
				Description: into.Opt[v1.OptString](p.Description),
				Tags:        into.OptNonnilArray[v1.OptNilStringArray](&p.Tags),
			},
		}, v1.UpdateCloudHSMLicenseParams{
			ResourceID: id,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetLicense()), nil
}

func (op *LicenseOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("License.Delete", func() error {
		return op.client.DeleteCloudHSMLicense(ctx, v1.DeleteCloudHSMLicenseParams{
			ResourceID: id,
		})
	})
}
