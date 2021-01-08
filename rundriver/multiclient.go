// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	"github.com/kasworld/goonlinescaffolding/lib/goslog"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_connwsgorilla"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_error"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_gob"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_obj"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_pid2rspfn"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statapierror"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statcallapi"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statnoti"
	"github.com/kasworld/multirun"
	"github.com/kasworld/prettystring"
	"github.com/kasworld/rangestat"
)

// service const
const (
	// for client
	readTimeoutSec  = 6 * time.Second
	writeTimeoutSec = 3 * time.Second
)

func main() {
	configurl := flag.String("i", "", "client config file or url")
	ads := argdefault.New(&MultiClientConfig{})
	ads.RegisterFlag()
	flag.Parse()
	config := &MultiClientConfig{}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, config); err != nil {
			goslog.Error("%v", err)
		}
	}
	ads.ApplyFlagTo(config)
	fmt.Println(prettystring.PrettyString(config, 4))

	chErr := make(chan error)
	go multirun.Run(
		context.Background(),
		config.Concurrent,
		config.AccountPool,
		config.AccountOverlap,
		config.LimitStartCount,
		config.LimitEndCount,
		func(config interface{}) multirun.ClientI {
			return NewApp(config.(AppArg), goslog.GlobalLogger)
		},
		func(i int) interface{} {
			return AppArg{
				ConnectToServer: config.ConnectToServer,
				Nickname:        fmt.Sprintf("%v_%v", config.PlayerNameBase, i),
				SessionUUID:     "",
				Auth:            "",
			}
		},
		chErr,
		rangestat.New("", 0, config.Concurrent),
	)
	for err := range chErr {
		fmt.Printf("%v\n", err)
	}
}

type MultiClientConfig struct {
	ConnectToServer string `default:"localhost:24101" argname:""`
	PlayerNameBase  string `default:"MC_" argname:""`
	Concurrent      int    `default:"1000" argname:""`
	AccountPool     int    `default:"0" argname:""`
	AccountOverlap  int    `default:"0" argname:""`
	LimitStartCount int    `default:"0" argname:""`
	LimitEndCount   int    `default:"0" argname:""`
}

type AppArg struct {
	ConnectToServer string
	Nickname        string
	SessionUUID     string
	Auth            string
}

type App struct {
	config            AppArg
	c2scWS            *gos_connwsgorilla.Connection
	EnqueueSendPacket func(pk *gos_packet.Packet) error
	runResult         error

	sendRecvStop func()
	apistat      *gos_statcallapi.StatCallAPI
	pid2statobj  *gos_statcallapi.PacketID2StatObj
	notistat     *gos_statnoti.StatNotification
	errstat      *gos_statapierror.StatAPIError
	pid2recv     *gos_pid2rspfn.PID2RspFn
}

func NewApp(config AppArg, log *goslog.LogBase) *App {
	app := &App{
		config:      config,
		apistat:     gos_statcallapi.New(),
		pid2statobj: gos_statcallapi.NewPacketID2StatObj(),
		notistat:    gos_statnoti.New(),
		errstat:     gos_statapierror.New(),
		pid2recv:    gos_pid2rspfn.New(),
	}
	return app
}

func (app *App) String() string {
	return fmt.Sprintf("App[%v %v]", app.config.Nickname, app.config.SessionUUID)
}

func (app *App) GetArg() interface{} {
	return app.config
}

func (app *App) GetRunResult() error {
	return app.runResult
}

func (app *App) Run(mainctx context.Context) {
	ctx, stopFn := context.WithCancel(mainctx)
	app.sendRecvStop = stopFn
	defer app.sendRecvStop()

	app.c2scWS = gos_connwsgorilla.New(10)
	if err := app.c2scWS.ConnectTo(app.config.ConnectToServer); err != nil {
		fmt.Printf("%v\n", err)
		app.sendRecvStop()
		app.runResult = err
		return
	}
	app.EnqueueSendPacket = app.c2scWS.EnqueueSendPacket
	go func(ctx context.Context) {
		app.runResult = app.c2scWS.Run(ctx,
			readTimeoutSec, writeTimeoutSec,
			gos_gob.MarshalBodyFn,
			app.handleRecvPacket,
			app.handleSentPacket,
		)
	}(ctx)

	// login
	var wg sync.WaitGroup
	var rtn *gos_obj.RspLogin_data
	wg.Add(1)
	app.ReqWithRspFn(
		gos_idcmd.Login,
		&gos_obj.ReqLogin_data{
			SessionKey: app.config.SessionUUID,
			NickName:   app.config.Nickname,
			AuthKey:    "",
		},
		func(hd gos_packet.Header, rsp interface{}) error {

			rtn = rsp.(*gos_obj.RspLogin_data)
			wg.Done()
			return nil
		},
	)
	wg.Wait()

	timerPingTk := time.NewTicker(time.Second)
	defer timerPingTk.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-timerPingTk.C:
			go app.reqHeartbeat()

		}
	}
}

func (app *App) reqHeartbeat() error {
	return app.ReqWithRspFn(
		gos_idcmd.Heartbeat,
		&gos_obj.ReqHeartbeat_data{
			Tick: time.Now().UnixNano(),
		},
		func(hd gos_packet.Header, rsp interface{}) error {
			// rpk := rsp.(*gos_obj.RspHeartbeat_data)
			// pingDur := time.Now().UnixNano() - rpk.Tick
			// app.PingDur = (app.PingDur + pingDur) / 2
			return nil
		},
	)
}

func (app *App) handleSentPacket(pk *gos_packet.Packet) error {
	if err := app.apistat.AfterSendReq(pk.Header); err != nil {
		return err
	}
	return nil
}

func (app *App) handleRecvPacket(header gos_packet.Header, body []byte) error {
	robj, err := gos_gob.UnmarshalPacket(header, body)
	if err != nil {
		return err
	}

	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, body)
	case gos_packet.Notification:
		// noti stat
		app.notistat.Add(header)
		//process noti here
		// robj, err := gos_gob.UnmarshalPacket(header, body)

	case gos_packet.Response:
		// error stat
		app.errstat.Inc(gos_idcmd.CommandID(header.Cmd), header.ErrorCode)
		// api stat
		if err := app.apistat.AfterRecvRsp(header); err != nil {
			fmt.Printf("%v %v\n", app, err)
			return err
		}
		psobj := app.pid2statobj.Get(header.ID)
		if psobj == nil {
			return fmt.Errorf("no statobj for %v", header.ID)
		}
		psobj.CallServerEnd(header.ErrorCode == gos_error.None)
		app.pid2statobj.Del(header.ID)

		// process response
		if err := app.pid2recv.HandleRsp(header, robj); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) ReqWithRspFn(cmd gos_idcmd.CommandID, body interface{},
	fn gos_pid2rspfn.HandleRspFn) error {

	pid := app.pid2recv.NewPID(fn)
	spk := gos_packet.Packet{
		Header: gos_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: gos_packet.Request,
		},
		Body: body,
	}

	// add api stat
	psobj, err := app.apistat.BeforeSendReq(spk.Header)
	if err != nil {
		return nil
	}
	app.pid2statobj.Add(spk.Header.ID, psobj)

	if err := app.EnqueueSendPacket(&spk); err != nil {
		fmt.Printf("End %v %v %v\n", app, spk, err)
		app.sendRecvStop()
		return fmt.Errorf("Send fail %v %v", app, err)
	}
	return nil
}
