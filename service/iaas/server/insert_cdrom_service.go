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

package server

import (
	"context"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) InsertCDROM(req *InsertCDROMRequest) error {
	return s.InsertCDROMWithContext(context.Background(), req)
}

func (s *Service) InsertCDROMWithContext(ctx context.Context, req *InsertCDROMRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	client := iaas.NewServerOp(s.caller)
	return client.InsertCDROM(ctx, req.Zone, req.ID, &iaas.InsertCDROMRequest{ID: req.CDROMID})
}
