// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import v1 "github.com/sacloud/nosql-api-go/apis/v1"

type Plan string

const (
	Plan40GB  Plan = "40GB"
	Plan100GB Plan = "100GB"
	Plan250GB Plan = "250GB"
)

func GetPlanFromID(planID int) Plan {
	switch planID {
	case 51142:
		return Plan40GB
	case 51143, 51145:
		return Plan100GB
	case 51144, 51146:
		return Plan250GB
	default:
		return ""
	}
}

func (Plan) AllValues() []Plan {
	return []Plan{
		Plan40GB,
		Plan100GB,
		Plan250GB,
	}
}

func (Plan) AllValuesAsString() []string {
	return []string{
		string(Plan40GB),
		string(Plan100GB),
		string(Plan250GB),
	}
}

func (p Plan) GetPlanID() int {
	switch p {
	case Plan40GB:
		return 51142
	case Plan100GB:
		return 51143
	case Plan250GB:
		return 51144
	default:
		return 0
	}
}

func (p Plan) GetPlanIDforNodes() int {
	switch p {
	case Plan100GB:
		return 51145
	case Plan250GB:
		return 51146
	default:
		return 0
	}
}

func (p Plan) GetServiceClass() v1.ServiceClass {
	switch p {
	case Plan40GB:
		return v1.ServiceClass("cloud/nosql/plan/1")
	case Plan100GB:
		return v1.ServiceClass("cloud/nosql/plan/2")
	case Plan250GB:
		return v1.ServiceClass("cloud/nosql/plan/3")
	default:
		return ""
	}
}

func (p Plan) GetServiceClassForNodes() v1.ServiceClass {
	switch p {
	case Plan100GB:
		return v1.ServiceClass("cloud/nosql/plan/2/node")
	case Plan250GB:
		return v1.ServiceClass("cloud/nosql/plan/3/node")
	default:
		return ""
	}
}

func (p Plan) GetMemoryMB() int {
	switch p {
	case Plan40GB:
		return 4096
	case Plan100GB:
		return 8192
	case Plan250GB:
		return 16384
	default:
		return 0
	}
}

func (p Plan) GetDiskSizeMB() int {
	switch p {
	case Plan40GB:
		return 40960
	case Plan100GB:
		return 102400
	case Plan250GB:
		return 256000
	default:
		return 0
	}
}

func (p Plan) GetVirtualCore() int {
	switch p {
	case Plan40GB:
		return 2
	case Plan100GB:
		return 3
	case Plan250GB:
		return 6
	default:
		return 0
	}
}

func (p Plan) GetNodes() int {
	switch p {
	case Plan100GB, Plan250GB:
		return 3
	default:
		return 1
	}
}

func (p Plan) GetNodesForNodes() int {
	switch p {
	case Plan100GB, Plan250GB:
		return 2
	default:
		return 0
	}
}
