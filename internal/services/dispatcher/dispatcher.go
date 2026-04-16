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

package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/sacloud/services"
	"github.com/sacloud/services/helper"
)

var registry = map[string]services.Service{}

func Register(platformName string, ss services.Services) {
	if platformName == "" {
		panic("platformName is required")
	}
	for _, s := range ss {
		args := []string{platformName, s.Info().FullName()}
		registry[key(args)] = s
	}
}

func Dispatch(arguments []string, parameter interface{}) (interface{}, error) {
	return DispatchWithContext(context.Background(), arguments, parameter)
}

func DispatchWithContext(ctx context.Context, arguments []string, parameter interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		panic("invalid arguments")
	}
	keys, operation := arguments[:len(arguments)-1], arguments[len(arguments)-1]
	service, ok := registry[key(keys)]
	if !ok {
		return nil, fmt.Errorf("service %s not found", key(keys))
	}
	for _, op := range service.Operations() {
		if op.EqualsByName(operation) {
			return dispatch(ctx, service, op, parameter)
		}
	}
	return nil, fmt.Errorf("operation %s#%s not found", key(keys), operation)
}

func serviceParameter(service services.Service, op services.SupportedOperation, parameter interface{}) (interface{}, error) {
	if parameter == nil {
		parameter = map[string]interface{}{}
	}

	param, err := helper.NewParameter(service, op.Name)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(parameter)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, param); err != nil {
		return nil, err
	}
	return param, nil
}

func dispatch(ctx context.Context, service services.Service, op services.SupportedOperation, parameter interface{}) (interface{}, error) {
	param, err := serviceParameter(service, op, parameter)
	if err != nil {
		return nil, err
	}

	method := reflect.ValueOf(service).MethodByName(op.WithContextFuncName())
	results := method.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)}) // xxxWithContextはctx+reqの2つだけを受け取るはず

	var value interface{}
	if op.OperationType.HasReturnValue() {
		if len(results) != 2 {
			return nil, fmt.Errorf("invalid results: want 2 results, but got %d: %+v", len(results), results)
		}
		// shift
		value = results[0].Interface()
		results = results[1:]
	} else if len(results) != 1 {
		return nil, fmt.Errorf("invalid results: want 1 results, but got %d: %+v", len(results), results)
	}

	if e, ok := results[0].Interface().(error); ok {
		err = e
	}
	return value, err
}

func key(arguments []string) string {
	return strings.Join(arguments, "/")
}
