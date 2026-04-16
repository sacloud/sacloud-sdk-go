// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package define

import (
	"net/http"

	"github.com/sacloud/iaas-api-go/internal/define/names"
	"github.com/sacloud/iaas-api-go/internal/define/ops"
	"github.com/sacloud/iaas-api-go/internal/dsl"
	"github.com/sacloud/iaas-api-go/internal/dsl/meta"
	"github.com/sacloud/iaas-api-go/naked"
	"github.com/sacloud/iaas-api-go/types"
)

const (
	diskAPIName     = "Disk"
	diskAPIPathName = "disk"
)

var diskAPI = &dsl.Resource{
	Name:       diskAPIName,
	PathName:   diskAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	Operations: dsl.Operations{
		// find
		ops.Find(diskAPIName, diskNakedType, findParameter, diskModel),

		// create
		{
			ResourceName: diskAPIName,
			Name:         "Create",
			PathFormat:   dsl.DefaultPathFormat,
			Method:       http.MethodPost,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: diskNakedType,
					Name: "Disk",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskDistantFromType,
					Name: "DistantFrom",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.KMSKey{}),
					Name: "KMSKey",
				},
			),
			Arguments: dsl.Arguments{
				{
					Name:       "createParam",
					MapConvTag: "Disk,recursive",
					Type:       diskCreateParam,
				},
				{
					Name:       "distantFrom",
					MapConvTag: "DistantFrom",
					Type:       diskDistantFromType,
				},
				{
					Name:       "kmeKeyID",
					MapConvTag: "KMSKey.ID",
					Type:       meta.TypeID,
				},
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: diskNakedType,
				Name: "Disk",
			}),
			Results: dsl.Results{
				{
					SourceField: "Disk",
					DestField:   diskModel.Name,
					IsPlural:    false,
					Model:       diskModel,
				},
			},
		},

		// create disk on dedicated storage
		{
			ResourceName: diskAPIName,
			Name:         "CreateOnDedicatedStorage",
			PathFormat:   dsl.DefaultPathFormat,
			Method:       http.MethodPost,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: diskNakedType,
					Name: "Disk",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskDistantFromType,
					Name: "DistantFrom",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.KMSKey{}),
					Name: "KMSKey",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.DedicatedStorageContract{}),
					Name: "TargetDedicatedStorageContract",
				},
			),
			Arguments: dsl.Arguments{
				{
					Name:       "createParam",
					MapConvTag: "Disk,recursive",
					Type:       diskCreateParam,
				},
				{
					Name:       "distantFrom",
					MapConvTag: "DistantFrom",
					Type:       diskDistantFromType,
				},
				{
					Name:       "kmeKeyID",
					MapConvTag: "KMSKey.ID",
					Type:       meta.TypeID,
				},
				{
					Name:       "dedicatedStorageContractID",
					MapConvTag: "TargetDedicatedStorageContract.ID",
					Type:       meta.TypeID,
				},
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: diskNakedType,
				Name: "Disk",
			}),
			Results: dsl.Results{
				{
					SourceField: "Disk",
					DestField:   diskModel.Name,
					IsPlural:    false,
					Model:       diskModel,
				},
			},
		},

		// config(DiskEdit)
		{
			ResourceName:    diskAPIName,
			Name:            "Config",
			PathFormat:      dsl.IDAndSuffixPathFormat("config"),
			Method:          http.MethodPut,
			RequestEnvelope: dsl.RequestEnvelopeFromModel(diskEditParam),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				dsl.PassthroughModelArgument("edit", diskEditParam),
			},
		},

		// create with config(DiskEdit)
		{
			ResourceName: diskAPIName,
			Name:         "CreateWithConfig",
			PathFormat:   dsl.DefaultPathFormat,
			Method:       http.MethodPost,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: diskNakedType,
					Name: "Disk",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskEditNakedType,
					Name: "Config",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.TypeFlag,
					Name: "BootAtAvailable",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskDistantFromType,
					Name: "DistantFrom",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.KMSKey{}),
					Name: "KMSKey",
				},
			),
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: diskNakedType,
				Name: "Disk",
			}),
			Arguments: dsl.Arguments{
				{
					Name:       "createParam",
					MapConvTag: "Disk,recursive",
					Type:       diskCreateParam,
				},
				{
					Name:       "editParam",
					MapConvTag: "Config,recursive",
					Type:       diskEditParam,
				},
				{
					Name:       "bootAtAvailable",
					Type:       meta.TypeFlag,
					MapConvTag: "BootAtAvailable",
				},
				{
					Name:       "distantFrom",
					Type:       diskDistantFromType,
					MapConvTag: "DistantFrom",
				},
				{
					Name:       "kmeKeyID",
					MapConvTag: "KMSKey.ID",
					Type:       meta.TypeID,
				},
			},
			Results: dsl.Results{
				{
					SourceField: "Disk",
					DestField:   diskModel.Name,
					IsPlural:    false,
					Model:       diskModel,
				},
			},
		},

		// create disk on dedicated storage with config(DiskEdit)
		{
			ResourceName: diskAPIName,
			Name:         "CreateOnDedicatedStorageWithConfig",
			PathFormat:   dsl.DefaultPathFormat,
			Method:       http.MethodPost,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: diskNakedType,
					Name: "Disk",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskEditNakedType,
					Name: "Config",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.TypeFlag,
					Name: "BootAtAvailable",
				},
				&dsl.EnvelopePayloadDesc{
					Type: diskDistantFromType,
					Name: "DistantFrom",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.KMSKey{}),
					Name: "KMSKey",
				},
				&dsl.EnvelopePayloadDesc{
					Type: meta.Static(&naked.DedicatedStorageContract{}),
					Name: "TargetDedicatedStorageContract",
				},
			),
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: diskNakedType,
				Name: "Disk",
			}),
			Arguments: dsl.Arguments{
				{
					Name:       "createParam",
					MapConvTag: "Disk,recursive",
					Type:       diskCreateParam,
				},
				{
					Name:       "editParam",
					MapConvTag: "Config,recursive",
					Type:       diskEditParam,
				},
				{
					Name:       "bootAtAvailable",
					Type:       meta.TypeFlag,
					MapConvTag: "BootAtAvailable",
				},
				{
					Name:       "distantFrom",
					Type:       diskDistantFromType,
					MapConvTag: "DistantFrom",
				},
				{
					Name:       "kmeKeyID",
					MapConvTag: "KMSKey.ID",
					Type:       meta.TypeID,
				},
				{
					Name:       "dedicatedStorageContractID",
					MapConvTag: "TargetDedicatedStorageContract.ID",
					Type:       meta.TypeID,
				},
			},
			Results: dsl.Results{
				{
					SourceField: "Disk",
					DestField:   diskModel.Name,
					IsPlural:    false,
					Model:       diskModel,
				},
			},
		},

		// resize partition
		{
			ResourceName: diskAPIName,
			Name:         "ResizePartition",
			PathFormat:   dsl.IDAndSuffixPathFormat("resize-partition"),
			Method:       http.MethodPut,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: meta.TypeFlag,
					Name: "Background",
				},
			),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				dsl.PassthroughModelArgument("param", &dsl.Model{
					Name: "DiskResizePartitionRequest",
					Fields: []*dsl.FieldDesc{
						fields.Def("Background", meta.TypeFlag),
					},
					NakedType: meta.Static(naked.ResizePartitionRequest{}),
				}),
			},
		},

		// connect to server
		ops.WithIDAction(diskAPIName, "ConnectToServer", http.MethodPut, "to/server/{{.serverID}}",
			&dsl.Argument{
				Name: "serverID",
				Type: meta.TypeID,
			},
		),

		// disconnect from server
		ops.WithIDAction(diskAPIName, "DisconnectFromServer", http.MethodDelete, "to/server"),

		// read
		ops.Read(diskAPIName, diskNakedType, diskModel),

		// update
		ops.Update(diskAPIName, diskNakedType, diskUpdateParam, diskModel),

		// delete
		ops.Delete(diskAPIName),

		// monitor
		ops.Monitor(diskAPIName, monitorParameter, monitors.diskModel()),
		ops.MonitorChild(diskAPIName, "Disk", "", monitorParameter, monitors.diskModel()),
	},
}

var (
	diskNakedType       = meta.Static(naked.Disk{})
	diskEditNakedType   = meta.Static(naked.DiskEdit{})
	diskDistantFromType = meta.Static([]types.ID{})

	diskModel = &dsl.Model{
		Name:      diskAPIName,
		NakedType: diskNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.Availability(),
			fields.DiskConnection(),
			fields.DiskConnectionOrder(),
			fields.DiskEncryptionAlgorithm(),
			fields.KMSKeyID(),
			fields.DiskReinstallCount(),
			fields.Def("JobStatus", models.migrationJobStatus()),
			fields.SizeMB(),
			fields.MigratedMB(),
			fields.DiskPlanID(),
			fields.DiskPlanName(),
			fields.DiskPlanStorageClass(),
			fields.SourceDiskID(),
			fields.SourceDiskAvailability(),
			fields.SourceArchiveID(),
			fields.SourceArchiveAvailability(),
			fields.BundleInfo(),
			fields.Storage(),
			fields.ServerID(),
			fields.ServerName(),
			fields.IconID(),
			fields.CreatedAt(),
			fields.ModifiedAt(),
		},
	}

	diskCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(diskAPIName),
		NakedType: diskNakedType,
		Fields: []*dsl.FieldDesc{
			fields.DiskPlanID(),
			fields.DiskConnection(),
			fields.DiskEncryptionAlgorithm(),
			fields.SourceDiskID(),
			fields.SourceArchiveID(),
			fields.ServerID(),
			fields.SizeMB(),
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),
		},
	}

	diskUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(diskAPIName),
		NakedType: diskNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),
			fields.DiskConnection(),
		},
	}

	diskEditParam = models.diskEdit()
)
