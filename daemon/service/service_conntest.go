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

package service

import (
	"fmt"

	"github.com/ivpn/desktop-app/daemon/service/conntest"
)

func (s *Service) ConnectionTestStop() error {
	// TODO: ...
	return nil
}

func (s *Service) ConnectionTestStart() error {
	if err := s.ConnectionTestStop(); err != nil {
		log.Error(err)
	}

	if err := s.Disconnect(); err != nil {
		log.Error(err)
	}
	// TODO: disabling killswitch is temporary (just for tests)!
	if err := s.SetKillSwitchIsPersistent(false); err != nil {
		log.Error(err)
	}
	if err := s.SetKillSwitchState(false); err != nil {
		log.Error(err)
		return err
	}

	// the function is asynchronous
	go func() {
		svrs, _ := s.ServersList()
		cTester := conntest.ConnectivityTester{}
		statusNotifyChan := make(chan conntest.StatusEvent)

		go func() {
			for {
				status, more := <-statusNotifyChan
				if !more {
					break
				}
				s._evtReceiver.OnConnectionTestStatus(status)

				msg := fmt.Sprintf("Server: %s (%s) Host=%s\t(%s:%d %s)...",
					status.Server.Country,
					status.Server.City,
					status.Host.Hostname,
					status.Host.Host,
					status.Port.Port,
					status.Port.Type)
				log.Info("Connection TEST: ", msg)
			}
		}()
		ci, err := cTester.Test(*svrs, s.Preferences().Session, s.GetConnectionParams(), statusNotifyChan)
		if err != nil {
			log.Info("Connection TEST Failed: ", err.Error())
			s._evtReceiver.OnConnectionTestResult(err, conntest.GoodConnectionInfo{})
		} else {
			log.Info("Connection TEST Success ", ci)
			s._evtReceiver.OnConnectionTestResult(nil, *ci)
		}
	}()

	return nil
}
