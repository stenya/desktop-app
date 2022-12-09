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

package wireguard

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ivpn/desktop-app/daemon/logger"
	"github.com/ivpn/desktop-app/daemon/netinfo"
	"github.com/ivpn/desktop-app/daemon/service/dns"
	"github.com/ivpn/desktop-app/daemon/service/platform"
	"github.com/ivpn/desktop-app/daemon/shell"
	"github.com/ivpn/desktop-app/daemon/vpn"
)

// TODO: BE CAREFUL! Constant string! (can be changed after WireGuard update)
const (
	strTriggerSuccessInit      string = "UAPI listener started"
	strTriggerAddrAlreadyInUse string = "Address already in use"
)

// IsUpdateDefaultRoute determines the style of default route configuration
//
// Normally, the "0/1 <VPN_GW>..." route is in use
// But such route configuration breaks iCloud synchronization (e.f. Safari bookmarks not synchronizing when VPN connected)
// Info: iCloud synchronization ignores the "0/1" route and uses "default" with higher priority, but this direct route is not workable when VPN is connected.
//
// To fix this problem, the routing must be configured in a different way:
//  1. The original default route must be scoped to the default interface (this will allow us to create a new default route to the VPN gateway)
//     *Note: do not forget to restore this to the original state after VPN disconnect (the standalone scoped 'default' route does not work for some reason)
//  2. Use "default <VPN_GW>..." route to access VPN gateway
//
// Example:
//
//	(IsUpdateDefaultRoute==false):
//	  Destination        Gateway            Flags           Netif Expire
//	  0/1                172.16.0.1         UGScg           utun6
//	  default            192.168.0.1        UGScg             en0
//	(IsUpdateDefaultRoute==true):
//	  Destination        Gateway            Flags           Netif Expire
//	  default            172.16.0.1         UGScg           utun6
//	  default            192.168.0.1        UGScIg            en0
const IsUpdateDefaultRoute bool = true

const subnetMask string = "255.0.0.0"
const subnetMaskPrefixLenIPv6 string = "64"

// internalVariables of wireguard implementation for macOS
type internalVariables struct {
	// WG running process (shell command)
	command       *exec.Cmd
	isGoingToStop bool
	defaultRoute  netinfo.Route

	isDefaultRouteUpdated bool // 'true' if the default route successfully updated (see: IsUpdateDefaultRoute == true)

	isPaused      bool
	omResumedChan chan struct{} // channel for 'On Resume' events
}

var logWgOut *logger.Logger

func (wg *WireGuard) init() error {
	logWgOut = logger.NewLogger("wg_out")
	return nil
}

// connect - SYNCHRONOUSLY execute openvpn process (wait until it finished)
func (wg *WireGuard) connect(stateChan chan<- vpn.StateInfo) (err error) {
	wg.internals.omResumedChan = make(chan struct{}, 1)
	defer func() {
		// The 'Pause' functionality is based on fact that connection will be re-connected by a service
		// if we disconnected without any 'disconnect' request.
		// Therefore, in case of 'pause' we just stopping real connection
		// and waiting for 'resume' command to return control to the owner service.
		if wg.internals.isPaused && !wg.internals.isGoingToStop {
			// waiting to 'resume' event
			<-wg.internals.omResumedChan
			err = &vpn.ReconnectionRequiredError{Err: err}
		}
	}()

	return wg.internalConnect(stateChan)
}

// connect - SYNCHRONOUSLY execute openvpn process (wait until it finished)
func (wg *WireGuard) internalConnect(stateChan chan<- vpn.StateInfo) error {

	var routineStopWaiter sync.WaitGroup

	// if we are trying to connect when no connectivity (WiFi off?) -
	// waiting until network appears
	// Retry to check each 5 seconds (sending RECONNECTING event)
	for !wg.internals.isGoingToStop {
		if dns.IsPrimaryInterfaceFound() {
			break
		}
		log.Info("No connectivity. Waiting 5 sec to retry...")

		stateChan <- vpn.NewStateInfo(vpn.RECONNECTING, "No connectivity")
		pauseEnd := time.Now().Add(time.Second * 5)
		for time.Now().Before(pauseEnd) && !wg.internals.isGoingToStop {
			time.Sleep(time.Millisecond * 50)
		}
	}

	// get default Gateway IP
	defaultRoute, err := netinfo.DefaultRoute()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to detect default gateway: %s", err))
		return err
	}
	wg.internals.defaultRoute = defaultRoute

	if wg.internals.isGoingToStop {
		return nil
	}

	defer func() {
		wg.removeRoutes()
		wg.removeDNS()

		// wait to stop all routines
		routineStopWaiter.Wait()

		log.Info("Stopped")
	}()

	utunName, err := getFreeTunInterfaceName()
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("unable to start WireGuard. Failed to obtain free utun interface: %w", err)
	}

	log.Info("Starting WireGuard in interface ", utunName)
	// LOG_LEVEL=verbose
	wg.internals.command = exec.Command(wg.binaryPath, "-f", utunName)
	wg.internals.command.Env = os.Environ()
	wg.internals.command.Env = append(wg.internals.command.Env, "LOG_LEVEL=verbose")

	isStartedChannel := make(chan bool)

	// output reader
	outPipe, err := wg.internals.command.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to start WireGuard: %w", err)
	}

	// wait for WG initialization + logging all output
	outPipeScanner := bufio.NewScanner(outPipe)
	routineStopWaiter.Add(1)
	go func() {
		defer routineStopWaiter.Done()

		isWaitingToStart := true
		for outPipeScanner.Scan() && wg.internals.command.ProcessState == nil {
			text := outPipeScanner.Text()
			logWgOut.Info(text) // logging the output

			if isWaitingToStart && strings.Contains(text, strTriggerSuccessInit) {
				isWaitingToStart = false
				isStartedChannel <- true
			}
		}
	}()

	// error reader
	errPipe, err := wg.internals.command.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to start WireGuard: %w", err)
	}
	errPipeScanner := bufio.NewScanner(errPipe)
	routineStopWaiter.Add(1)
	go func() {
		defer routineStopWaiter.Done()

		for errPipeScanner.Scan() {
			logWgOut.Info("[err] ", errPipeScanner.Text())
		}
	}()

	// start
	if err := wg.internals.command.Start(); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed to start WireGuard process: %w", err)
	}

	var initError error = nil

	// waiting to start and initialize
	routineStopWaiter.Add(1)
	go func() {
		defer routineStopWaiter.Done()
		isHaveToBeStopped := false

		select {
		case <-isStartedChannel:
			// Process started. Perform initialization...
			if initError = wg.initialize(utunName); initError != nil {
				// (return initialization error as a result of connect)
				log.ErrorTrace(initError)
				isHaveToBeStopped = true
			} else {
				log.Info("Started")
				// CONNECTED
				wg.notifyConnectedStat(stateChan)
			}

		case <-time.After(time.Second * 5):
			// stop process if WG not successfully started during 5 sec
			err = fmt.Errorf("WireGuard process initialization timeout")
			if initError == nil {
				initError = err
			}
			log.Error(err)
			isHaveToBeStopped = true
		}

		if isHaveToBeStopped {
			log.Error("Stopping process manually...")
			if err := wg.disconnect(); err != nil {
				log.Error("Failed to stop process: ", err)
			}
		}
	}()

	if wg.internals.isGoingToStop {
		wg.disconnect()
	}

	if err := wg.internals.command.Wait(); err != nil {
		// error will be received anyway. We are logging it only if process was stopped unexpectedly
		if !wg.internals.isGoingToStop {
			log.Error(err.Error())
			return fmt.Errorf("WireGuard process error: %w", err)
		}
	}
	return initError
}

func (wg *WireGuard) disconnect() error {
	wg.internals.isGoingToStop = true
	log.Info("Stopping")
	wg.resume()
	return wg.internalDisconnect()
}

func (wg *WireGuard) internalDisconnect() error {
	cmd := wg.internals.command

	// ProcessState contains information about an exited process,
	// available after a call to Wait or Run.
	// NOT nil = process finished
	if cmd == nil || cmd.Process == nil || cmd.ProcessState != nil {
		return nil // nothing to stop
	}

	log.Info("Stopping")
	return cmd.Process.Kill()
}

func (wg *WireGuard) isPaused() bool {
	return wg.internals.isPaused
}

func (wg *WireGuard) pause() error {
	wg.internals.isPaused = true
	return wg.internalDisconnect()
}

func (wg *WireGuard) resume() error {
	// send 'resumed' event
	resumeCh := wg.internals.omResumedChan
	if resumeCh != nil {
		select {
		case resumeCh <- struct{}{}:
		default:
		}
	}

	wg.internals.isPaused = false
	return nil
}

func (wg *WireGuard) setManualDNS(dnsCfg dns.DnsSettings) error {
	return dns.SetManual(dnsCfg, nil)
}

func (wg *WireGuard) resetManualDNS() error {
	return dns.DeleteManual(wg.DefaultDNS(), nil)
}

func (wg *WireGuard) initialize(utunName string) error {

	// Init IPv6 DNS resolver (if necessary);
	// It should be done before initialization of the tunnel interface
	if err := wg.initIPv6DNSResolver(utunName); err != nil {
		log.Error(fmt.Errorf("failed to initialize IPv6 DNS resolver: %w", err))
	}

	if err := wg.initializeConfiguration(utunName); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	if err := wg.setRoutes(); err != nil {
		return fmt.Errorf("failed to set routes: %w", err)
	}

	err := wg.setDNS()
	if err != nil {
		return fmt.Errorf("failed to set DNS: %w", err)
	}
	return nil
}

func (wg *WireGuard) initializeConfiguration(utunName string) error {
	log.Info("Configuring ", utunName, " interface...")

	// Configure WireGuard interface
	// example command:	ipconfig set utun7 MANUAL 10.0.0.121 255.255.255.0
	//				 	ifconfig utun2 inet 172.26.22.146/8 172.26.22.146 alias
	if err := wg.initializeUnunInterface(utunName); err != nil {
		return fmt.Errorf("failed to initialize interface: %w", err)
	}

	// WireGuard configuration
	if err := wg.setWgConfiguration(utunName); err != nil {
		return err
	}

	if wg.connectParams.mtu > 0 {
		// Custom MTU
		log.Info(fmt.Sprintf("Configuring custom MTU = %d ...", wg.connectParams.mtu))
		err := shell.Exec(log, "/sbin/ifconfig", utunName, "mtu", strconv.Itoa(wg.connectParams.mtu))
		if err != nil {
			return fmt.Errorf("failed to set custom MTU (%d): %w", wg.connectParams.mtu, err)
		}
	}

	return nil
}

// Configure WireGuard interface
// example command: ipconfig set utun7 MANUAL 10.0.0.121 12
// example command: ipconfig set utun7 MANUAL-V6 fd00:4956:504e:ffff::ac1a:704b 96
func (wg *WireGuard) initializeUnunInterface(utunName string) error {
	// initialize IPv4 interface for tunnel
	if err := shell.Exec(log, "/usr/sbin/ipconfig", "set", utunName, "MANUAL", wg.connectParams.clientLocalIP.String(), subnetMask); err != nil {
		return fmt.Errorf("failed to set the IPv4 address for interface: %w", err)
	}

	// initialize IPv6 interface for tunnel
	ipv6LocalIP := wg.connectParams.GetIPv6ClientLocalIP()
	if ipv6LocalIP != nil {
		if err := shell.Exec(log, "/usr/sbin/ipconfig", "set", utunName, "MANUAL-V6", ipv6LocalIP.String(), subnetMaskPrefixLenIPv6); err != nil {
			return fmt.Errorf("failed to set the IPv6 address for interface: %w", err)
		}
	}
	return nil
}

// WireGuard configuration
func (wg *WireGuard) setWgConfiguration(utunName string) error {
	// do not forget to remove config file after finishing configuration
	defer os.Remove(wg.configFilePath)

	for retries := 0; ; retries++ {
		// few retries if local port is already in use
		if retries >= 5 {
			// not more than 5 retries
			return fmt.Errorf("failed to set wireguard configuration")
		}

		// generate configuration
		err := wg.generateAndSaveConfigFile(wg.configFilePath)
		if err != nil {
			return fmt.Errorf("failed to save WG config file: %w", err)
		}

		// define output processing function
		isPortInUse := false
		errParse := func(text string, isError bool) {
			if isError {
				log.Debug("[wgconf error] ", text)
			} else {
				log.Debug("[wgconf out] ", text)
			}
			if strings.Contains(text, strTriggerAddrAlreadyInUse) {
				isPortInUse = true
			}
		}

		// Configure WireGuard
		// example command: wg setconf utun7 wireguard.conf
		err = shell.ExecAndProcessOutput(log, errParse, "", wg.toolBinaryPath,
			"setconf", utunName, wg.configFilePath)

		if !isPortInUse {
			return err
		}
	}
}

func (wg *WireGuard) setRoutes() error {
	log.Info("Modifying routing table...")

	if net.IPv4(127, 0, 0, 1).Equal(wg.connectParams.hostIP) {
		return fmt.Errorf("WG server IP error (unable to use '127.0.0.1' as WG server IP)")
	}

	// Update main route
	wg.internals.isDefaultRouteUpdated = false
	if IsUpdateDefaultRoute {
		func() (retErr error) {
			isDefaultRouteDeleted := false
			defer func() {
				if retErr == nil {
					wg.internals.isDefaultRouteUpdated = true
				} else {
					log.Error(retErr)
					// In case of error - revert changes were made
					shell.Exec(log, "/sbin/route", "-n", "delete", "default", wg.internals.defaultRoute.GatewayIP.String(), "-ifscope", wg.internals.defaultRoute.InterfaceName)
					shell.Exec(log, "/sbin/route", "-n", "delete", "default", wg.connectParams.hostLocalIP.String())
					if isDefaultRouteDeleted {
						shell.Exec(log, "/sbin/route", "-n", "add", "default", wg.internals.defaultRoute.GatewayIP.String())
					}
				}
			}()
			// Example:
			//	  Destination        Gateway            Flags           Netif Expire
			//	  default            172.16.0.1         UGScg           utun6
			//	  default            192.168.0.1        UGScIg            en0
			//
			//	route add default ${_argDefGwIp} -ifscope ${_argDefGwInterface}
			//	route delete default {DEF_INTERFACE}
			//	route add default ${_argLocalVpnHostIp}

			//  The original default route must be scoped to the default interface (this will allow us to create a new default route to the VPN gateway)
			//    *Note: do not forget to restore this to the original state after VPN disconnect (the standalone scoped 'default' route does not work for some reason)
			if err := shell.Exec(log, "/sbin/route", "-n", "add", "default", wg.internals.defaultRoute.GatewayIP.String(), "-ifscope", wg.internals.defaultRoute.InterfaceName); err != nil {
				return fmt.Errorf("route shell command error : %w", err)
			}
			// 	Remove original default route
			if err := shell.Exec(log, "/sbin/route", "-n", "delete", "default", wg.internals.defaultRoute.GatewayIP.String()); err != nil {
				return fmt.Errorf("route shell command error : %w", err)
			} else {
				isDefaultRouteDeleted = true
			}
			//  Use "default <VPN_GW>..." route to access VPN gateway
			if err := shell.Exec(log, "/sbin/route", "-n", "add", "default", wg.connectParams.hostLocalIP.String()); err != nil {
				return fmt.Errorf("route shell command error : %w", err)
			}
			return nil
		}()
	}
	if !wg.internals.isDefaultRouteUpdated {
		// example command:	route	-n	add	-net	0/1			10.0.0.1
		// 					route	-n	add	-inet	0.0.0.0/1	-interface utun2
		if err := shell.Exec(log, "/sbin/route", "-n", "add", "-inet", "-net", "0/1", wg.connectParams.hostLocalIP.String()); err != nil {
			return fmt.Errorf("adding route shell command error : %w", err)
		}
	}

	// Update routing to remote server (remote_server default_router 255.255.255)
	// example command:	route	-n	add	-net	145.239.239.55	192.168.1.1	255.255.255.255
	//					route	-n	add	-inet	51.77.91.106	-gateway	192.168.1.1
	if err := shell.Exec(log, "/sbin/route", "-n", "add", "-inet", "-net", wg.connectParams.hostIP.String(), wg.internals.defaultRoute.GatewayIP.String(), "255.255.255.255"); err != nil {
		return fmt.Errorf("adding route shell command error : %w", err)
	}

	// Update routing table
	// example command:	route	-n	add	-net	128.0.0.0	10.0.0.1	128.0.0.0
	// 					route	-n	add	-inet	128.0.0.0/1	-interface	utun2
	if err := shell.Exec(log, "/sbin/route", "-n", "add", "-inet", "-net", "128.0.0.0", wg.connectParams.hostLocalIP.String(), "128.0.0.0"); err != nil {
		return fmt.Errorf("adding route shell command error : %w", err)
	}

	ipv6HostLocalIP := wg.connectParams.GetIPv6HostLocalIP()
	if ipv6HostLocalIP != nil {
		// Using the default gateway (a ::/0 netmask) as two /1 networks: ::/1 and 8000::/1.
		// Since a more specific route always wins, this forces traffic to be routed via the VPN instead of over the default gateway.
		// Additionally, this does not change the current 'default' route (do not break users configuration after disconnection).
		if err := shell.Exec(log, "/sbin/route", "-n", "add", "-inet6", "-net", "::/1", ipv6HostLocalIP.String()); err != nil {
			return fmt.Errorf("adding route shell command error : %w", err)
		}
		if err := shell.Exec(log, "/sbin/route", "-n", "add", "-inet6", "-net", "8000::/1", ipv6HostLocalIP.String()); err != nil {
			return fmt.Errorf("adding route shell command error : %w", err)
		}
	}

	return nil
}

func (wg *WireGuard) removeRoutes() error {
	log.Info("Restoring routing table...")

	if wg.internals.isDefaultRouteUpdated {
		shell.Exec(log, "/sbin/route", "-n", "delete", "default", wg.connectParams.hostLocalIP.String())
		shell.Exec(log, "/sbin/route", "-n", "add", "default", wg.internals.defaultRoute.GatewayIP.String())
		shell.Exec(log, "/sbin/route", "-n", "delete", "default", wg.internals.defaultRoute.GatewayIP.String(), "-ifscope", wg.internals.defaultRoute.InterfaceName)
	}

	shell.Exec(log, "/sbin/route", "-n", "delete", "-inet", "-net", "0/1", wg.connectParams.hostLocalIP.String())
	shell.Exec(log, "/sbin/route", "-n", "delete", "-inet", "-net", wg.connectParams.hostIP.String())
	shell.Exec(log, "/sbin/route", "-n", "delete", "-inet", "-net", "128.0.0.0", wg.connectParams.hostLocalIP.String())

	ipv6HostLocalIP := wg.connectParams.GetIPv6HostLocalIP()
	if ipv6HostLocalIP != nil {
		// Using the default gateway (a ::/0 netmask) as two /1 networks: ::/1 and 8000::/1.
		// Since a more specific route always wins, this forces traffic to be routed via the VPN instead of over the default gateway.
		// Additionally, this does not change the current 'default' route (do not break users configuration after disconnection).
		shell.Exec(log, "/sbin/route", "-n", "delete", "-inet6", "-net", "::/1", ipv6HostLocalIP.String())
		shell.Exec(log, "/sbin/route", "-n", "delete", "-inet6", "-net", "8000::/1", ipv6HostLocalIP.String())
	}
	return nil
}

func (wg *WireGuard) onRoutingChanged() error {
	defRoute, err := netinfo.DefaultRoute()
	if err != nil {
		log.Warning(fmt.Sprintf("onRoutingChanged: %v", err))
		return err
	}

	needUpdateRoutes := false
	if wg.internals.isDefaultRouteUpdated {
		if !wg.connectParams.hostLocalIP.Equal(defRoute.GatewayIP) {
			log.Info(fmt.Sprintf("Default gateway changed: %s -> %s. Updating routes...", wg.connectParams.hostLocalIP.String(), defRoute.GatewayIP.String()))
			needUpdateRoutes = true

		}
	} else if !wg.internals.defaultRoute.GatewayIP.Equal(defRoute.GatewayIP) {
		log.Info(fmt.Sprintf("Default gateway changed: %s -> %s. Updating routes...", wg.internals.defaultRoute.GatewayIP.String(), defRoute.GatewayIP.String()))
		needUpdateRoutes = true
	}

	if needUpdateRoutes {
		wg.internals.defaultRoute = defRoute
		wg.removeRoutes()
		wg.setRoutes()
	}
	return nil
}

func (wg *WireGuard) setDNS() error {
	defaultDNS := wg.DefaultDNS()
	log.Info("Updating DNS server to " + defaultDNS.String() + "...")
	err := shell.Exec(log, platform.DNSScript(), "-up_set_dns", defaultDNS.String())
	if err != nil {
		return fmt.Errorf("failed to change DNS: %w", err)
	}
	return nil
}

func (wg *WireGuard) initIPv6DNSResolver(utunName string) error {
	// required to be able to resolve IPv6 DNS addresses by the default macOS's domain name resolver
	ipv6LocalIP := wg.connectParams.GetIPv6ClientLocalIP()
	if ipv6LocalIP != nil && len(utunName) > 0 {
		err := shell.Exec(log, platform.DNSScript(), "-up_init_ipv6_resolver", ipv6LocalIP.String(), utunName)
		if err != nil {
			return fmt.Errorf("failed to change DNS: %w", err)
		}
	}
	return nil
}

func (wg *WireGuard) removeDNS() error {
	log.Info("Restoring DNS server.")
	err := shell.Exec(log, platform.DNSScript(), "-down", wg.DefaultDNS().String())
	if err != nil {
		return fmt.Errorf("failed to restore DNS: %w", err)
	}

	return nil
}

func getFreeTunInterfaceName() (string, error) {
	utunNameRegExp := regexp.MustCompile("^utun([0-9]+)")

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	maxUtunNo := 0
	for _, ifs := range ifaces {
		strs := utunNameRegExp.FindStringSubmatch(ifs.Name)
		if len(strs) == 2 {
			if utunNo, _ := strconv.Atoi(strs[1]); utunNo > maxUtunNo {
				maxUtunNo = utunNo
			}
		}
	}

	return fmt.Sprintf("utun%d", maxUtunNo+1), nil
}

func (wg *WireGuard) getOSSpecificConfigParams() (interfaceCfg []string, peerCfg []string) {

	// TODO: check if we need it for this platform
	// Same as "0.0.0.0/0" but such type of configuration is disabling internal WireGuard-s Firewall
	// It blocks everything except WireGuard traffic.
	// We need to disable WireGuard-s firewall because we have our own implementation of firewall.
	//  For details, refer to WireGuard-windows sources: tunnel\ifaceconfig.go (enableFirewall(...) method)

	if len(wg.connectParams.GetIPv6HostLocalIP()) > 0 {
		peerCfg = append(peerCfg, "AllowedIPs = 128.0.0.0/1, 0.0.0.0/1, ::/0")
	} else {
		peerCfg = append(peerCfg, "AllowedIPs = 128.0.0.0/1, 0.0.0.0/1")
	}

	return interfaceCfg, peerCfg
}
