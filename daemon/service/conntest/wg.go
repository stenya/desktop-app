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
	"time"

	api_types "github.com/ivpn/desktop-app/daemon/api/types"
	"github.com/ivpn/desktop-app/daemon/service/platform"
	"github.com/ivpn/desktop-app/daemon/vpn"
	"github.com/ivpn/desktop-app/daemon/vpn/wireguard"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type ConnectivityTesterWireguard struct {
	wg *wireguard.WireGuard

	localIP          net.IP
	privateKey       string
	privateKeyParsed wgtypes.Key

	// channel closes when exiting from synchronous 'Connect' function
	connectionFinishedChan chan struct{}
}

func InitTesterWireguard(onUninitialise chan struct{}, localIP net.IP, privateKey string) (*ConnectivityTesterWireguard, error) {
	if localIP == nil || localIP.IsUnspecified() {
		return nil, fmt.Errorf("wireguard local IP not specified")
	}

	privKey, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return nil, err
	}

	obj := &ConnectivityTesterWireguard{
		localIP:                localIP,
		privateKey:             privateKey,
		privateKeyParsed:       privKey,
		connectionFinishedChan: make(chan struct{})}

	// init wireguard device
	if err := obj.initWireguardDevice(); err != nil {
		return nil, err
	}

	go func() {
		<-onUninitialise
		obj.wg.Disconnect()
	}()

	return obj, nil
}

func (wct *ConnectivityTesterWireguard) WaitForDisconnect() {
	<-wct.connectionFinishedChan
}

func (wct *ConnectivityTesterWireguard) initWireguardDevice() error {
	// basic initialisation parameters
	// we do not care about real connectivity; we need only initiate WG device
	wgConnParams := wireguard.CreateConnectionParams(
		"",                       // miltihop exit host name
		2049,                     // host port
		net.ParseIP("127.0.0.1"), // host IP - use local
		"rg+GGDmjM4Vxo1hURvKmgm9yonb6qcoKbPCP/DNDBnI=", // host public key - any random key
		net.ParseIP("172.16.0.1"),                      // host local IP
		"",                                             // ipv6 pefix
		0)                                              // mtu

	if len(wct.privateKey) == 0 || wct.localIP.IsUnspecified() {
		return fmt.Errorf("WireGuard credentials are not defined (please, regenerate WG credentials or re-login)")
	}
	wgConnParams.SetCredentials(wct.privateKey, wct.localIP)

	// Init WG interface and first connection
	wg, err := wireguard.NewWireGuardObject(platform.WgBinaryPath(),
		platform.WgToolBinaryPath(),
		platform.WGConfigFilePath(), wgConnParams)
	if err != nil {
		return err
	}
	wct.wg = wg

	// Mark connection as only for tests. It is important to not change any connectivity parameters on the device
	wg.MarkAsTestConnection()

	if err := wg.Init(); err != nil {
		return fmt.Errorf("failed to initialize VPN object: %w", err)
	}

	// channel notifies when device initialised
	onInitializedChan := make(chan error)

	go func() {
		// init 'status' channel and start reading it
		stateChan := make(chan vpn.StateInfo)
		defer func() {
			close(wct.connectionFinishedChan)
			//close(stateChan)
		}()
		go func() {
			for {
				select {
				case _, more := <-wct.connectionFinishedChan:
					if !more {
						return // channel is closed
					}
				case state := <-stateChan:
					if state.State == vpn.CONNECTED {
						select {
						case onInitializedChan <- nil:
						default:
						}
					}
				}
			}
		}()

		// Start connection just to initialize WireGuard device (we do not care about real connectivity here)
		err = wg.Connect(stateChan)
		if err != nil {
			onInitializedChan <- err
		}
	}()

	if connErr := <-onInitializedChan; connErr != nil {
		return connErr
	}

	return nil
}

func (wct ConnectivityTesterWireguard) Test(host api_types.WireGuardServerHostInfo, port int) error {
	devName := wct.wg.GetTunnelName()

	// Wireguard control client
	wgCtrlClient, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgCtrlClient.Close()

	wgDev, err := wireguard.GetCtrlDevice(devName, wgCtrlClient)
	if err != nil {
		return err
	}

	listenPort := wgDev.ListenPort
	kaInterval := time.Second * 60

	pubKey, err := wgtypes.ParseKey(host.PublicKey)
	if err != nil {
		return err
	}
	ep, err := net.ResolveUDPAddr("udp", host.Host+":"+fmt.Sprint(port))
	if err != nil {
		return err
	}
	cfg := wgtypes.Config{}
	cfg.PrivateKey = &wct.privateKeyParsed
	cfg.ListenPort = &listenPort
	cfg.ReplacePeers = true

	pcfg := wgtypes.PeerConfig{}
	pcfg.PublicKey = pubKey
	pcfg.Endpoint = ep
	pcfg.PersistentKeepaliveInterval = &kaInterval

	cfg.Peers = []wgtypes.PeerConfig{pcfg}

	err = wgCtrlClient.ConfigureDevice(devName, cfg)
	if err != nil {
		return err
	}

	var isDisconnectRequested bool
	defer func() {
		isDisconnectRequested = true
	}()

	return wireguard.WaitForWireguardFirstHanshake(devName, time.Millisecond*200, &isDisconnectRequested, nil)
}
