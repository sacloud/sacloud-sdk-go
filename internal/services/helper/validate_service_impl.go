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

package helper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/sacloud/services"
)

func TestServiceImplementation(t *testing.T, ss ...services.Service) bool {
	for _, s := range ss {
		err, warn := testServiceImpl(s)
		if err != nil {
			t.Error(err)
		}
		if warn != nil {
			t.Logf("[WARN] %s", warn)
		}
	}
	return t.Failed()
}

func testServiceImpl(svc services.Service) (error, error) {
	errors := &multierror.Error{}
	warnings := &multierror.Error{}

	// 見出しにするためにInfo().Nameは先にチェックしておく
	info := svc.Info()
	if info == nil {
		errors = multierror.Append(fmt.Errorf("requires a return value: Info()"))
		return errors.ErrorOrNil(), warnings.ErrorOrNil()
	}
	if info.Name == "" {
		errors = multierror.Append(fmt.Errorf("requires a return value: Info().Name"))
		return errors.ErrorOrNil(), warnings.ErrorOrNil()
	}

	validateConfig(svc, errors, warnings)
	validateOperations(svc, errors, warnings)

	return multierror.Prefix(errors.ErrorOrNil(), fmt.Sprintf("services[%s]: ", info.Name)), warnings.ErrorOrNil()
}

func appendErrors(err *multierror.Error, errors ...error) {
	newError := multierror.Append(err, errors...)
	*err = *newError
}

func validateConfig(svc services.Service, errors, warnings *multierror.Error) {
	config := svc.Config()
	if config == nil {
		appendErrors(errors, fmt.Errorf("requires a return value: Config()"))
	}
}

func validateOperations(svc services.Service, errors, warnings *multierror.Error) {
	operations := svc.Operations()
	if len(operations) == 0 {
		appendErrors(errors, fmt.Errorf("requires a return value: Operations()"))
	}

	names := make(map[string]struct{})
	for _, op := range operations {
		validateOperation(svc, op, errors, warnings)

		// 重複チェック
		_, exist := names[op.Name]
		if exist {
			appendErrors(errors, fmt.Errorf("unique value required: operation %s is duplicated", op.Name))
		}
		names[op.Name] = struct{}{}
	}
}

func validateOperation(svc services.Service, op services.SupportedOperation, errors, warnings *multierror.Error) {
	if op.Name == "" {
		appendErrors(errors, fmt.Errorf("required: Operations().Name"))
	}
	if op.OperationType == services.OperationTypeUnknown {
		appendErrors(errors, fmt.Errorf("operation[%s]: value must be set: Operations().OperationType", op.Name))
	}

	// funcが定義されているか?
	method, exist := reflect.TypeOf(svc).MethodByName(op.FuncName())
	if !exist {
		appendErrors(errors, fmt.Errorf("func %s() required", op.FuncName()))
		return
	}
	methodWithContext, exist := reflect.TypeOf(svc).MethodByName(op.WithContextFuncName())
	if !exist {
		appendErrors(errors, fmt.Errorf("func %s() required", op.WithContextFuncName()))
		return
	}

	// 引数
	if method.Type.NumIn() != 2 {
		appendErrors(errors, fmt.Errorf("func %s() must have 1 argument", op.FuncName()))
	}
	if methodWithContext.Type.NumIn() != 3 {
		appendErrors(errors, fmt.Errorf("func %s() must have 2 arguments(ctx + request)", op.FuncName()))
	}

	// 戻り値
	for _, m := range []reflect.Method{method, methodWithContext} {
		switch op.OperationType {
		case services.OperationTypeList:
			// 戻り値は[]xxx + error
			if m.Type.NumOut() != 2 {
				appendErrors(errors, fmt.Errorf("func %s() must return value + error", m.Name))
				break
			}

		case services.OperationTypeCreate, services.OperationTypeRead, services.OperationTypeUpdate:
			// 戻り値はxxx + error
			if m.Type.NumOut() != 2 {
				appendErrors(errors, fmt.Errorf("func %s() must return value + error", m.Name))
				break
			}
		case services.OperationTypeAction, services.OperationTypeDelete:
			// 戻り値はerrorのみ
			if m.Type.NumOut() != 1 {
				appendErrors(errors, fmt.Errorf("func %s() must return error", m.Name))
				break
			}
		}

		// 最後の戻り値はerror型
		if m.Type.NumOut() > 0 && !m.Type.Out(m.Type.NumOut()-1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			appendErrors(errors, fmt.Errorf("func %s(): return-values[%d]: required error type", m.Name, m.Type.NumOut()-1))
		}
	}
}
