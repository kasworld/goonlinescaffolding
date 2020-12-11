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

package wasmclient

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"syscall/js"
	"time"

	"github.com/kasworld/goonlinescaffolding/lib/clientcookie"
	"github.com/kasworld/goonlinescaffolding/lib/jsobj"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_connwasm"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_gob"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_obj"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_pid2rspfn"
	"github.com/kasworld/gowasmlib/jslog"
	"github.com/kasworld/gowasmlib/wasmcookie"
)

func getConnURL() string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}
	u.Path = "ws"
	u.Scheme = "ws"
	return u.String()
}

func (app *WasmClient) NetInit(ctx context.Context) (*gos_obj.RspLogin_data, error) {
	app.wsConn = gos_connwasm.New(
		getConnURL(),
		gos_gob.MarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket)

	fmt.Println(getConnURL())

	var wg sync.WaitGroup

	// connect
	wg.Add(1)
	go func() {
		err := app.wsConn.Connect(ctx, &wg)
		if err != nil {
			jslog.Errorf("wsConn.Connect err %v", err)
			app.DoClose()
		}
	}()
	authkey := clientcookie.GetQuery().Get("authkey")
	nick := jsobj.GetTextValueFromInputText("nickname")
	ck := wasmcookie.GetMap()
	sessionkey := ck[clientcookie.SessionKeyName()]
	wg.Wait()
	jslog.Info("connected")

	// login
	var rtn *gos_obj.RspLogin_data
	wg.Add(1)
	app.ReqWithRspFn(
		gos_idcmd.Login,
		&gos_obj.ReqLogin_data{
			SessionKey: sessionkey,
			NickName:   nick,
			AuthKey:    authkey,
		},
		func(hd gos_packet.Header, rsp interface{}) error {
			rtn = rsp.(*gos_obj.RspLogin_data)
			wg.Done()
			return nil
		},
	)
	wg.Wait()
	jslog.Info("logined")

	return rtn, nil
}

func (app *WasmClient) Cleanup() {
	app.wsConn.SendRecvStop()
}

func (app *WasmClient) handleSentPacket(pk *gos_packet.Packet) error {
	return nil
}

func (app *WasmClient) handleRecvPacket(header gos_packet.Header, body []byte) error {
	robj, err := gos_gob.UnmarshalPacket(header, body)
	if err != nil {
		return err
	}
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, robj)
	case gos_packet.Response:
		if err := app.pid2recv.HandleRsp(header, robj); err != nil {
			return err
		}
	case gos_packet.Notification:
		fn := DemuxNoti2ObjFnMap[header.Cmd]
		if err := fn(app, header, robj); err != nil {
			return err
		}
	}
	return nil
}

func (app *WasmClient) ReqWithRspFn(cmd gos_idcmd.CommandID, body interface{},
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
	if err := app.wsConn.EnqueueSendPacket(&spk); err != nil {
		app.wsConn.SendRecvStop()
		return fmt.Errorf("Send fail %s %v:%v %v", app, cmd, pid, err)
	}
	return nil
}

func (app *WasmClient) reqHeartbeat() error {
	return app.ReqWithRspFn(
		gos_idcmd.Heartbeat,
		&gos_obj.ReqHeartbeat_data{
			Tick: time.Now().UnixNano(),
		},
		func(hd gos_packet.Header, rsp interface{}) error {
			rpk := rsp.(*gos_obj.RspHeartbeat_data)
			pingDur := time.Now().UnixNano() - rpk.Tick
			app.PingDur = (app.PingDur + pingDur) / 2
			return nil
		},
	)
}

func (app *WasmClient) ReqWithRspFnWithAuth(cmd gos_idcmd.CommandID, body interface{},
	fn gos_pid2rspfn.HandleRspFn) error {
	if !app.CanUseCmd(cmd) {
		return fmt.Errorf("Cmd not allowed %v", cmd)
	}
	return app.ReqWithRspFn(cmd, body, fn)
}

func (app *WasmClient) CanUseCmd(cmd gos_idcmd.CommandID) bool {
	if app.loginData == nil {
		return false
	}
	return app.loginData.CmdList[cmd]
}

func (app *WasmClient) sendPacket(cmd gos_idcmd.CommandID, arg interface{}) {
	app.ReqWithRspFnWithAuth(
		cmd, arg,
		func(hd gos_packet.Header, rsp interface{}) error {
			return nil
		},
	)
}
