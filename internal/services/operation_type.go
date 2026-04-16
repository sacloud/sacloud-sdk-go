// Copyright 2022-2025 The sacloud/services Authors
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

package services

// OperationType サービスがサポートする操作の種別
type OperationType int

const (
	OperationTypeUnknown OperationType = iota
	OperationTypeCreate                // Id:不要	/ 戻り値:単体
	OperationTypeRead                  // Id:要 	/ 戻り値:単体
	OperationTypeUpdate                // Id:要 	/ 戻り値:単体
	OperationTypeDelete                // Id:要 	/ 戻り値:なし
	OperationTypeList                  // Id:不要	/ 戻り値: スライス
	OperationTypeAction                // Id:要 	/ 戻り値: なし
)

func (o OperationType) String() string {
	switch o {
	case OperationTypeUnknown:
		return "unknown"
	case OperationTypeCreate:
		return "create"
	case OperationTypeRead:
		return "read"
	case OperationTypeUpdate:
		return "update"
	case OperationTypeDelete:
		return "delete"
	case OperationTypeList:
		return "list"
	case OperationTypeAction:
		return "action"
	default:
		panic("got unknown OperationType")
	}
}

func (o OperationType) HasReturnValue() bool {
	switch o {
	case OperationTypeCreate, OperationTypeRead, OperationTypeUpdate, OperationTypeList:
		return true
	case OperationTypeAction, OperationTypeDelete:
		return false
	default:
		panic("got unknown OperationType")
	}
}
