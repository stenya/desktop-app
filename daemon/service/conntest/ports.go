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
	"strings"
	"sync"

	atypes "github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/helpers"
	stypes "github.com/ivpn/desktop-app/daemon/service/types"
	"github.com/ivpn/desktop-app/daemon/vpn"
)

func (ct *connectivityTester) TestPorts(customPorts []stypes.PortData, getGeolookup func(timeoutMs int) (*atypes.GeoLookupResponse, error)) (
	ret map[stypes.PortData]bool,
	retRrr error) {

	if len(ct.servers.WireguardServers) == 0 || len(ct.servers.OpenvpnServers) == 0 {
		return nil, fmt.Errorf("servers not defined")
	}

	// Ports accessibility check based on checking ports connectivity to server(host):
	//	1. If client geolocation is known: use server nearest to the client's geolocation but from another country
	//	2. If client geolocation is NOT known: use server from the selected configuration
	//		2.1. If selected configuration not defined: use random server

	ret = make(map[stypes.PortData]bool)

	/*
		var err error
		var geolocation *api_types.GeoLookupResponse

		if getGeolookup != nil {
			// get geolocation into (API request)
			geolocation, _ = getGeolookup(1500) // 1500ms timeout
		}
	*/

	portsOvpn := getApplicablePorts(ct.servers.Config.Ports.OpenVPN, customPorts)
	portsOvpnTcp := getPortsByType(portsOvpn, stypes.TCP)

	// get selected server (server from last used configuration)
	selectedWgSvr, selectedOvpnSvr := ct.getSelectedServer()
	// if selected server not defined - get random server
	ct.initRandomServersIfNil(&selectedWgSvr, &selectedOvpnSvr)

	// Test OpenVPN TCP ports

	appendResults(ret, ct.testServerPorts_OpenvpnTcp(selectedOvpnSvr, portsOvpnTcp))

	// Test WireGuard UDP ports
	// ...

	log.Info(fmt.Sprintf("Ports test result: %v", ret))
	return ret, nil
}

func (ct *connectivityTester) testServerPorts_OpenvpnTcp(server atypes.ServerGeneric, ports []stypes.PortData) (ret map[stypes.PortData]bool) {
	ret = make(map[stypes.PortData]bool)
	hosts := server.GetHostsInfoBase()
	if len(hosts) == 0 || len(ports) == 0 {
		return
	}

	mutex := sync.Mutex{}

	testTcpPort := func(h atypes.HostInfoBase, p stypes.PortData) bool {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.Host, p.Port), constTimeout)
		if err != nil {
			return false
		}
		defer conn.Close()
		return true
	}

	guard := make(chan struct{}, 10) // maximum number of gotoutines running in same time
	wg := sync.WaitGroup{}

	host := hosts[helpers.RndInt(len(hosts))]
	log.Info(fmt.Sprintf("Testing OpenVPN TCP ports (destination: %s [%s:%s])", server.GetServerInfoBase().City, host.Hostname, host.Host))
	for _, port := range ports {
		guard <- struct{}{} // would block if guard channel is already filled
		wg.Add(1)
		go func(port stypes.PortData) {
			defer func() {
				wg.Done()
				<-guard
			}()
			appendResult(&mutex, ret, port, testTcpPort(host, port))
		}(port)
	}

	wg.Wait() // wait all routines to stop

	return
}

func (ct connectivityTester) getSelectedServer() (wgSvr, ovpnSvr atypes.ServerGeneric) {
	if err := ct.connParams.CheckIsDefined(); err != nil {
		return nil, nil
	}

	hostsBase := make([]atypes.HostInfoBase, 1)
	serversWg := ct.servers.ServersGenericWireguard()
	serversOvpn := ct.servers.ServersGenericOpenvpn()

	if ct.connParams.VpnType == vpn.WireGuard {
		for _, h := range ct.connParams.WireGuardParameters.EntryVpnServer.Hosts {
			hostsBase = append(hostsBase, h.HostInfoBase)
		}
		wgSvr = findServerByHosts(serversWg, hostsBase)
		if wgSvr != nil {
			ovpnSvr = findServerByID(wgSvr.GetServerInfoBase().Gateway, serversOvpn)
		}
	} else if ct.connParams.VpnType == vpn.OpenVPN {
		for _, h := range ct.connParams.OpenVpnParameters.EntryVpnServer.Hosts {
			hostsBase = append(hostsBase, h.HostInfoBase)
		}
		ovpnSvr = findServerByHosts(serversOvpn, hostsBase)
		if ovpnSvr != nil {
			wgSvr = findServerByID(ovpnSvr.GetServerInfoBase().Gateway, serversWg)
		}
	}

	return
}

func (ct connectivityTester) initRandomServersIfNil(wgSvr, ovpnSvr *atypes.ServerGeneric) {
	if *wgSvr != nil && *ovpnSvr != nil {
		return
	}

	serversWg := ct.servers.ServersGenericWireguard()
	serversOvpn := ct.servers.ServersGenericOpenvpn()

	if *wgSvr == nil && *ovpnSvr == nil {
		// random WG server
		svrs := ct.servers.WireguardServers
		*wgSvr = svrs[helpers.RndInt(len(svrs))]
		// use same server for OpenVPN
		*ovpnSvr = findServerByID((*wgSvr).GetServerInfoBase().Gateway, serversOvpn)
	} else if *wgSvr == nil {
		// OpenVPN server defined - use same GW server
		*wgSvr = findServerByID((*ovpnSvr).GetServerInfoBase().Gateway, serversWg)
	} else if *ovpnSvr == nil {
		// WG server defined - use same OpenVPN server
		*ovpnSvr = findServerByID((*wgSvr).GetServerInfoBase().Gateway, serversOvpn)
	}
	if *ovpnSvr == nil { // if OpenVPN server stil not defined - use random server
		svrs := ct.servers.OpenvpnServers
		*ovpnSvr = svrs[helpers.RndInt(len(svrs))]
	}
}

func findServerByID(gatewayID string, allServers []atypes.ServerGeneric) atypes.ServerGeneric {
	gatewayID = strings.Split(gatewayID, ".")[0]
	for _, s := range allServers {
		if gatewayID == strings.Split(s.GetServerInfoBase().Gateway, ".")[0] {
			return s
		}
	}
	return nil
}

func findServerByHosts(allServers []atypes.ServerGeneric, hosts []atypes.HostInfoBase) atypes.ServerGeneric {
	// hashing hosts dnsName-s (we will be matching this data with server hosts)
	hostDnsNames := make(map[string]struct{})
	for _, h := range hosts {
		hostDnsNames[h.DnsName] = struct{}{}
	}

	// if any of hosts belongs to a server - return this server
	for _, s := range allServers {
		for _, h := range s.GetHostsInfoBase() {
			if _, exists := hostDnsNames[h.DnsName]; exists {
				return s
			}
		}
	}
	return nil
}

func appendResult(mutex *sync.Mutex, dst map[stypes.PortData]bool, port stypes.PortData, isOk bool) {
	if mutex != nil {
		mutex.Lock()
		defer mutex.Unlock()
	}

	if val, exists := dst[port]; exists && val {
		return // do nothing if port already defined and accessible
	}

	dst[port] = isOk
}

func appendResults(dst map[stypes.PortData]bool, src map[stypes.PortData]bool) {

	for k, v := range src {
		if val, exists := dst[k]; exists && val {
			continue // do nothing if port already defined and accessible
		}
		dst[k] = v
	}
}

// getPortsByType() returns UDP or TCP ports
func getPortsByType(ports []stypes.PortData, portType stypes.PortType) (ret []stypes.PortData) {
	for _, p := range ports {
		if p.Protocol == portType {
			ret = append(ret, p)
		}
	}
	return
}

// Returns all applicable ports (default + custom). Custom ports are checked to fit allowed ranges from 'allPortsInfo' list.
func getApplicablePorts(allPortsInfo []atypes.PortInfo, customPorts []stypes.PortData) (ret []stypes.PortData) {
	ret = getPorts(allPortsInfo)
	ranges := getPortRanges(allPortsInfo)
	for _, cp := range customPorts {
		for _, r := range ranges {
			if (r.IsTCP() == (cp.Protocol == stypes.TCP)) &&
				cp.Port >= r.Range.Min && cp.Port <= r.Range.Max {
				ret = append(ret, cp)
				break
			}
		}
	}
	return
}

// getPortRanges() returns ports ranges (single port definitions are skipped)
func getPortRanges(all []atypes.PortInfo) (ret []atypes.PortInfo) {
	for _, pi := range all {
		if pi.Range.Max > 0 {
			ret = append(ret, pi)
		}
	}
	return
}

// getPorts() returns ports data (port ranges are skipped)
func getPorts(all []atypes.PortInfo) (ret []stypes.PortData) {
	for _, pi := range all {
		if pi.Port > 0 {
			proto := stypes.UDP
			if pi.IsTCP() {
				proto = stypes.TCP
			}
			ret = append(ret, stypes.PortData{Port: pi.Port, Protocol: proto})
		}
	}
	return
}
