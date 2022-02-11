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

package dns

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/ivpn/desktop-app/daemon/netinfo"
	"github.com/ivpn/desktop-app/daemon/service/platform"
)

var (
	_fSetDNSByLocalIP      *syscall.LazyProc // DWORD _cdecl SetDNSByLocalIP(const char* interfaceLocalAddr, const char* dnsIP, byte operation, byte isDoH, const char* dohTemplateUrl, byte isIpv6)
	_fIsCanUseDnsOverHttps *syscall.LazyProc // DWORD _cdecl IsCanUseDnsOverHttps()
)

var dnsMutex sync.Mutex

// Operation enumerates possible DNS operations
type Operation uint32

// DNS operations
const (
	OperationSet Operation = 0
	OperationAdd Operation = 1
	OperationDel Operation = 2
)

// implInitialize doing initialization stuff (called on application start)
func implInitialize() error {
	helpersDllPath := platform.WindowsNativeHelpersDllPath()
	if len(helpersDllPath) == 0 {
		return fmt.Errorf("unable to initialize DNS wrapper: helpers dll path not initialized")
	}
	if _, err := os.Stat(helpersDllPath); err != nil {
		return fmt.Errorf("unable to initialize DNS wrapper (helpers dll not found) : '%s'", helpersDllPath)
	}

	dll := syscall.NewLazyDLL(helpersDllPath)
	_fSetDNSByLocalIP = dll.NewProc("SetDNSByLocalIP")           // DWORD _cdecl SetDNSByLocalIP(const char* interfaceLocalAddr, const char* dnsIP, byte operation, byte isDoH, const char* dohTemplateUrl, byte isIpv6)
	_fIsCanUseDnsOverHttps = dll.NewProc("IsCanUseDnsOverHttps") // DWORD _cdecl IsCanUseDnsOverHttps()
	return nil
}

func fSetDNSByLocalIP(interfaceLocalAddr net.IP, dnsCfg DnsSettings, ipv6 bool, op Operation) error {

	isDoH := uint32(0)
	switch dnsCfg.Encryption {
	case EncryptionDnsOverTls:
		return fmt.Errorf("DnsOverTls settings not supported by Windows. Please, try to use DnsOverHttps")
	case EncryptionDnsOverHttps:
		isDoH = 1
	default:
		isDoH = 0
	}

	dohTemplateUrl := dnsCfg.DohTemplate

	dnsIpString := ""
	if !dnsCfg.IsEmpty() {
		isAddrIpv6, _ := dnsCfg.IsIPv6()
		if isAddrIpv6 != ipv6 {
			return fmt.Errorf("unable to apply DNS configuration. IP address type mismatch to the IPv6 parameter")
		}
		dnsIpString = dnsCfg.Ip().String()
	}

	isIpv6 := uint32(0)
	if ipv6 {
		isIpv6 = 1
	}

	dnsMutex.Lock()
	defer dnsMutex.Unlock()

	retval, _, err := _fSetDNSByLocalIP.Call(
		uintptr(unsafe.Pointer(syscall.StringBytePtr(interfaceLocalAddr.String()))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(dnsIpString))),
		uintptr(op),
		uintptr(isDoH),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(dohTemplateUrl))),
		uintptr(isIpv6))

	return checkDefaultAPIResp(retval, err)
}

func fIsCanUseDnsOverHttps() bool {
	retval, _, err := _fIsCanUseDnsOverHttps.Call()
	if retval == 0 || err != syscall.Errno(0) {
		return false
	}
	return true
}

func checkDefaultAPIResp(retval uintptr, err error) error {
	if err != syscall.Errno(0) {
		return err
	}
	if retval != 0 {
		return fmt.Errorf("DNS change error: 0x%X", retval)
	}
	return nil
}

// last custom-DNS info which was enabled
var (
	_lastDNS DnsSettings
)

func catchPanic(err *error) {
	if r := recover(); r != nil {
		log.Error("PANIC (recovered): ", r)
		if e, ok := r.(error); ok {
			*err = e
		} else {
			*err = errors.New(fmt.Sprint(r))
		}
	}
}

// Pause - (on vpn paused) temporary restore OS default DNS parameters
func implPause() error {
	// Not in use for Windows implementation
	// In paused state we are simply switching to the main network interface (to default routes)

	// TODO: in case of custom DNS from local network - necessary to remove custom-DNS configuration from main (non-ivpn) network interface ???

	return nil
}

// Resume - (on vpn resumed) set VPN-defined DNS parameters
func implResume(defaultDNS DnsSettings) error {
	// Not in use for Windows implementation
	// In paused state we are simply switching to the main network interface (to default routes)

	// TODO: in case of custom DNS from local network - necessary to add (restore) custom-DNS configuration to main (non-ivpn) network interface ???

	return nil
}

func implGetDnsEncryptionAbilities() (dnsOverHttps, dnsOverTls bool, err error) {
	defer catchPanic(&err)
	return fIsCanUseDnsOverHttps(), false, err
}

func implSetManual(dnsCfg DnsSettings, localInterfaceIP net.IP) (err error) {
	defer catchPanic(&err)

	if dnsCfg.Equal(_lastDNS) {
		return nil
	}

	if !_lastDNS.IsEmpty() {
		// if there was defined DNS - remove it from non-VPN interfaces (if necessary)
		// (skipping VPN interface, because its data will be overwritten)
		if err := implDeleteManual(nil); err != nil {
			return fmt.Errorf("failed to set DNS: %w", err)
		}
	}

	isIpv6 := false
	if dnsCfg.IsEmpty() {
		return fmt.Errorf("unable to change DNS (configuration is not defined)")
	}
	isIpv6, _ = dnsCfg.IsIPv6()

	// non-VPN interfaces to update (if DNS located in local network)
	notVpnInterfacesToUpdate, err := getInterfacesIPsWhichContainsIP(dnsCfg.Ip(), localInterfaceIP)

	if localInterfaceIP == nil && len(notVpnInterfacesToUpdate) <= 0 {
		return nil
	}

	start := time.Now()
	log.Info(fmt.Sprintf("Changing DNS to %s ...", dnsCfg.InfoString()))
	defer func() {
		if err != nil {
			log.Error(fmt.Sprintf("Changing DNS to %s done (%dms) with error: %s", dnsCfg.InfoString(), time.Since(start).Milliseconds(), err.Error()))
		} else {
			log.Info(fmt.Sprintf("Changing DNS to %s: done (%dms)", dnsCfg.InfoString(), time.Since(start).Milliseconds()))
		}
	}()

	if localInterfaceIP != nil {
		// SET DNS to VPN interface
		if err := fSetDNSByLocalIP(localInterfaceIP, dnsCfg, isIpv6, OperationSet); err != nil {
			return fmt.Errorf("failed to set DNS for local interface: %w", err)
		}
	}

	if len(notVpnInterfacesToUpdate) > 0 {
		// ADD DNS to non-VPN interface (if necessary, when DNS is in local network)
		for _, ifcAddr := range notVpnInterfacesToUpdate {
			if err := fSetDNSByLocalIP(ifcAddr.IP, dnsCfg, isIpv6, OperationAdd); err != nil {
				return fmt.Errorf("failed to set DNS for non-VPN interface: %w", err)
			}
		}
	}

	// save last changed DNS address
	_lastDNS = dnsCfg

	return nil
}

func implDeleteManual(localInterfaceIP net.IP) (retErr error) {
	defer catchPanic(&retErr)

	// non-VPN interfaces to update (if DNS server is in local network)
	notVpnInterfacesToUpdate, err := getInterfacesIPsWhichContainsIP(_lastDNS.Ip(), localInterfaceIP)

	if localInterfaceIP == nil && len(notVpnInterfacesToUpdate) <= 0 {
		return nil
	}

	start := time.Now()
	log.Info("Restoring default DNS...")
	defer func() {
		if err != nil {
			log.Info(fmt.Sprintf("Restoring default DNS done (%dms) with error: %s", time.Since(start).Milliseconds(), err.Error()))
		} else {
			log.Info(fmt.Sprintf("Restoring default DNS: done (%dms)", time.Since(start).Milliseconds()))
		}
	}()

	isIpv6 := false
	if !_lastDNS.IsEmpty() {
		isIpv6, _ = _lastDNS.IsIPv6()
	}

	if localInterfaceIP != nil {
		// RESET DNS for VPN interface
		if err := fSetDNSByLocalIP(localInterfaceIP, DnsSettings{}, false, OperationSet); err != nil {
			retErr = fmt.Errorf("failed to reset DNS for local interface: %w", err)
		}
		if isIpv6 {
			if err := fSetDNSByLocalIP(localInterfaceIP, DnsSettings{}, true, OperationSet); err != nil {
				log.Error(fmt.Errorf("failed to reset DNS (IPv6) for local interface: %w", err))
			}
		}
	}

	if len(notVpnInterfacesToUpdate) > 0 {
		// REMOVE DNS from non-VPN interface (if necessary, when DNS is in local network)
		for _, ifcAddr := range notVpnInterfacesToUpdate {
			if err := fSetDNSByLocalIP(ifcAddr.IP, _lastDNS, isIpv6, OperationDel); err != nil {
				log.Error(fmt.Errorf("failed to remove previously applied DNS configuration for non-VPN interface (ipv6:%v): %w", isIpv6, err))
			}
		}
	}

	_lastDNS = DnsSettings{}

	return retErr
}

func implGetPredefinedDnsConfigurations() ([]DnsSettings, error) {
	return []DnsSettings{}, nil
}

// getInterfacesIPsWhichContainsIP - get IP addresses of local network interfaces to which belongs an IP address
// (interface which is in same network as 'addr')
// 		addr - IP address from local network (which can be accessed by interface)
//		localAddrToSkip - local IP of interface which can be excluded from output (e.g. VPN interface)
func getInterfacesIPsWhichContainsIP(addr net.IP, localAddrToSkip net.IP) (ret []net.IPNet, err error) {
	if addr == nil {
		return ret, nil
	}

	// get interfaces which must be modified by new DNS value
	networks, err := netinfo.GetAllLocalV4Addresses()
	if err != nil {
		return nil, fmt.Errorf("error receiving local V4 addresses : %w", err)
	}

	for _, network := range networks {

		if network.IP.Equal(localAddrToSkip) {
			continue
		}

		if network.Contains(addr) { // 'addr' is in 'network'
			ret = append(ret, network)
		}
	}

	return ret, nil
}
