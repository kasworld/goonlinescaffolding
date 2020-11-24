// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gos_obj

import "github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"

// Invalid not used, make empty packet error
type ReqInvalid_data struct {
	Dummy uint8 // change as you need
}

// Invalid not used, make empty packet error
type RspInvalid_data struct {
	Dummy uint8 // change as you need
}

type ReqLogin_data struct {
	SessionKey   string
	NickName     string
	AuthKey      string
	StageToEnter string // stage number or uuid, empty or unknown to random
}
type RspLogin_data struct {
	Version         string
	ProtocolVersion string
	DataVersion     string

	SessionKey string
	StageUUID  string
	NickName   string
	CmdList    [gos_idcmd.CommandID_Count]bool
}

// Heartbeat prevent connection timeout
type ReqHeartbeat_data struct {
	Tick int64
}

// Heartbeat prevent connection timeout
type RspHeartbeat_data struct {
	Tick int64
}

// Chat chat to stage
type ReqChat_data struct {
	Chat string
}

// Chat chat to stage
type RspChat_data struct {
	Dummy uint8
}

// Act send user action
type ReqAct_data struct {
	Dummy uint8 // change as you need
}

// Act send user action
type RspAct_data struct {
	Dummy uint8 // change as you need
}

//////////////////////////////////////////////////////////////////////////////

type NotiInvalid_data struct {
	Dummy uint8
}

type NotiStageInfo_data struct {
	Tick int64
}

type NotiStageChat_data struct {
	SenderNick string
	Chat       string
}
