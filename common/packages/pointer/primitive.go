// Copyright 2022-2025 The sacloud/packages-go Authors
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

package pointer

// NewBool returns a pointer to the given bool value.
//
//go:fix inline
func NewBool(b bool) *bool { return new(b) }

// NewString returns a pointer to the given string value.
//
//go:fix inline
func NewString(s string) *string { return new(s) }

// NewInt returns a pointer to the given int value.
//
//go:fix inline
func NewInt(i int) *int { return new(i) }

// NewInt8 returns a pointer to the given int8 value.
//
//go:fix inline
func NewInt8(i int8) *int8 { return new(i) }

// NewInt16 returns a pointer to the given int16 value.
//
//go:fix inline
func NewInt16(i int16) *int16 { return new(i) }

// NewInt32 returns a pointer to the given int32 value.
//
//go:fix inline
func NewInt32(i int32) *int32 { return new(i) }

// NewInt64 returns a pointer to the given int64 value.
//
//go:fix inline
func NewInt64(i int64) *int64 { return new(i) }

// NewUint returns a pointer to the given uint value.
//
//go:fix inline
func NewUint(i uint) *uint { return new(i) }

// NewUint8 returns a pointer to the given uint8 value.
//
//go:fix inline
func NewUint8(i uint8) *uint8 { return new(i) }

// NewUint16 returns a pointer to the given uint16 value.
//
//go:fix inline
func NewUint16(i uint16) *uint16 { return new(i) }

// NewUint32 returns a pointer to the given uint32 value.
//
//go:fix inline
func NewUint32(i uint32) *uint32 { return new(i) }

// NewUint64 returns a pointer to the given uint64 value.
//
//go:fix inline
func NewUint64(i uint64) *uint64 { return new(i) }

// NewFloat32 returns a pointer to the given float32 value.
//
//go:fix inline
func NewFloat32(f float32) *float32 { return new(f) }

// NewFloat64 returns a pointer to the given float64 value.
//
//go:fix inline
func NewFloat64(f float64) *float64 { return new(f) }

// NewByte returns a pointer to the given byte value.
//
//go:fix inline
func NewByte(b byte) *byte { return new(b) }

// NewRune returns a pointer to the given rune value.
//
//go:fix inline
func NewRune(r rune) *rune { return new(r) }
