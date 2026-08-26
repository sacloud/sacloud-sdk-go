// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package fake

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/structs"
	"github.com/mitchellh/go-homedir"
	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
)

const defaultJSONFilePath = "libsacloud-fake-store.json"

// JSONFileStore .
type JSONFileStore struct {
	Path       string
	Ctx        context.Context
	NoInitData bool

	mu    sync.Mutex
	cache JSONFileStoreData
}

// JSONFileStoreData .
type JSONFileStoreData map[string]map[string]any

// MarshalJSON .
func (d JSONFileStoreData) MarshalJSON() ([]byte, error) {
	var transformed []map[string]any
	for cacheKey, resources := range d {
		resourceKey, zone := d.parseKey(cacheKey)
		for id, value := range resources {
			var mapValue map[string]any
			if d.isArrayOrSlice(value) {
				mapValue = map[string]any{
					"Values": value,
				}
			} else {
				mapValue = structs.Map(value)
			}

			mapValue["ID"] = id
			mapValue["ZoneName"] = zone
			mapValue["ResourceType"] = resourceKey

			transformed = append(transformed, mapValue)
		}
	}

	sort.Slice(transformed, func(i, j int) bool {
		rt1 := transformed[i]["ResourceType"].(string)
		rt2 := transformed[j]["ResourceType"].(string)
		if rt1 == rt2 {
			id1 := types.StringID(transformed[i]["ID"].(string))
			id2 := types.StringID(transformed[j]["ID"].(string))
			return id1 < id2
		}
		return rt1 < rt2
	})

	return json.MarshalIndent(transformed, "", "\t")
}

// UnmarshalJSON .
func (d *JSONFileStoreData) UnmarshalJSON(data []byte) error {
	var transformed []map[string]any
	if err := json.Unmarshal(data, &transformed); err != nil {
		return err
	}

	dest := JSONFileStoreData{}
	for _, mapValue := range transformed {
		rawID, ok := mapValue["ID"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ID' field is missing: %v", mapValue)
		}
		id := rawID.(string)

		rawZone, ok := mapValue["ZoneName"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ZoneName' field is missing: %v", mapValue)
		}
		zone := rawZone.(string)

		rawRt, ok := mapValue["ResourceType"]
		if !ok {
			return fmt.Errorf("invalid JSON: 'ResourceType' field is missing: %v", mapValue)
		}
		rt := rawRt.(string)

		var resources map[string]any
		r, ok := dest[d.key(rt, zone)]
		if ok {
			resources = r
		} else {
			resources = map[string]any{}
		}
		if v, ok := mapValue["Values"]; ok {
			resources[id] = v
		} else {
			resources[id] = mapValue
		}

		dest[d.key(rt, zone)] = resources
	}

	*d = dest
	return nil
}

func (d *JSONFileStoreData) isArrayOrSlice(v any) bool {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	case reflect.Pointer:
		return d.isArrayOrSlice(reflect.ValueOf(v).Elem().Interface())
	}
	return false
}

func (d *JSONFileStoreData) key(resourceKey, zone string) string {
	return fmt.Sprintf("%s/%s", resourceKey, zone)
}

func (d *JSONFileStoreData) parseKey(k string) (string, string) {
	ss := strings.Split(k, "/")
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	return "", ""
}

// NewJSONFileStore .
func NewJSONFileStore(path string) *JSONFileStore {
	return &JSONFileStore{
		Path:  path,
		cache: make(map[string]map[string]any),
	}
}

// Init .
func (s *JSONFileStore) Init() error {
	if s.Ctx == nil {
		s.Ctx = context.Background()
	}
	if s.Path == "" {
		s.Path = defaultJSONFilePath
	}

	// expand filepath
	path, err := homedir.Expand(s.Path)
	if err != nil {
		return err
	}
	s.Path = path

	if stat, err := os.Stat(s.Path); err == nil {
		if stat.IsDir() {
			return fmt.Errorf("path %q is directory", s.Path)
		}
	} else {
		if _, err := os.Create(s.Path); err != nil {
			return err
		}
	}

	if err := s.load(); err != nil {
		return err
	}
	s.startWatcher()
	return nil
}

// NeedInitData .
func (s *JSONFileStore) NeedInitData() bool {
	if s.NoInitData {
		return false
	}
	return len(s.cache) < 2
}

// Put .
func (s *JSONFileStore) Put(resourceKey, zone string, id types.ID, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values == nil {
		values = map[string]any{}
	}
	values[id.String()] = value
	s.cache[s.key(resourceKey, zone)] = values

	s.store() //nolint
}

// Get .
func (s *JSONFileStore) Get(resourceKey, zone string, id types.ID) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values == nil {
		return nil
	}
	return values[id.String()]
}

// List .
func (s *JSONFileStore) List(resourceKey, zone string) []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	var ret []any
	for _, v := range values {
		ret = append(ret, v)
	}
	return ret
}

// Delete .
func (s *JSONFileStore) Delete(resourceKey, zone string, id types.ID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	values := s.values(resourceKey, zone)
	if values != nil {
		delete(values, id.String())
	}
	s.store() //nolint
}

var jsonResourceTypeMap = map[string]func() any{
	ResourceArchive:           func() any { return &iaas.Archive{} },
	ResourceAuthStatus:        func() any { return &iaas.AuthStatus{} },
	ResourceAutoBackup:        func() any { return &iaas.AutoBackup{} },
	ResourceBill:              func() any { return &iaas.Bill{} },
	ResourceBridge:            func() any { return &iaas.Bridge{} },
	ResourceCDROM:             func() any { return &iaas.CDROM{} },
	ResourceContainerRegistry: func() any { return &iaas.ContainerRegistry{} },
	ResourceCoupon:            func() any { return &iaas.Coupon{} },
	ResourceDatabase:          func() any { return &iaas.Database{} },
	ResourceDisk:              func() any { return &iaas.Disk{} },
	ResourceDiskPlan:          func() any { return &iaas.DiskPlan{} },
	ResourceDNS:               func() any { return &iaas.DNS{} },
	ResourceEnhancedDB:        func() any { return &iaas.EnhancedDB{} },
	ResourceESME:              func() any { return &iaas.ESME{} },
	ResourceGSLB:              func() any { return &iaas.GSLB{} },
	ResourceIcon:              func() any { return &iaas.Icon{} },
	ResourceInterface:         func() any { return &iaas.Interface{} },
	ResourceInternet:          func() any { return &iaas.Internet{} },
	ResourceInternetPlan:      func() any { return &iaas.InternetPlan{} },
	ResourceIPAddress:         func() any { return &iaas.IPAddress{} },
	ResourceIPv6Net:           func() any { return &iaas.IPv6Net{} },
	ResourceIPv6Addr:          func() any { return &iaas.IPv6Addr{} },
	ResourceLicense:           func() any { return &iaas.License{} },
	ResourceLicenseInfo:       func() any { return &iaas.LicenseInfo{} },
	ResourceLoadBalancer:      func() any { return &iaas.LoadBalancer{} },
	ResourceLocalRouter:       func() any { return &iaas.LocalRouter{} },
	ResourceMobileGateway:     func() any { return &iaas.MobileGateway{} },
	ResourceNFS:               func() any { return &iaas.NFS{} },
	ResourceNote:              func() any { return &iaas.Note{} },
	ResourcePacketFilter:      func() any { return &iaas.PacketFilter{} },
	ResourcePrivateHost:       func() any { return &iaas.PrivateHost{} },
	ResourcePrivateHostPlan:   func() any { return &iaas.PrivateHostPlan{} },
	ResourceProxyLB:           func() any { return &iaas.ProxyLB{} },
	ResourceRegion:            func() any { return &iaas.Region{} },
	ResourceServer:            func() any { return &iaas.Server{} },
	ResourceServerPlan:        func() any { return &iaas.ServerPlan{} },
	ResourceServiceClass:      func() any { return &iaas.ServiceClass{} },
	ResourceSIM:               func() any { return &iaas.SIM{} },
	ResourceSimpleMonitor:     func() any { return &iaas.SimpleMonitor{} },
	ResourceSubnet:            func() any { return &iaas.Subnet{} },
	ResourceSSHKey:            func() any { return &iaas.SSHKey{} },
	ResourceSwitch:            func() any { return &iaas.Switch{} },
	ResourceVPCRouter:         func() any { return &iaas.VPCRouter{} },
	ResourceZone:              func() any { return &iaas.Zone{} },

	valuePoolResourceKey:         func() any { return &valuePool{} },
	"BillDetails":                func() any { return &[]*iaas.BillDetail{} },
	"ContainerRegistryUsers":     func() any { return &[]*iaas.ContainerRegistryUser{} },
	"DatabaseParameter":          func() any { return map[string]any{} },
	"ESMELogs":                   func() any { return &[]*iaas.ESMELogs{} },
	"LocalRouterStatus":          func() any { return &iaas.LocalRouterHealth{} },
	"MobileGatewayDNS":           func() any { return &iaas.MobileGatewayDNSSetting{} },
	"MobileGatewaySIMRoutes":     func() any { return &[]*iaas.MobileGatewaySIMRoute{} },
	"MobileGatewaySIMs":          func() any { return &[]*iaas.MobileGatewaySIMInfo{} },
	"MobileGatewayTrafficConfig": func() any { return &iaas.MobileGatewayTrafficControl{} },
	"ProxyLBStatus":              func() any { return &iaas.ProxyLBHealth{} },
	"SIMNetworkOperator":         func() any { return &[]*iaas.SIMNetworkOperatorConfig{} },
}

func (s *JSONFileStore) unmarshalResource(resourceKey string, data []byte) (any, error) {
	f, ok := jsonResourceTypeMap[resourceKey]
	if !ok {
		panic(fmt.Errorf("type %q is not registered", resourceKey))
	}
	v := f()
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *JSONFileStore) store() error {
	data, err := json.MarshalIndent(s.cache, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0600)
}

func (s *JSONFileStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	var cache = JSONFileStoreData{}
	if err := json.Unmarshal(data, &cache); err != nil {
		return err
	}

	var loaded = make(map[string]map[string]any)
	for cacheKey, values := range cache {
		resourceKey, _ := s.parseKey(cacheKey)

		var dest = make(map[string]any)
		for id, v := range values {
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			cv, err := s.unmarshalResource(resourceKey, data)
			if err != nil {
				return err
			}
			dest[id] = cv
		}
		loaded[cacheKey] = dest
	}
	s.cache = loaded
	return nil
}

func (s *JSONFileStore) key(resourceKey, zone string) string {
	return fmt.Sprintf("%s/%s", resourceKey, zone)
}

func (s *JSONFileStore) parseKey(k string) (string, string) {
	ss := strings.Split(k, "/")
	if len(ss) == 2 {
		return ss[0], ss[1]
	}
	return "", ""
}

func (s *JSONFileStore) values(resourceKey, zone string) map[string]any {
	return s.cache[s.key(resourceKey, zone)]
}
