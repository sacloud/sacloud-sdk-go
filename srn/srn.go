// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package srn

import (
	"fmt"
	"strconv"
	"strings"
)

// SRN struct stores the fields of Sakura resource name.
type SRN struct {
	// Version of the SRN. Only integer values are expected, derived from the srnvX.
	Version int32
	// Location of the resource, e.g. is1a/sakura/etc.
	Location string
	// Resource identifier for the service. It can be a string like sakura.gslb or sakura.iaas.server.
	Resource string
	// The ID assigned to the resource. It can be a number/sqids/uuid, so it is stored as a string.
	ID string
}

// Parse parses a string in the form srnv<version>:<location>:<resource>:<id>.
func Parse(str string) (SRN, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 4 {
		return SRN{}, fmt.Errorf("invalid SRN format: %q", str)
	}

	versionPart := parts[0]
	if versionPart != "srnv1" { // TODO: support future versions
		return SRN{}, fmt.Errorf("invalid SRN version prefix: %q", versionPart)
	}

	version, _ := strconv.ParseInt(strings.TrimPrefix(versionPart, "srnv"), 10, 32)
	srn := SRN{
		Version:  int32(version),
		Location: parts[1],
		Resource: parts[2],
		ID:       parts[3],
	}
	if err := srn.Validate(); err != nil {
		return SRN{}, err
	}
	return srn, nil
}

// IsSRN returns whether the given value is an SRN or not.
func IsSRN(str string) bool {
	_, err := Parse(str)
	return err == nil
}

// String returns the canonical SRN string representation.
func (s SRN) String() string {
	return fmt.Sprintf("srnv%d:%s:%s:%s", s.Version, s.Location, s.Resource, s.ID)
}

// Validate checks if the SRN fields are valid or not. It returns an error if any field is invalid.
func (s SRN) Validate() error {
	if s.Version != 1 {
		return fmt.Errorf("invalid SRN version: %d", s.Version)
	}
	locLen := len(s.Location)
	if locLen < 4 || 21 < locLen {
		return fmt.Errorf("invalid SRN: invalid location: %s", s.Location)
	}
	resourceLen := len(s.Resource)
	if resourceLen < 1 || 256 < resourceLen {
		return fmt.Errorf("invalid SRN: invalid resource: %s", s.Resource)
	}
	idLen := len(s.ID)
	if idLen < 1 || 99 < idLen {
		return fmt.Errorf("invalid SRN: invalid ID: %s", s.ID)
	}
	return nil
}
