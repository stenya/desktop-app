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

package wireguard

import (
	"fmt"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
)

func (wg *WireGuard) WaitForFirstHanshake(timeout time.Duration) error {
	endTime := time.Now().Add(timeout)
	tunnelName := wg.GetTunnelName()

	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to check handshake info: %w", err)
	}

	for {
		dev, err := client.Device(tunnelName)
		if err != nil {
			return fmt.Errorf("failed to check handshake info: %w", err)
		}

		for _, peer := range dev.Peers {
			if !peer.LastHandshakeTime.IsZero() {
				return nil // handshake detected
			}
		}

		if time.Now().After(endTime) {
			return fmt.Errorf("WireGuard handshake timeout")
		}
		time.Sleep(time.Millisecond * 10)
	}
}
