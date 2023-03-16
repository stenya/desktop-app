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
	"sort"
	"time"

	api_types "github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/helpers"
	"github.com/ivpn/desktop-app/daemon/logger"
	"github.com/ivpn/desktop-app/daemon/service/preferences"
	service_types "github.com/ivpn/desktop-app/daemon/service/types"
	"github.com/ivpn/desktop-app/daemon/vpn"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("ctest")
}

const constTimeout time.Duration = time.Millisecond * 300

type ConnectivityTester interface {
	Stop() error
	Test(currentConnParams service_types.ConnectionParams, statusNotifyChan chan<- StatusEvent) (
		*GoodConnectionInfo, error)
	TestPorts(customPorts []service_types.PortData, getGeolookup func(timeoutMs int) (*api_types.GeoLookupResponse, error)) (
		testedPortsResult map[service_types.PortData]bool, // map[<tested_port>]<is_accessible>
		err error)
}

type connectivityTester struct {
	servers          api_types.ServersInfoResponse
	session          preferences.SessionStatus
	connParams       service_types.ConnectionParams
	statusNotifyChan chan<- StatusEvent
	isStopRequested  bool
}

type GoodConnectionInfo struct {
	Gateway  string // Server gateway
	HostName string // Host name (empty if all server hosts are OK)
	Port     int    // Port number
	PortType string // udp/tcp
}

type StatusEvent struct {
	Server api_types.ServerInfoBase
	Host   api_types.HostInfoBase
	Port   api_types.PortInfo
}

func CreateConnectivityTester(
	servers api_types.ServersInfoResponse,
	connParams service_types.ConnectionParams,
	session preferences.SessionStatus) (ConnectivityTester, error) {

	ret := &connectivityTester{
		servers:    servers,
		connParams: connParams,
		session:    session}

	return ret, nil
}

func (ct *connectivityTester) Stop() error {
	ct.isStopRequested = true
	return nil
}

func (ct connectivityTester) Test(
	connParams service_types.ConnectionParams,
	statusNotifyChan chan<- StatusEvent) (*GoodConnectionInfo, error) {

	if len(ct.servers.OpenvpnServers) == 0 && len(ct.servers.WireguardServers) == 0 {
		return nil, fmt.Errorf("internal error: object not initialised")
	}

	ct.connParams = connParams
	ct.statusNotifyChan = statusNotifyChan

	defer close(ct.statusNotifyChan)

	if !ct.session.IsLoggedIn() {
		return nil, fmt.Errorf("account not initialised (please, login)")
	}

	wgLocalIP := net.ParseIP(ct.session.WGLocalIP)
	if wgLocalIP == nil {
		return nil, fmt.Errorf("error updating WG connection preferences (failed parsing local IP for WG connection)")
	}

	wct, err := InitTesterWireguard(wgLocalIP, ct.session.WGPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise WireGuard tester object: %w", err)
	}
	defer wct.Disconnect()

	// sorting servers by distance to currently serlected server
	svrs := ct.sortServersByDistance(ct.servers.WireguardServers)

	// sort ports (e.g. the default port (selected by user) has highest priority and must be checked first)
	ports := ct.sortPorts(ct.servers.Config.Ports.WireGuard)

	for _, svr := range svrs {
		for _, host := range svr.Hosts {
			for _, port := range ports {
				if port.Port == 0 {
					continue
				}
				if ct.isStopRequested {
					return nil, fmt.Errorf("cancelled")
				}

				// notify current status
				select {
				case ct.statusNotifyChan <- StatusEvent{
					Server: svr.ServerInfoBase,
					Host:   host.HostInfoBase,
					Port:   port}:
				default: // channel is full
				}

				err := wct.Test(host, port.Port)

				if err == nil {
					return &GoodConnectionInfo{
						Gateway:  svr.Gateway,   // Server gateway
						HostName: host.Hostname, // Host name (empty if all server hosts are OK)
						Port:     port.Port,     // Port number
						PortType: port.Type,     // udp/tcp
					}, nil

				}
			}
		}
		break
	}

	return nil, fmt.Errorf("no good connection parameters found")
}

func (ct connectivityTester) sortServersByDistance(svrs []api_types.WireGuardServerInfo) (ret []api_types.WireGuardServerInfo) {
	ret = make([]api_types.WireGuardServerInfo, len(svrs))
	copy(ret, svrs) // do not change original slice

	if ct.connParams.CheckIsDefined() != nil {
		return ret
	}
	if ct.connParams.VpnType != vpn.WireGuard {
		return ret
	}

	baseSvr := ct.getServerByHostDnsName(ct.connParams.WireGuardParameters.EntryVpnServer.Hosts[0].DnsName)
	if baseSvr == nil {
		return ret
	}

	cLat := float64(baseSvr.Latitude)
	cLot := float64(baseSvr.Longitude)
	sort.Slice(ret, func(i, j int) bool {
		di := helpers.GetDistanceFromLatLonInKm(cLat, cLot, float64(ret[i].Latitude), float64(ret[i].Longitude))
		dj := helpers.GetDistanceFromLatLonInKm(cLat, cLot, float64(ret[j].Latitude), float64(ret[j].Longitude))
		return di < dj
	})
	return ret
}

func (ct connectivityTester) getServerByHostDnsName(hostDnsName string) *api_types.ServerInfoBase {
	for _, s := range ct.servers.WireguardServers {
		for _, h := range s.Hosts {
			if h.DnsName == hostDnsName {
				return &s.ServerInfoBase
			}
		}
	}
	return nil
}

// sortPorts() returns ports slice in port priority way
//
//	E.g. The default port (selected by user) has highest priority and must be checked first
func (ct connectivityTester) sortPorts(ports []api_types.PortInfo) (ret []api_types.PortInfo) {

	defaultPort := api_types.PortInfo{Port: ct.connParams.WireGuardParameters.Port.Port, Type: "UDP"}
	ret = append(ret, defaultPort)
	for _, p := range ports {
		if p.Port != 0 && !p.Equal(defaultPort) {
			ret = append(ret, p)
		}
	}

	return ret
}
