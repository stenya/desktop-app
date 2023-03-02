//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2023 Privatus Limited.
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

package conntest

import (
	"fmt"
	"net"

	api_types "github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/service/preferences"
	service_types "github.com/ivpn/desktop-app/daemon/service/types"
)

type ConnectivityTester struct {
	servers         api_types.ServersInfoResponse
	session         preferences.SessionStatus
	connParams      service_types.ConnectionParams
	isStopRequested bool
}

func (ct *ConnectivityTester) Stop() error {
	ct.isStopRequested = true
	return nil
}

func (ct ConnectivityTester) Test(servers api_types.ServersInfoResponse, session preferences.SessionStatus, connParams service_types.ConnectionParams) error {
	ct.servers = servers
	ct.session = session

	if !ct.session.IsLoggedIn() {
		return fmt.Errorf("account not initialised (please, login)")
	}

	wgLocalIP := net.ParseIP(ct.session.WGLocalIP)
	if wgLocalIP == nil {
		return fmt.Errorf("error updating WG connection preferences (failed parsing local IP for WG connection)")
	}

	var (
		err error
		wct *ConnectivityTesterWireguard
	)
	stopNotifierChan := make(chan struct{})
	defer func() {
		close(stopNotifierChan)
		if wct != nil {
			wct.WaitForDisconnect()
		}
	}()

	wct, err = InitTesterWireguard(stopNotifierChan, wgLocalIP, session.WGPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to initialise WireGuard tester object: %w", err)
	}

	for _, svr := range ct.servers.WireguardServers {
		for _, host := range svr.Hosts {
			for _, port := range ct.servers.Config.Ports.WireGuard {
				if port.Port == 0 {
					continue
				}
				if ct.isStopRequested {
					return nil
				}
				err := wct.Test(host, port.Port)

				msg := fmt.Sprintf("SVR: %s (%s) Host=%s\t(%s:%d)", svr.Country, svr.City, host.Hostname, host.Host, port.Port)
				if err != nil {
					fmt.Printf("*** ERROR    : %s, err=%v\n", msg, err)
				} else {
					fmt.Printf("*** Connected: %s\n", msg)
				}
			}
		}
	}

	return nil
}
