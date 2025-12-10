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

package saclient

import (
	"encoding/json"
	"flag"
	"io"
	"maps"
	"net/http"
)

// The API
type ClientAPI interface {
	// Populate settings from environment variables and flags
	//
	// Note that once populated, the settings gets immutable.
	// Also note that first call to [Do] implicityly implies population,
	// means you cannot change settings afterwards.
	Populate() error

	// A copy without population.  You can modify settings and repopulate.
	Dup() ClientAPI

	// ```golang
	//
	//	import (
	//		"os"
	//
	//	   "github.com/sacloud/saclient-go"
	//	)
	//
	//	var client saclient.Client
	//
	//	func main() {
	//		client.SetEnviron(os.Environ())
	//		client.Populate()
	//		// ...
	//	}
	// ```
	SetEnviron(env []string) error

	// ```golang
	//
	//	import (
	//		"context"
	//		"os"
	//
	//		"github.com/hashicorp/terraform-plugin-framework/provider"
	//		"github.com/sacloud/saclient-go"
	//	)
	//
	//	type providerModel struct {
	//		// ...
	//	}
	//
	//	var _ saclient.TerraformProviderInterface = (*providerModel)(nil)
	//
	//	func (p *provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	//		var config providerModel
	//		var client saclient.Client
	//		diags := req.Config.Get(ctx, &config)
	//		resp.Diagnostics.Append(diags...)
	//		if resp.Diagnostics.HasError() {
	//			return
	//		}
	//
	//		client.SetEnviron(os.Environ())
	//		client.SettingsFromTerraformProvider(&config)
	//		client.Populate()
	//	}
	// ```
	SettingsFromTerraformProvider(config TerraformProviderInterface) error

	// ```golang
	//
	//	import (
	//	   "github.com/spf13/cobra"
	//	   "github.com/spf13/pflag"
	//	   "github.com/sacloud/saclient-go"
	//	)
	//
	//	var client saclient.Client
	//	var command = &cobra.Command{
	//	    RunE: func(cmd *cobra.Command, args []string) error {
	//	        if err := client.Populate(); err != nil {
	//	            return err
	//	        }
	//	        // ...
	//	        return nil
	//	    },
	//	}
	//
	//	func init() {
	//	   client.SetEnviron(os.Environ())
	//	   command.PersistentFlags().AddGoFlagSet(client.FlagSet(flag.PanicOnError))
	//	}
	// ```
	FlagSet(eh flag.ErrorHandling) *flag.FlagSet

	// Returns the currently selected profile, or nil if absent.
	// Profile historically includes much more than client configuration,
	// like usacloud's ArgumentMatchMode.
	//
	// Note that it's completely normal to have nil here.
	// The user can just opt not to have profiles at all.
	// (typical situation for CI environments etc.)
	Profile() (*Profile, error)

	// _Calculated_ name of the current profile.
	//
	// On normal situations this is identical to Profile().Name.  But
	// there are cases when the specified profile does not exist, or
	// is broken.  This method returns something, regardless of such
	// situations.
	//
	// Note however that this can still return nil, in case there is
	// no profile specified at all.
	ProfileName() (dir, name *string)

	// CRUD-style API operator for Profile resource
	ProfileOp() (ProfileAPI, error)

	// HTTP request doer
	Do(req *http.Request) (*http.Response, error)
}

// impmlementation of ClientAPI
type Client struct {
	params parameter
	once   once[inner]
}

func (c *Client) Populate() error {
	_, err := c.ensurePopulated()
	return err
}

func (c *Client) Dup() ClientAPI {
	if c == nil {
		return (ClientAPI)(nil)
	} else {
		return &Client{params: c.params}
	}
}

//nolint:gocritic
func (c *Client) SetEnviron(env []string) error {
	if c == nil {
		return NewErrorf("nil client")
	} else if c.once.Done() {
		return NewErrorf("client already populated; cannot change settings")
	} else {
		return c.params.setEnviron(env)
	}
}

func (c *Client) FlagSet(eh flag.ErrorHandling) *flag.FlagSet {
	if c == nil {
		return nil
	} else {
		return c.params.flagSet(eh)
	}
}

//nolint:gocritic
func (c *Client) SettingsFromTerraformProvider(p TerraformProviderInterface) error {
	if c == nil {
		return NewErrorf("nil client")
	} else if c.once.Done() {
		return NewErrorf("client already populated; cannot change settings")
	} else {
		c.params.setHCL(p)
		return nil
	}
}

// :NODOC: This is mainly for tests.
// Not stopping you from using it though.  Maybe inspection?
func (c *Client) JSON() map[string]any {
	if c == nil {
		return map[string]any(nil)
	} else {
		q, _ := c.ensurePopulated()
		w := maps.All(*q)
		e := rejectSeq2(w, func(k string, v any) bool { return k == "Profile" || !isJSONMarshalable(v) })
		r := maps.Collect(e)

		return r
	}
}

func (c *Client) Profile() (*Profile, error) {
	var p *Profile

	if q, err := c.ensurePopulated(); err != nil {
		return p, err
	} else if result := obtainFromConfig[*Profile](q, "Profile"); result.isErr() {
		return p, result.error()
	} else {
		return result.unwrapOr(nil), nil
	}
}

func (c *Client) ProfileName() (dir, name *string) {
	i, _ := c.__populate__() // This __populate__ can fail, but we go ahead.
	result := obtainFromConfig[string](&i.c, "ProfileName")
	op, err := c.ProfileOp()

	if str, ok := result.some(); ok {
		name = &str
	}

	if err == nil {
		dir = Ptr(op.Dir())
	}

	return
}

func (c *Client) ProfileOp() (ProfileAPI, error) {
	if c == nil {
		return nil, NewErrorf("nil client")
	} else if ret := c.params.profileOp; ret == nil {
		return nil, NewErrorf("ProfileOp uninitialized")
	} else {
		return ret, nil
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil {
		return nil, NewErrorf("nil client")
	} else if doer, err := c.ensureDoer(); err != nil {
		return nil, err
	} else {
		return doer.Do(req)
	}
}

func (c *Client) ensurePopulated() (*config, error) {
	if c == nil {
		return nil, NewErrorf("nil client")
	} else if i, err := c.__populate__(); err != nil {
		return nil, err
	} else {
		return &i.c, nil
	}
}

func (c *Client) ensureDoer() (HttpRequestDoer, error) {
	if c == nil {
		return nil, NewErrorf("nil client")
	} else if i, err := c.__populate__(); err != nil {
		return nil, err
	} else {
		return i.d, nil
	}
}

func (c *Client) __populate__() (*inner, error) {
	return c.once.Do(func(i *inner) error {
		i.c = make(config)

		if err := c.params.populate(&i.c); err != nil {
			return err
		} else if i.d, err = newHttpRequestDoer(&i.c); err != nil {
			return err
		} else {
			return nil
		}
	})
}

// :NODOC:
type inner struct {
	c config
	d HttpRequestDoer
}

func isJSONMarshalable(v any) bool {
	enc := json.NewEncoder(io.Discard)
	err := enc.Encode(v)

	return err == nil
}

var _ ClientAPI = (*Client)(nil)
var _ HttpRequestDoer = (*Client)(nil)
