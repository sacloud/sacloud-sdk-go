// Copyright 2025- The sacloud/saclient-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/sacloud/iam-api-go"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func iamrole(name string) (this v1.IamPolicyRole) {
	this.SetType(v1.NewOptIamPolicyRoleType(v1.IamPolicyRoleTypePreset))
	this.SetID(v1.NewOptString(name))
	return
}

func idrole(name string) (this v1.IdPolicyRole) {
	this.SetType(v1.NewOptIdPolicyRoleType(v1.IdPolicyRoleTypePreset))
	this.SetID(v1.NewOptString(name))
	return
}

func principalOf(id int) (this v1.Principal) {
	this.SetType(v1.NewOptString("service-principal"))
	this.SetID(v1.NewOptInt(id))
	return
}

func iampolicy(name string, id int) (this v1.IamPolicy) {
	var role = iamrole(name)
	var principal = principalOf(id)
	var principals = make([]v1.Principal, 0, 1)

	this.SetRole(v1.NewOptIamPolicyRole(role))
	this.SetPrincipals(append(principals, principal))
	return
}

func idpolicy(name string, id int) (this v1.IdPolicy) {
	var role = idrole(name)
	var principal = principalOf(id)
	var principals = make([]v1.Principal, 0, 1)

	this.SetRole(v1.NewOptIdPolicyRole(role))
	this.SetPrincipals(append(principals, principal))
	return
}

func inspectPrincipal(p *v1.Principal) {
	fmt.Printf("(")
	if t, ok := p.GetType().Get(); ok {
		fmt.Printf("%s: ", t)
	}
	if i, ok := p.GetID().Get(); ok {
		fmt.Printf("%d", i)
	}
	fmt.Printf(")")
}

func inspectIAMRole(r *v1.IamPolicyRole) {
	fmt.Printf("Role (")
	if t, ok := r.GetType().Get(); ok {
		fmt.Printf("%s:", t)
	}
	if i, ok := r.GetID().Get(); ok {
		fmt.Printf(" \"%s\"", i)
	}
	fmt.Printf(")")
}

func inspectIDRole(r *v1.IdPolicyRole) {
	fmt.Printf("Role (")
	if t, ok := r.GetType().Get(); ok {
		fmt.Printf("%s:", t)
	}
	if i, ok := r.GetID().Get(); ok {
		fmt.Printf(" \"%s\"", i)
	}
	fmt.Printf(")")
}

func organizationIAMPolicies(ctx context.Context, client *v1.Client, id int) {
	var op iam.IAMPolicyAPI = iam.NewIAMPolicyOp(client)
	roles := []string{
		"owner",
		"organization-admin",
		"servicepolicy-admin",
		"folder-admin",
		"project-creator",
	}

	if policies, err := op.ReadOrganizationPolicy(ctx); err != nil {
		panic(err)
	} else {
		for _, p := range roles {
			policies = append(policies, iampolicy(p, id))
		}
		if actual, err := op.UpdateOrganizationPolicy(ctx, policies); err != nil {
			panic(err)
		} else {
			fmt.Printf("Current (updated) IAM policies assigned to the organization:\n")
			for _, p := range actual {
				fmt.Printf("- ")
				if r, ok := p.GetRole().Get(); ok {
					inspectIAMRole(&r)
				}
				fmt.Printf(" => ")
				for _, pr := range p.GetPrincipals() {
					inspectPrincipal(&pr)
				}
				fmt.Printf("\n")
			}
		}
	}
}

func organizationIDPolicies(ctx context.Context, client *v1.Client, id int) {
	var op iam.IDPolicyAPI = iam.NewIDPolicyOp(client)
	roles := []string{
		"identity-admin",
	}

	if policies, err := op.ReadOrganizationIdPolicy(ctx); err != nil {
		panic(err)
	} else {
		for _, p := range roles {
			policies = append(policies, idpolicy(p, id))
		}
		if actual, err := op.UpdateOrganizationIdPolicy(ctx, policies); err != nil {
			panic(err)
		} else {
			fmt.Printf("Current (updated) ID policies assigned to the organization:\n")
			for _, p := range actual {
				fmt.Printf("- ")
				if r, ok := p.GetRole().Get(); ok {
					inspectIDRole(&r)
				}
				fmt.Printf(" => ")
				for _, pr := range p.GetPrincipals() {
					inspectPrincipal(&pr)
				}
				fmt.Printf("\n")
			}
		}
	}
}

func main() {
	ctx := context.Background()
	tok := flag.String("token", "", "developer token")
	principalId := flag.Int("principal", 0, "resource # of principal")

	flag.Parse()

	if *tok == "" {
		panic("token is required")
	}
	if *principalId == 0 {
		panic("principal is required")
	}

	if err := theClient.SetWith(saclient.WithBasicAuth1(*tok)); err != nil {
		panic(err)
	}

	if err := theClient.Populate(); err != nil {
		panic(err)
	}

	v1, err := iam.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	organizationIAMPolicies(ctx, v1, *principalId)
	fmt.Printf("\n")
	organizationIDPolicies(ctx, v1, *principalId)
}
