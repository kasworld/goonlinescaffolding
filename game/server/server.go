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
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/kasworld/actpersec"
	"github.com/kasworld/goonlinescaffolding/config/serverconfig"
	"github.com/kasworld/goonlinescaffolding/game/stage"
	"github.com/kasworld/goonlinescaffolding/game/stagemanager"
	"github.com/kasworld/goonlinescaffolding/lib/goslog"
	"github.com/kasworld/goonlinescaffolding/lib/sessionmanager"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_connbytemanager"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statapierror"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statnoti"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statserveapi"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/weblib/retrylistenandserve"
)

type Server struct {
	rnd       *rand.Rand      `prettystring:"hide"`
	log       *goslog.LogBase `prettystring:"hide"`
	config    serverconfig.Config
	adminWeb  *http.Server `prettystring:"simple"`
	clientWeb *http.Server `prettystring:"simple"`
	startTime time.Time    `prettystring:"simple"`

	sendRecvStop func()
	SendStat     *actpersec.ActPerSec `prettystring:"simple"`
	RecvStat     *actpersec.ActPerSec `prettystring:"simple"`

	apiStat   *gos_statserveapi.StatServeAPI
	notiStat  *gos_statnoti.StatNotification
	errorStat *gos_statapierror.StatAPIError

	marshalBodyFn          func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error)
	unmarshalPacketFn      func(h gos_packet.Header, bodyData []byte) (interface{}, error)
	DemuxReq2BytesAPIFnMap [gos_idcmd.CommandID_Count]func(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error)

	connManager    *gos_connbytemanager.Manager
	sessionManager *sessionmanager.SessionManager

	stageManager *stagemanager.Manager
}

func New(config serverconfig.Config) *Server {
	if config.BaseLogDir != "" {
		log, err := goslog.NewWithDstDir(
			"",
			config.MakeLogDir(),
			logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
			config.LogLevel,
			config.SplitLogLevel,
		)
		if err == nil {
			goslog.GlobalLogger = log
		} else {
			fmt.Printf("%v\n", err)
			goslog.GlobalLogger.SetFlags(
				goslog.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
			goslog.GlobalLogger.SetLevel(
				config.LogLevel)
		}
	} else {
		goslog.GlobalLogger.SetFlags(
			goslog.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
		goslog.GlobalLogger.SetLevel(
			config.LogLevel)
	}
	svr := &Server{
		config: config,
		log:    goslog.GlobalLogger,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),

		SendStat: actpersec.New(),
		RecvStat: actpersec.New(),

		apiStat:        gos_statserveapi.New(),
		notiStat:       gos_statnoti.New(),
		errorStat:      gos_statapierror.New(),
		connManager:    gos_connbytemanager.New(),
		sessionManager: sessionmanager.New("", config.ConcurrentConnections, goslog.GlobalLogger),

		stageManager: stagemanager.New(goslog.GlobalLogger),
	}
	svr.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call\n")
	}
	return svr
}

// called from signal handler
func (svr *Server) GetServiceLockFilename() string {
	return svr.config.MakePIDFileFullpath()
}

// called from signal handler
// return implement signalhandle.LoggerI
func (svr *Server) GetLogger() interface{} {
	return goslog.GlobalLogger
}

// called from signal handler
func (svr *Server) ServiceInit() error {
	return nil
}

// called from signal handler
func (svr *Server) ServiceCleanup() {
}

// called from signal handler
func (svr *Server) ServiceMain(mainctx context.Context) {
	fmt.Println(prettystring.PrettyString(svr.config, 4))
	svr.startTime = time.Now()

	ctx, stopFn := context.WithCancel(mainctx)
	svr.sendRecvStop = stopFn
	defer svr.sendRecvStop()

	svr.initAdminWeb()
	svr.initServiceWeb(ctx)

	fmt.Printf("WebAdmin  : %v:%v id:%v pass:%v\n",
		svr.config.ServiceHostBase, svr.config.AdminPort, svr.config.WebAdminID, svr.config.WebAdminPass)
	fmt.Printf("WebClient : %v:%v/\n", svr.config.ServiceHostBase, svr.config.ServicePort)

	go retrylistenandserve.RetryListenAndServe(svr.adminWeb, svr.log, "serveAdminWeb")
	go retrylistenandserve.RetryListenAndServe(svr.clientWeb, svr.log, "serveServiceWeb")

	for i := 0; i < svr.config.StageCount; i++ {
		stg := stage.New(svr.log, svr.config)
		svr.stageManager.Add(stg)
		go stg.Run(ctx)
	}

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerInfoTk.C:
			svr.SendStat.UpdateLap()
			svr.RecvStat.UpdateLap()
		}
	}
}
