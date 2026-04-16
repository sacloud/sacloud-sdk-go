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

	"github.com/sacloud/services"
)

// PrintServiceMeta サービスメタデータを標準出力に整形して書き込む
func PrintServiceMeta(service services.Service) {
	ops, err := ServiceMeta(service)
	if err != nil {
		panic(err)
	}

	info := service.Info()
	fmt.Println("---")
	fmt.Printf("# Service[%s]\n", info.Name)
	fmt.Printf("Name: %s\n", info.Name)
	fmt.Printf("PkgPath: %s\n", ServicePkgPath(service))
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("ParentServices: %s\n", info.ParentServices)

	fmt.Println("Operations:")
	for _, op := range ops {
		fmt.Printf("  %s(%s):\n", op.Operation.Name, op.Operation.OperationType)
		for _, field := range op.Parameters {
			fmt.Printf("    %s:\n", field.FieldName)
			fmt.Printf("      StructField:\n")
			fmt.Printf("        PkgPath: %s\n", field.StructField.PkgPath)
			fmt.Printf("        Tag: %s\n", field.StructField.Tag)
			fmt.Printf("        Type: %q\n", field.StructField.Type.String())

			fmt.Printf("      Tag:\n")
			fmt.Printf("        FlagName: %s\n", field.Tag.FlagName)
			fmt.Printf("        FieldName: %s\n", field.Tag.FieldName)
			fmt.Printf("        Aliases: %s\n", field.Tag.Aliases)
			fmt.Printf("        Shorthand: %s\n", field.Tag.Shorthand)
			fmt.Printf("        Description: %s\n", field.Tag.Description)
			fmt.Printf("        Squash: %t\n", field.Tag.Squash)
			fmt.Printf("        Ignore: %t\n", field.Tag.Ignore)
			fmt.Printf("        Category: %s\n", field.Tag.Category)
			fmt.Printf("        Order: %d\n", field.Tag.Order)
			fmt.Printf("        Options: %s\n", field.Tag.Options)
			fmt.Printf("        DisplayOptions: %s\n", field.Tag.DisplayOptions)
			fmt.Printf("      LongDescription(): %#+v\n", field.Tag.LongDescription())
			fmt.Printf("      AliasesString(): %#+v\n", field.Tag.AliasesString())
			fmt.Printf("      OptionsString(): %#+v\n", field.Tag.OptionsString())
		}
	}
}
