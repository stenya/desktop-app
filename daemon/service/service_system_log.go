//
//  Daemon for IVPN Client Desktop
//  https://github.com/ivpn/desktop-app
//
//  Created by Stelnykovych Alexandr.
//  Copyright (c) 2022 Privatus Limited.
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

import "fmt"

type SystemLogMessageType int

const (
	Info    SystemLogMessageType = iota
	Warning SystemLogMessageType = iota
	Error   SystemLogMessageType = iota
)

type SystemLogMessage struct {
	Type    SystemLogMessageType
	Message string
}

func (s *Service) systemLog(mesType SystemLogMessageType, message string) bool {
	ch := s._systemLog
	if ch == nil {
		switch mesType {
		case Info:
			log.Info(fmt.Sprintf("(syslog not initialized) INFO: %s", message))
		case Warning:
			log.Info(fmt.Sprintf("(syslog not initialized) WARNING: %s", message))
		case Error:
			log.Info(fmt.Sprintf("(syslog not initialized) ERROR: %s", message))
		default:
		}

		return false
	}

	ch <- SystemLogMessage{Message: message, Type: mesType}
	return true
}
