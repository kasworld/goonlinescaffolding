// Copyright 2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package authdata

import (
	"fmt"

	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_authorize"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
)

var Authkey2Admin = map[string][2][]string{}

func AddAdminKey(key string) error {
	var err error
	if _, exist := Authkey2Admin[key]; exist {
		err = fmt.Errorf("key %v exist, overwright", key)
	}
	Authkey2Admin[key] = [2][]string{
		[]string{"Login", "Admin"}, []string{"DelAfterLogin"},
	}
	return err
}

var allAuthorizationSet = map[string]*gos_authorize.AuthorizedCmds{
	"PreLogin": gos_authorize.NewByCmdIDList([]gos_idcmd.CommandID{
		gos_idcmd.Login,
	}),

	"DelAfterLogin": gos_authorize.NewByCmdIDList([]gos_idcmd.CommandID{
		gos_idcmd.Login,
	}),

	"Login": gos_authorize.NewByCmdIDList([]gos_idcmd.CommandID{
		gos_idcmd.Heartbeat,
		gos_idcmd.Chat,
		gos_idcmd.Act,
	}),
	"Admin": gos_authorize.NewByCmdIDList([]gos_idcmd.CommandID{}),
}

func NewPreLoginAuthorCmdIDList() *gos_authorize.AuthorizedCmds {
	return allAuthorizationSet["PreLogin"].Duplicate()
}

func UpdateByAuthKey(acicl *gos_authorize.AuthorizedCmds, key string) error {
	ag, exist := Authkey2Admin[key]
	if !exist {
		ag = [2][]string{[]string{"Login"}, []string{"DelAfterLogin"}}
	}
	// process include
	for _, authgroupname := range ag[0] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.Union(cmdidList)
	}
	// process exclude
	for _, authgroupname := range ag[1] {
		cmdidList := allAuthorizationSet[authgroupname]
		if cmdidList == nil {
			return fmt.Errorf("Can't Found authgroup %v", authgroupname)
		}
		acicl.SubIntersection(cmdidList)
	}
	return nil
}
