package nosql_test

import (
	"reflect"
	"testing"

	"github.com/sacloud/nosql-api-go"
	v1 "github.com/sacloud/nosql-api-go/apis/v1"
)

func TestPlan_AllValues(t *testing.T) {
	var p nosql.Plan
	got := p.AllValues()
	want := []nosql.Plan{nosql.Plan40GB, nosql.Plan100GB, nosql.Plan250GB}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AllValues() = %v, want %v", got, want)
	}

	gotStr := p.AllValuesAsString()
	wantStr := []string{string(nosql.Plan40GB), string(nosql.Plan100GB), string(nosql.Plan250GB)}
	if !reflect.DeepEqual(gotStr, wantStr) {
		t.Errorf("AllValuesAsString() = %v, want %v", gotStr, wantStr)
	}
}

func TestPlan_Attributes(t *testing.T) {
	cases := []struct {
		name         string
		p            nosql.Plan
		planID       int
		serviceClass v1.ServiceClass
		memoryMB     int
		diskMB       int
		vcore        int
		nodes        int
	}{
		{"40GB", nosql.Plan40GB, 51142, v1.ServiceClass("cloud/nosql/plan/1"), 4096, 40960, 2, 1},
		{"100GB", nosql.Plan100GB, 51143, v1.ServiceClass("cloud/nosql/plan/2"), 8192, 102400, 3, 3},
		{"250GB", nosql.Plan250GB, 51144, v1.ServiceClass("cloud/nosql/plan/3"), 16384, 256000, 6, 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.GetPlanID(); got != tc.planID {
				t.Errorf("GetPlanID(%s) = %d, want %d", tc.p, got, tc.planID)
			}
			if got := tc.p.GetServiceClass(); got != tc.serviceClass {
				t.Errorf("GetServiceClass(%s) = %q, want %q", tc.p, got, tc.serviceClass)
			}
			if got := tc.p.GetMemoryMB(); got != tc.memoryMB {
				t.Errorf("GetMemoryMB(%s) = %d, want %d", tc.p, got, tc.memoryMB)
			}
			if got := tc.p.GetDiskSizeMB(); got != tc.diskMB {
				t.Errorf("GetDiskSizeMB(%s) = %d, want %d", tc.p, got, tc.diskMB)
			}
			if got := tc.p.GetVirtualCore(); got != tc.vcore {
				t.Errorf("GetVirtualCore(%s) = %d, want %d", tc.p, got, tc.vcore)
			}
			if got := tc.p.GetNodes(); got != tc.nodes {
				t.Errorf("GetNodes(%s) = %d, want %d", tc.p, got, tc.nodes)
			}
		})
	}
}

func TestPlan_ForNodesAttributes(t *testing.T) {
	cases := []struct {
		name           string
		p              nosql.Plan
		planIDForNodes int
		serviceClass   v1.ServiceClass
		nodesForNodes  int
	}{
		{"40GB", nosql.Plan40GB, 0, v1.ServiceClass(""), 0},
		{"100GB", nosql.Plan100GB, 51145, v1.ServiceClass("cloud/nosql/plan/2/node"), 2},
		{"250GB", nosql.Plan250GB, 51146, v1.ServiceClass("cloud/nosql/plan/3/node"), 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.p.GetPlanIDforNodes(); got != tc.planIDForNodes {
				t.Errorf("GetPlanIDforNodes(%s) = %d, want %d", tc.p, got, tc.planIDForNodes)
			}
			if got := tc.p.GetServiceClassForNodes(); got != tc.serviceClass {
				t.Errorf("GetServiceClassForNodes(%s) = %q, want %q", tc.p, got, tc.serviceClass)
			}
			if got := tc.p.GetNodesForNodes(); got != tc.nodesForNodes {
				t.Errorf("GetNodesForNodes(%s) = %d, want %d", tc.p, got, tc.nodesForNodes)
			}
		})
	}
}
