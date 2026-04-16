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

package validate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
)

type FormatErrorFunc func(target interface{}, err validator.FieldError) string

type Validator struct {
	FormatErrorFuncMap map[string]FormatErrorFunc

	instance *validator.Validate
	initOnce sync.Once
}

func New() *Validator {
	v := &Validator{
		FormatErrorFuncMap: map[string]FormatErrorFunc{
			"file": func(_ interface{}, err validator.FieldError) string {
				return fmt.Sprintf("invalid file path: %v", err.Value())
			},
		},
	}
	v.init()
	return v
}

func (v *Validator) init() {
	v.initOnce.Do(func() {
		if v.instance == nil {
			v.instance = validator.New()
		}
		v.instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
			if name == "" {
				// nameタグがない場合はyamlタグを参照
				name = strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})
	})
}

// RegisterCollectionValidator 要素/コレクションを示す単数名/複数名に対しそれぞれoneof/diveエイリアスを登録する
//
// ("zone", "zones", []string{"v1", "v2"}) とした場合、以下のエイリアスがバリデーターに登録される
//   - "zone"  => "oneof=v1 v2"のエイリアス
//   - "zones" => "dive,zone"のエイリアス
func (v *Validator) RegisterCollectionValidator(singularName, pluralName string, allowedValues []string) {
	v.init()

	v.instance.RegisterAlias(singularName, fmt.Sprintf("oneof=%s", strings.Join(allowedValues, " ")))
	v.instance.RegisterAlias(pluralName, fmt.Sprintf("dive,%s", singularName))
}

// Struct 対象structを検証しerrorを返す
func (v *Validator) Struct(value interface{}) error {
	v.init()

	errors := v.StructWithMultiError(value)
	if errors != nil {
		return errors.ErrorOrNil()
	}
	return nil
}

// StructWithMultiError 対象structを検証し、*multierror.Errorを返す
func (v *Validator) StructWithMultiError(value interface{}) *multierror.Error {
	v.init()

	err := v.instance.Struct(value)
	if err != nil {
		if err != nil {
			// see https://github.com/go-playground/validator/blob/f6584a41c8acc5dfc0b62f7962811f5231c11530/_examples/simple/main.go#L59-L65
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return &multierror.Error{Errors: []error{err}}
			}

			errors := &multierror.Error{}
			for _, err := range err.(validator.ValidationErrors) {
				errors = multierror.Append(errors, v.errorFromValidationErr(v, err))
			}
			return errors
		}
	}

	return nil
}

func (v *Validator) errorFromValidationErr(target interface{}, err validator.FieldError) error {
	namespaces := strings.Split(err.Namespace(), ".")
	actualName := namespaces[len(namespaces)-1] // .で区切った末尾の要素

	param := err.Param()
	detail := err.ActualTag()
	if param != "" {
		detail += "=" + param
	}

	return v.newError(actualName, v.formatErrorDetail(detail, target, err))
}

func (v *Validator) formatErrorDetail(detail string, target interface{}, err validator.FieldError) string {
	if fn, ok := v.FormatErrorFuncMap[detail]; ok {
		return fn(target, err)
	}
	return detail
}

func (v *Validator) newError(name, message string) error {
	return fmt.Errorf("%s: %s", name, message)
}
