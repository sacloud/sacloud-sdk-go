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

package disk

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

type dummyDiskPlanReader struct {
	diskPlan *iaas.DiskPlan
	err      error
}

func (d *dummyDiskPlanReader) Read(ctx context.Context, zone string, id types.ID) (*iaas.DiskPlan, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.diskPlan, nil
}

type dummyNoteHandler struct {
	note *iaas.Note
	err  error
}

func (d *dummyNoteHandler) Read(ctx context.Context, id types.ID) (*iaas.Note, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.note, nil
}

func (d *dummyNoteHandler) Create(ctx context.Context, param *iaas.NoteCreateRequest) (*iaas.Note, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.note, nil
}

func (d *dummyNoteHandler) Delete(ctx context.Context, id types.ID) error {
	return d.err
}
