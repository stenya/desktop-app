//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2020 Privatus Limited.
//
//  This file is part of the Daemon for IVPN Client Desktop.
//
//  The Daemon for IVPN Client Desktop is free software: you can redistribute it and/or
//  modify it under the terms of the GNU General Public License as published by the Free
//  Software Foundation, either version 3 of the License, or (at your option) any later version.
//
//  The Daemon for IVPN Client Desktop is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
//  or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
//  details.
//
//  You should have received a copy of the GNU General Public License
//  along with the Daemon for IVPN Client Desktop. If not, see <https://www.gnu.org/licenses/>.
//

package platform

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ivpn/desktop-app/daemon/helpers"
)

var (
	firewallScript string
	splitTunScript string
	logDir         string = "/opt/ivpn/log"
	tmpDir         string = "/opt/ivpn/mutable"

	// path to the readonly servers.json file bundled into the package
	serversFileBundled string
)

// SnapEnvInfo contains values of SNAP environment variables
// (applicable only if running in SNAP)
// https://snapcraft.io/docs/environment-variables
type SnapEnvInfo struct {
	// Directory where the snap is mounted. This is where all the files in your snap are visible in the filesystem.
	// All of the data in the snap is read-only and cannot be changed.
	SNAP string
	// Directory for system data that is common across revisions of a snap.
	// This directory is owned and writable by root and is meant to be used by background applications (daemons, services).
	// Unlike SNAP_DATA this directory is not backed up and restored across snap refresh and revert operations.
	SNAP_COMMON string
	// Directory for system data of a snap.
	// This directory is owned and writable by root and is meant to be used by background applications (daemons, services).
	// Unlike SNAP_COMMON this directory is backed up and restored across snap refresh and snap revert operations.
	SNAP_DATA string
}

// GetSnapEnvs returns SNAP environment variables (or nil)
func GetSnapEnvs() *SnapEnvInfo {
	snap := os.Getenv("SNAP")
	snapCommon := os.Getenv("SNAP_COMMON")
	snapData := os.Getenv("SNAP_DATA")
	if len(snap) == 0 || len(snapCommon) == 0 || len(snapData) == 0 {
		return nil
	}
	if ex, err := os.Executable(); err == nil && len(ex) > 0 {
		if !strings.HasPrefix(ex, snap) {
			// if snap environment - the binary must be located in "$SNAP"
			return nil
		}
	}

	return &SnapEnvInfo{
		SNAP:        snap,
		SNAP_COMMON: snapCommon,
		SNAP_DATA:   snapData,
	}
}

// initialize all constant values (e.g. servicePortFile) which can be used in external projects (IVPN CLI)
func doInitConstants() {
	openVpnBinaryPath = "/usr/sbin/openvpn"
	routeCommand = "/sbin/ip route"

	// check if we are running in snap environment
	if envs := GetSnapEnvs(); envs != nil {
		logDir = path.Join(envs.SNAP_COMMON, "/opt/ivpn/log")
		tmpDir = path.Join(envs.SNAP_COMMON, "/opt/ivpn/mutable")
		openVpnBinaryPath = path.Join(envs.SNAP, openVpnBinaryPath)
	}

	fwInitialValueAllowApiServers = false
	serversFile = path.Join(tmpDir, "servers.json")
	servicePortFile = path.Join(tmpDir, "port.txt")
	paranoidModeSecretFile = path.Join(tmpDir, "eaa")

	logFile = path.Join(logDir, "IVPN_Agent.log")
	openvpnLogFile = path.Join(logDir, "openvpn.log")

	openvpnUserParamsFile = path.Join(tmpDir, "ovpn_extra_params.txt")
}

func doOsInit() (warnings []string, errors []error, logInfo []string) {
	warnings, errors, logInfo = doOsInitForBuild()

	if errors == nil {
		errors = make([]error, 0)
	}

	if err := checkFileAccessRightsExecutable("firewallScript", firewallScript); err != nil {
		errors = append(errors, err)
	}
	if err := checkFileAccessRightsExecutable("splitTunScript", splitTunScript); err != nil {
		errors = append(errors, err)
	}

	return warnings, errors, logInfo
}

//func doInitOperations() (w string, e error) { return "", nil }
func doInitOperations() (w string, e error) {
	serversFile := ServersFile()
	if _, err := os.Stat(serversFile); err != nil {
		if os.IsNotExist(err) {
			if len(serversFileBundled) == 0 {
				return fmt.Sprintf("'%s' not exists and the 'serversFileBundled' path not defined", serversFile), nil
			}

			srcStat, err := os.Stat(serversFileBundled)
			if err != nil {
				return fmt.Sprintf("'%s' not exists and the serversFileBundled='%s' access error: %s", serversFile, serversFileBundled, err.Error()), nil
			}

			fmt.Printf("File '%s' does not exists. Copying from bundle (%s)...\n", serversFile, serversFileBundled)
			// Servers file is not exists on required place
			// Probably, it is first start after clean install
			// Copying it from a bundle
			os.MkdirAll(filepath.Base(serversFile), os.ModePerm)
			if err = helpers.CopyFile(serversFileBundled, serversFile); err != nil {
				return err.Error(), nil
			}

			// keep file mode same as source file
			err = os.Chmod(serversFile, srcStat.Mode())
			if err != nil {
				return err.Error(), nil
			}

			return "", nil
		}

		return err.Error(), nil
	}
	return "", nil
}

// FirewallScript returns path to firewal script
func FirewallScript() string {
	return firewallScript
}

// SplitTunScript returns path to script which control split-tunneling functionality
func SplitTunScript() string {
	return splitTunScript
}
