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

package server

import (
	"runtime"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goonlinescaffolding/config/dataversion"
	"github.com/kasworld/goonlinescaffolding/config/serverconfig"
	"github.com/kasworld/goonlinescaffolding/game/stagemanager"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_connbytemanager"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statapierror"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statnoti"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statserveapi"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_version"
	"github.com/kasworld/version"
	"github.com/kasworld/wrapper"
)

func (svr *Server) BuildDate() time.Time {
	return version.GetBuildDate()
}

func (svr *Server) GetVersion() string {
	return version.GetVersion()
}

func (svr *Server) GetProtocolVersion() string {
	return gos_version.ProtocolVersion
}

func (svr *Server) GetDataVersion() string {
	return dataversion.DataVersion
}

func (svr *Server) NumGoroutine() int {
	return runtime.NumGoroutine()
}

func (svr *Server) WrapInfo() string {
	return wrapper.G_WrapperInfo()
}

func (svr *Server) GetRunDur() time.Duration {
	return time.Now().Sub(svr.startTime)
}
func (svr *Server) GetStartTime() time.Time {
	return svr.startTime
}

func (svr *Server) GetSendStat() *actpersec.ActPerSec {
	return svr.SendStat
}
func (svr *Server) GetRecvStat() *actpersec.ActPerSec {
	return svr.RecvStat
}
func (svr *Server) GetProtocolStat() *gos_statserveapi.StatServeAPI {
	return svr.apiStat
}
func (svr *Server) GetNotiStat() *gos_statnoti.StatNotification {
	return svr.notiStat
}
func (svr *Server) GetErrorStat() *gos_statapierror.StatAPIError {
	return svr.errorStat
}
func (svr *Server) Config() serverconfig.Config {
	return svr.config
}

func (svr *Server) GetConnMan() *gos_connbytemanager.Manager {
	return svr.connManager
}

func (svr *Server) GetStageMan() *stagemanager.Manager {
	return svr.stageManager
}
