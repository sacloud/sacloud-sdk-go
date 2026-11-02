// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package database

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	iaasapi "github.com/sacloud/sacloud-sdk-go/api/iaas/helper/api"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/helper/power"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/helper/wait"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/testutil"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	"github.com/sacloud/sacloud-sdk-go/common/packages/size"
	builder2 "github.com/sacloud/sacloud-sdk-go/service/iaas/database/builder"
	diskbuilder "github.com/sacloud/sacloud-sdk-go/service/iaas/disk/builder"
	serverbuilder "github.com/sacloud/sacloud-sdk-go/service/iaas/server/builder"
	"github.com/sacloud/sacloud-sdk-go/service/iaas/setup"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestDatabaseService_UpdateWithContext(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("This test runs only with TESTACC=1")
	}

	passwordBytes := make([]byte, 15)
	if _, err := rand.Read(passwordBytes); err != nil {
		t.Fatal(err)
	}
	newPassword := hex.EncodeToString(passwordBytes)
	replicaPasswordBytes := make([]byte, 15)
	if _, err := rand.Read(replicaPasswordBytes); err != nil {
		t.Fatal(err)
	}
	newReplicaPassword := hex.EncodeToString(replicaPasswordBytes)

	zone := databasePasswordTestZone()
	if zone == "tk1v" {
		t.Fatal("database password acceptance test requires a bootable zone; set SAKURACLOUD_ZONE (for example, is1b)")
	}
	caller, err := iaasapi.NewCaller()
	require.NoError(t, err)

	name := testutil.ResourceName("database-update-password")
	clientServerName := testutil.ResourceName("database-update-password-client")
	updateName := name + "-upd"
	var switchID types.ID
	var testClient *databasePasswordTestClient
	svc := New(caller)
	weekdays := make([]types.EDayOfTheWeek, 0)
	initialPassword := "initial-password"
	initialReplicaPassword := "initial-replica-password"

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: func() iaas.APICaller { return caller },
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			if err := testutil.ComposeCleanupResourceFunc(
				clientServerName,
				testutil.CleanupTargets.Server,
				testutil.CleanupTargets.Disk,
			)(ctx, caller); err != nil {
				return err
			}
			if err := testutil.ComposeCleanupResourceFunc(
				name,
				testutil.CleanupTargets.Database,
				testutil.CleanupTargets.Switch,
			)(ctx, caller); err != nil {
				return err
			}
			sw, err := iaas.NewSwitchOp(caller).Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
			if err != nil {
				return err
			}
			switchID = sw.ID
			testClient, err = newDatabasePasswordTestClient(ctx, caller, zone, switchID, clientServerName)
			if err != nil {
				return err
			}
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (any, error) {
				return (&builder2.Builder{
					Zone:           zone,
					PlanID:         types.DatabasePlans.DB10GB,
					SwitchID:       switchID,
					IPAddresses:    []string{"192.168.0.101"},
					NetworkMaskLen: 24,
					DefaultRoute:   "192.168.0.1",
					Conf: &iaas.DatabaseRemarkDBConfCommon{
						DatabaseName: types.RDBMSTypesPostgreSQL.String(),
						DefaultUser:  "default",
						UserPassword: initialPassword,
					},
					CommonSetting: &iaas.DatabaseSettingCommon{
						ServicePort:     5432,
						SourceNetwork:   []string{"192.168.0.0/24"},
						DefaultUser:     "default",
						UserPassword:    initialPassword,
						ReplicaUser:     "replica",
						ReplicaPassword: "initial-replica-password",
					},
					ReplicationSetting: &iaas.DatabaseReplicationSetting{
						Model: types.DatabaseReplicationModels.MasterSlave,
					},
					Name:   name,
					Client: builder2.NewAPIClient(caller),
				}).Build(ctx)
			},
			CheckFunc: func(_ testutil.TestT, _ *testutil.CRUDTestContext, value any) error {
				db := value.(*iaas.Database)
				if err := testClient.waitForPostgreSQLClient(); err != nil {
					return err
				}
				if err := testClient.assertPassword(db.IPAddresses[0], initialPassword, true); err != nil {
					return err
				}
				return testClient.assertReplicaPassword(db.IPAddresses[0], initialReplicaPassword, true)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (any, error) {
				return iaas.NewDatabaseOp(caller).Read(ctx, zone, ctx.ID)
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
					t.Log("updating database without setting Password or ReplicaUserPassword")
					db, err := svc.UpdateWithContext(ctx, &UpdateRequest{
						Zone:             zone,
						ID:               ctx.ID,
						Name:             new(updateName),
						Backupv2Weekdays: &weekdays,
						NoWait:           true,
					})
					if err != nil {
						return nil, err
					}
					host := db.IPAddresses[0]
					if err := testClient.assertPassword(host, initialPassword, true); err != nil {
						return nil, err
					}
					if err := testClient.assertPassword(host, newPassword, false); err != nil {
						return nil, err
					}
					if err := testClient.assertReplicaPassword(host, initialReplicaPassword, true); err != nil {
						return nil, err
					}
					if err := testClient.assertReplicaPassword(host, newReplicaPassword, false); err != nil {
						return nil, err
					}
					return db, nil
				},
				CheckFunc: func(t testutil.TestT, _ *testutil.CRUDTestContext, value any) error {
					db := value.(*iaas.Database)
					return testutil.AssertEqual(t, updateName, db.Name, "Database.Name")
				},
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
					t.Log("updating database with new Password and ReplicaUserPassword")
					db, err := svc.UpdateWithContext(ctx, &UpdateRequest{
						Zone:                zone,
						ID:                  ctx.ID,
						Password:            &newPassword,
						ReplicaUserPassword: &newReplicaPassword,
						Backupv2Weekdays:    &weekdays,
					})
					if err != nil {
						return nil, err
					}
					host := db.IPAddresses[0]
					if err := testClient.waitForPassword(host, newPassword); err != nil {
						return nil, err
					}
					if err := testClient.assertPassword(host, initialPassword, false); err != nil {
						return nil, err
					}
					if err := testClient.waitForReplicaPassword(host, newReplicaPassword); err != nil {
						return nil, err
					}
					if err := testClient.assertReplicaPassword(host, initialReplicaPassword, false); err != nil {
						return nil, err
					}
					return db, nil
				},
				CheckFunc: func(t testutil.TestT, _ *testutil.CRUDTestContext, value any) error {
					db := value.(*iaas.Database)
					return testutil.AssertEqual(t, updateName, db.Name, "Database.Name")
				},
			},
		},
		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			return power.ShutdownDatabase(ctx, iaas.NewDatabaseOp(caller), zone, ctx.ID, true)
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				return iaas.NewDatabaseOp(caller).Delete(ctx, zone, ctx.ID)
			},
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			var cleanupErrors []error
			if err := testClient.close(); err != nil {
				cleanupErrors = append(cleanupErrors, err)
			}
			if err := testutil.ComposeCleanupResourceFunc(
				clientServerName,
				testutil.CleanupTargets.Server,
				testutil.CleanupTargets.Disk,
			)(ctx, caller); err != nil {
				cleanupErrors = append(cleanupErrors, err)
			}
			if err := testutil.ComposeCleanupResourceFunc(
				name,
				testutil.CleanupTargets.Database,
				testutil.CleanupTargets.Switch,
			)(ctx, caller); err != nil {
				cleanupErrors = append(cleanupErrors, err)
			}
			return errors.Join(cleanupErrors...)
		},
	})
}

const (
	databaseTestNetworkMask = 24
	databaseTestClientIP    = "192.168.0.102"
)

func databasePasswordTestZone() string {
	if zone := os.Getenv("SAKURACLOUD_ZONE"); zone != "" {
		return zone
	}
	return "is1b"
}

type databasePasswordTestClient struct {
	client *ssh.Client
}

func newDatabasePasswordTestClient(ctx context.Context, caller iaas.APICaller, zone string, switchID types.ID, serverName string) (*databasePasswordTestClient, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	archive, err := findUbuntuCloudimgArchive(ctx, caller, zone)
	if err != nil {
		return nil, err
	}
	diskSizeGB := (archive.SizeMB + size.GiB - 1) / size.GiB

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, err
	}

	buildResult, err := (&serverbuilder.Builder{
		Name:            serverName,
		CPU:             1,
		MemoryGB:        1,
		BootAfterCreate: true,
		NIC:             &serverbuilder.SharedNICSetting{},
		AdditionalNICs: []serverbuilder.AdditionalNICSettingHolder{
			&serverbuilder.ConnectedNICSetting{SwitchID: switchID},
		},
		DiskBuilders: []diskbuilder.Builder{
			&diskbuilder.FromDiskOrArchiveBuilder{
				SourceArchiveID: archive.ID,
				Name:            serverName + "-disk",
				SizeGB:          diskSizeGB,
				PlanID:          types.DiskPlans.SSD,
				Connection:      types.DiskConnections.VirtIO,
				Client:          diskbuilder.NewBuildersAPIClient(caller),
			},
		},
		UserData: fmt.Sprintf(`#cloud-config
package_update: true
packages:
  - openssh-server
  - postgresql-client
ssh_pwauth: false
disable_root: true
users:
  - default
  - name: database-test
    shell: /bin/bash
    lock_passwd: true
    sudo: ALL=(ALL) NOPASSWD:ALL
    ssh_authorized_keys:
      - %s
runcmd:
  - [systemctl, enable, --now, ssh]
`, strings.TrimSpace(string(ssh.MarshalAuthorizedKey(signer.PublicKey())))),
		Client: serverbuilder.NewBuildersAPIClient(caller),
	}).Build(ctx, zone)
	if err != nil {
		return nil, err
	}

	server, err := iaas.NewServerOp(caller).Read(ctx, zone, buildResult.ServerID)
	if err != nil {
		return nil, err
	}

	publicIP := ""
	privateMAC := ""
	for _, iface := range server.Interfaces {
		if iface.SwitchID == switchID {
			privateMAC = iface.MACAddress
		}
		if iface.SwitchScope == types.Scopes.Shared {
			publicIP = iface.IPAddress
		}
	}
	if net.ParseIP(publicIP) == nil {
		return nil, fmt.Errorf("test server public IP is invalid: %q", publicIP)
	}
	if _, err := net.ParseMAC(privateMAC); err != nil {
		return nil, fmt.Errorf("test server private NIC MAC is invalid: %w", err)
	}

	client, err := connectDatabaseTestServer(ctx, publicIP, signer)
	if err != nil {
		return nil, err
	}
	linkOutput, err := runDatabaseTestServerCommand(client, "ip -o link", "")
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("reading test server interfaces failed: %w", err)
	}
	privateInterface, err := interfaceNameForMAC(linkOutput, privateMAC)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	if !isSafeInterfaceName(privateInterface) {
		_ = client.Close()
		return nil, fmt.Errorf("invalid test server interface name: %q", privateInterface)
	}
	command := fmt.Sprintf(
		"sudo -n ip link set %s up && sudo -n ip address add %s/%d dev %s",
		privateInterface,
		databaseTestClientIP,
		databaseTestNetworkMask,
		privateInterface,
	)
	if _, err := runDatabaseTestServerCommand(client, command, ""); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("configuring test server private NIC failed: %w", err)
	}

	return &databasePasswordTestClient{client: client}, nil
}

func findUbuntuCloudimgArchive(ctx context.Context, caller iaas.APICaller, zone string) (*iaas.Archive, error) {
	const pageSize = 100

	archiveOp := iaas.NewArchiveOp(caller)
	var archives []*iaas.Archive
	for from := 0; ; {
		result, err := archiveOp.Find(ctx, zone, &iaas.FindCondition{
			Count: pageSize,
			From:  from,
		})
		if err != nil {
			return nil, err
		}
		archives = append(archives, result.Archives...)
		if len(result.Archives) == 0 {
			break
		}
		from += len(result.Archives)
		if from >= result.Total {
			break
		}
	}

	return selectUbuntuCloudimgArchive(archives)
}

func selectUbuntuCloudimgArchive(archives []*iaas.Archive) (*iaas.Archive, error) {
	var latest *iaas.Archive
	for _, archive := range archives {
		if archive == nil ||
			archive.Availability != types.Availabilities.Available ||
			archive.Scope != types.Scopes.Shared {
			continue
		}
		name := strings.ToLower(archive.Name)
		if !strings.Contains(name, "ubuntu") || !strings.Contains(name, "(cloudimg)") {
			continue
		}
		if latest == nil || archive.ID > latest.ID {
			latest = archive
		}
	}
	if latest == nil {
		return nil, errors.New("available Ubuntu cloudimg archive not found")
	}
	return latest, nil
}

func TestSelectUbuntuCloudimgArchive(t *testing.T) {
	archives := []*iaas.Archive{
		{
			ID:           1,
			Name:         "Ubuntu Server 24.04 LTS 64bit",
			Availability: types.Availabilities.Available,
			Scope:        types.Scopes.Shared,
		},
		{
			ID:           2,
			Name:         "Ubuntu Server 24.04 LTS 64bit (cloudimg)",
			Availability: types.Availabilities.Available,
			Scope:        types.Scopes.Shared,
		},
		{
			ID:           3,
			Name:         "Ubuntu Server 26.04 LTS 64bit (cloudimg)",
			Availability: types.Availabilities.Available,
			Scope:        types.Scopes.Shared,
		},
		{
			ID:           4,
			Name:         "Debian 13 64bit (cloudimg)",
			Availability: types.Availabilities.Available,
			Scope:        types.Scopes.Shared,
		},
	}

	selected, err := selectUbuntuCloudimgArchive(archives)
	require.NoError(t, err)
	require.Equal(t, types.ID(3), selected.ID)
}

func TestSelectUbuntuCloudimgArchiveMissing(t *testing.T) {
	_, err := selectUbuntuCloudimgArchive([]*iaas.Archive{
		{
			ID:           1,
			Name:         "Ubuntu Server 26.04 LTS 64bit",
			Availability: types.Availabilities.Available,
			Scope:        types.Scopes.Shared,
		},
	})
	require.ErrorContains(t, err, "available Ubuntu cloudimg archive not found")
}

func connectDatabaseTestServer(ctx context.Context, ip string, signer ssh.Signer) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: "database-test",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         10 * time.Second,
	}
	address := net.JoinHostPort(ip, "22")
	deadline := time.Now().Add(5 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		client, err := ssh.Dial("tcp", address, config)
		if err == nil {
			return client, nil
		}
		lastErr = err
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	return nil, fmt.Errorf("connecting to test server over SSH timed out: %w", lastErr)
}

func runDatabaseTestServerCommand(client *ssh.Client, command string, input string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close() //nolint:errcheck
	if input != "" {
		session.Stdin = strings.NewReader(input)
	}
	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("%w: %s", err, output)
	}
	return string(output), nil
}

func interfaceNameForMAC(output, mac string) (string, error) {
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		for i, field := range fields {
			if field == "link/ether" && i+1 < len(fields) && strings.EqualFold(fields[i+1], mac) {
				name := strings.TrimSuffix(fields[1], ":")
				return strings.Split(name, "@")[0], nil
			}
		}
	}
	return "", fmt.Errorf("private NIC %s not found on test server", mac)
}

func isSafeInterfaceName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !(r == '-' || r == '_' || r == '.' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func (c *databasePasswordTestClient) close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *databasePasswordTestClient) waitForPostgreSQLClient() error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		_, lastErr = runDatabaseTestServerCommand(c.client, "command -v psql", "")
		if lastErr == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("postgresql client was not installed on test server: %w", lastErr)
}

func (c *databasePasswordTestClient) runSQL(host, user, password, query string) (string, error) {
	if net.ParseIP(host) == nil {
		return "", fmt.Errorf("invalid database IP address: %q", host)
	}
	command := fmt.Sprintf(
		`PGPASSWORD="$(cat)" PGCONNECT_TIMEOUT=5 psql -w -v ON_ERROR_STOP=1 -h %s -p 5432 -U %s -d postgres -tAc %s`,
		shellQuote(host),
		shellQuote(user),
		shellQuote(query),
	)
	output, err := runDatabaseTestServerCommand(c.client, command, password)
	if err != nil {
		return output, err
	}
	return strings.TrimSpace(output), nil
}

func (c *databasePasswordTestClient) runReplicationCommand(host, password string) (string, error) {
	if net.ParseIP(host) == nil {
		return "", fmt.Errorf("invalid database IP address: %q", host)
	}
	conninfo := fmt.Sprintf("host=%s port=5432 user=replica replication=database", host)
	command := fmt.Sprintf(
		`PGPASSWORD="$(cat)" PGCONNECT_TIMEOUT=5 psql -w -d %s -tAc %s`,
		shellQuote(conninfo),
		shellQuote("IDENTIFY_SYSTEM"),
	)
	output, err := runDatabaseTestServerCommand(c.client, command, password)
	if err != nil {
		return output, err
	}
	return strings.TrimSpace(output), nil
}

func (c *databasePasswordTestClient) assertPassword(host, password string, shouldWork bool) error {
	output, err := c.runSQL(host, "default", password, "SELECT 1")
	if shouldWork {
		if err != nil {
			return fmt.Errorf("database authentication failed: %w", err)
		}
		if output != "1" {
			return fmt.Errorf("unexpected database query result: %q", output)
		}
		return nil
	}
	if err == nil {
		return errors.New("database authentication unexpectedly succeeded")
	}
	if !strings.Contains(strings.ToLower(output), "password authentication failed") {
		return fmt.Errorf("expected password authentication failure, got: %w", err)
	}
	return nil
}

func (c *databasePasswordTestClient) waitForPassword(host, password string) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		lastErr = c.assertPassword(host, password, true)
		if lastErr == nil {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("database did not accept the password before timeout: %w", lastErr)
}

func (c *databasePasswordTestClient) assertReplicaPassword(host, password string, shouldWork bool) error {
	output, err := c.runReplicationCommand(host, password)
	if shouldWork {
		if err != nil {
			return fmt.Errorf("replication authentication failed: %w", err)
		}
		if output == "" {
			return errors.New("replication command returned no system identifier")
		}
		return nil
	}
	if err == nil {
		return errors.New("replication authentication unexpectedly succeeded")
	}
	if !strings.Contains(strings.ToLower(output), "password authentication failed") {
		return fmt.Errorf("expected replication password authentication failure, got: %w", err)
	}
	return nil
}

func (c *databasePasswordTestClient) waitForReplicaPassword(host, password string) error {
	deadline := time.Now().Add(5 * time.Minute)
	var lastErr error
	for time.Now().Before(deadline) {
		lastErr = c.assertReplicaPassword(host, password, true)
		if lastErr == nil {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("replication password was not accepted before timeout: %w", lastErr)
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func TestDatabaseService_convertUpdateRequest(t *testing.T) {
	if testutil.IsAccTest() {
		t.Skip("This test runs only without TESTACC=1")
	}
	for _, tc := range []struct {
		databaseName string
		databaseType string
	}{
		{databaseName: "MariaDB", databaseType: "mariadb"},
		{databaseName: "PostgreSQL", databaseType: "postgres"},
		{databaseName: "postgres", databaseType: "postgres"},
	} {
		t.Run(tc.databaseName, func(t *testing.T) {
			testDatabaseServiceConvertUpdateRequest(t, tc.databaseName, tc.databaseType)
		})
	}
}

func testDatabaseServiceConvertUpdateRequest(t *testing.T, databaseName, databaseType string) {
	t.Helper()
	ctx := context.Background()
	zone := testutil.TestZone()
	name := testutil.ResourceName("database-service-update")
	caller := testutil.SingletonAPICaller()

	sw, err := iaas.NewSwitchOp(caller).Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}
	builder := &builder2.Builder{
		Zone:           zone,
		PlanID:         types.DatabasePlans.DB10GB,
		SwitchID:       sw.ID,
		IPAddresses:    []string{"192.168.0.101"},
		NetworkMaskLen: 24,
		DefaultRoute:   "192.168.0.1",
		Conf: &iaas.DatabaseRemarkDBConfCommon{
			DatabaseName: databaseName,
		},
		SourceID: 0,
		CommonSetting: &iaas.DatabaseSettingCommon{
			WebUI:           types.ToWebUI(true),
			ServicePort:     5432,
			SourceNetwork:   []string{"192.168.0.0/24"},
			DefaultUser:     "default",
			UserPassword:    "password1",
			ReplicaUser:     "replica",
			ReplicaPassword: "password2",
		},
		BackupSetting: &iaas.DatabaseSettingBackup{
			Rotate:    0,
			Time:      "10:00",
			DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
		},
		ReplicationSetting: &iaas.DatabaseReplicationSetting{
			Model: types.DatabaseReplicationModels.MasterSlave,
		},
		Name:         name,
		Description:  "description",
		Tags:         types.Tags{"tag1", "tag2"},
		SettingsHash: "",
		NoWait:       false,
		Parameters: map[string]any{
			"max_connections": float64(100),
		},
		SetupOptions: &setup.Options{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: time.Millisecond,
			PollingInterval:       time.Millisecond,
		},
		Client: builder2.NewAPIClient(caller),
	}
	db, err := builder.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dbOp := iaas.NewDatabaseOp(caller)
		if err := dbOp.Shutdown(ctx, zone, db.ID, &iaas.ShutdownOption{Force: true}); err != nil {
			return
		}
		if _, err := wait.UntilDatabaseIsDown(ctx, dbOp, zone, db.ID); err != nil {
			return
		}
		if err := dbOp.Delete(ctx, zone, db.ID); err != nil {
			return
		}
		iaas.NewSwitchOp(caller).Delete(ctx, zone, sw.ID) //nolint
	}()

	cases := []struct {
		in     *UpdateRequest
		expect *ApplyRequest
	}{
		{
			in: &UpdateRequest{
				Zone:   zone,
				ID:     db.ID,
				Name:   new(db.Name + "-upd"),
				NoWait: true,
			},
			expect: &ApplyRequest{
				Zone:                  zone,
				ID:                    db.ID,
				Name:                  db.Name + "-upd",
				Description:           "description",
				Tags:                  types.Tags{"tag1", "tag2"},
				PlanID:                types.DatabasePlans.DB10GB,
				SwitchID:              sw.ID,
				IPAddresses:           []string{"192.168.0.101"},
				NetworkMaskLen:        24,
				DefaultRoute:          "192.168.0.1",
				Port:                  5432,
				SourceNetwork:         []string{"192.168.0.0/24"},
				DatabaseType:          databaseType,
				Username:              "default",
				Password:              "",
				EnableReplication:     true,
				ReplicaUserPassword:   "",
				EnableWebUI:           true,
				EnableBackup:          true,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]any{
					"max_connections": float64(100),
				},
				NoWait: true,
			},
		},
		{
			in: &UpdateRequest{
				Zone:              zone,
				ID:                db.ID,
				Password:          new("new-password"),
				EnableReplication: new(false),
				EnableBackup:      new(false),
				Parameters: &map[string]any{
					"work_mem": float64(4096),
				},
				NoWait: true,
			},
			expect: &ApplyRequest{
				Zone:                  zone,
				ID:                    db.ID,
				Name:                  db.Name,
				Description:           "description",
				Tags:                  types.Tags{"tag1", "tag2"},
				PlanID:                types.DatabasePlans.DB10GB,
				SwitchID:              sw.ID,
				IPAddresses:           []string{"192.168.0.101"},
				NetworkMaskLen:        24,
				DefaultRoute:          "192.168.0.1",
				Port:                  5432,
				SourceNetwork:         []string{"192.168.0.0/24"},
				DatabaseType:          databaseType,
				Username:              "default",
				Password:              "new-password",
				EnableReplication:     false,
				ReplicaUserPassword:   "",
				EnableWebUI:           true,
				EnableBackup:          false,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]any{
					"max_connections": float64(100),
					"work_mem":        float64(4096),
				},
				NoWait: true,
			},
		},
	}

	for _, tc := range cases {
		got, err := tc.in.ApplyRequest(ctx, caller)
		require.NoError(t, err)
		require.NoError(t, got.Validate())

		// データベースバージョンは比較対象としない
		got.DatabaseVersion = ""
		require.EqualValues(t, tc.expect, got)
	}

	newReplicaPassword := "new-replica-password"
	got, err := (&UpdateRequest{
		Zone:                zone,
		ID:                  db.ID,
		EnableReplication:   new(true),
		ReplicaUserPassword: &newReplicaPassword,
	}).ApplyRequest(ctx, caller)
	require.NoError(t, err)
	require.Equal(t, newReplicaPassword, got.ReplicaUserPassword)
}
