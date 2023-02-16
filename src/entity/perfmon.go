/*
 *   sonic-android-supply  Supply of ADB.
 *   Copyright (C) 2022  SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package entity

import "encoding/json"

type PerfmonData struct {
	System  *SystemInfo  `json:"system,omitempty"`
	Process *ProcessInfo `json:"process,omitempty"`
}

func (p *PerfmonData) ToJson() string {
	str, _ := json.Marshal(p)
	return string(str)
}
func (p *PerfmonData) ToString() string {
	return p.ToJson()
}
func (p *PerfmonData) ToFormat() string {
	str, _ := json.MarshalIndent(p, "", "\t")
	return string(str)
}
