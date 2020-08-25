// Code generated by "genprotocol.exe -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_serveconnbyte

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_authorize"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_const"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idnoti"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_looptcp"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_loopwsgorilla"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statapierror"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statnoti"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_statserveapi"
)

type CounterI interface {
	Inc()
}

func (scb *ServeConnByte) String() string {
	return fmt.Sprintf("ServeConnByte[SendCh:%v/%v]",
		len(scb.sendCh), cap(scb.sendCh))
}

type ServeConnByte struct {
	connData       interface{} // custom data for this conn
	sendCh         chan gos_packet.Packet
	sendRecvStop   func()
	authorCmdList  *gos_authorize.AuthorizedCmds
	pid2ApiStatObj *gos_statserveapi.PacketID2StatObj
	apiStat        *gos_statserveapi.StatServeAPI
	notiStat       *gos_statnoti.StatNotification
	errorStat      *gos_statapierror.StatAPIError
	sendCounter    CounterI
	recvCounter    CounterI

	demuxReq2BytesAPIFnMap [gos_idcmd.CommandID_Count]func(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error)
}

// New with stats local
func New(
	connData interface{},
	sendBufferSize int,
	authorCmdList *gos_authorize.AuthorizedCmds,
	sendCounter, recvCounter CounterI,
	demuxReq2BytesAPIFnMap [gos_idcmd.CommandID_Count]func(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error),
) *ServeConnByte {
	scb := &ServeConnByte{
		connData:               connData,
		sendCh:                 make(chan gos_packet.Packet, sendBufferSize),
		pid2ApiStatObj:         gos_statserveapi.NewPacketID2StatObj(),
		apiStat:                gos_statserveapi.New(),
		notiStat:               gos_statnoti.New(),
		errorStat:              gos_statapierror.New(),
		sendCounter:            sendCounter,
		recvCounter:            recvCounter,
		authorCmdList:          authorCmdList,
		demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
	}
	scb.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", scb)
	}
	return scb
}

// NewWithStats with stats global
func NewWithStats(
	connData interface{},
	sendBufferSize int,
	authorCmdList *gos_authorize.AuthorizedCmds,
	sendCounter, recvCounter CounterI,
	apiStat *gos_statserveapi.StatServeAPI,
	notiStat *gos_statnoti.StatNotification,
	errorStat *gos_statapierror.StatAPIError,
	demuxReq2BytesAPIFnMap [gos_idcmd.CommandID_Count]func(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error),
) *ServeConnByte {
	scb := &ServeConnByte{
		connData:               connData,
		sendCh:                 make(chan gos_packet.Packet, sendBufferSize),
		pid2ApiStatObj:         gos_statserveapi.NewPacketID2StatObj(),
		apiStat:                apiStat,
		notiStat:               notiStat,
		errorStat:              errorStat,
		sendCounter:            sendCounter,
		recvCounter:            recvCounter,
		authorCmdList:          authorCmdList,
		demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
	}
	scb.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", scb)
	}
	return scb
}

func (scb *ServeConnByte) Disconnect() {
	scb.sendRecvStop()
}
func (scb *ServeConnByte) GetConnData() interface{} {
	return scb.connData
}
func (scb *ServeConnByte) GetAPIStat() *gos_statserveapi.StatServeAPI {
	return scb.apiStat
}
func (scb *ServeConnByte) GetNotiStat() *gos_statnoti.StatNotification {
	return scb.notiStat
}
func (scb *ServeConnByte) GetErrorStat() *gos_statapierror.StatAPIError {
	return scb.errorStat
}
func (scb *ServeConnByte) GetAuthorCmdList() *gos_authorize.AuthorizedCmds {
	return scb.authorCmdList
}
func (scb *ServeConnByte) StartServeWS(
	mainctx context.Context, conn *websocket.Conn,
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
) error {
	var returnerr error
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	scb.sendRecvStop = sendRecvCancel
	go func() {
		err := gos_loopwsgorilla.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
			readTimeoutSec, scb.handleRecvPacket)
		if err != nil {
			returnerr = fmt.Errorf("end RecvLoop %v", err)
		}
	}()
	go func() {
		err := gos_loopwsgorilla.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
			writeTimeoutSec, scb.sendCh,
			marshalfn, scb.handleSentPacket)
		if err != nil {
			returnerr = fmt.Errorf("end SendLoop %v", err)
		}
	}()
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		}
	}
	return returnerr
}
func (scb *ServeConnByte) StartServeTCP(
	mainctx context.Context, conn *net.TCPConn,
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
) error {
	var returnerr error
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	scb.sendRecvStop = sendRecvCancel
	go func() {
		err := gos_looptcp.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
			readTimeoutSec, scb.handleRecvPacket)
		if err != nil {
			returnerr = fmt.Errorf("end RecvLoop %v", err)
		}
	}()
	go func() {
		err := gos_looptcp.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
			writeTimeoutSec, scb.sendCh,
			marshalfn, scb.handleSentPacket)
		if err != nil {
			returnerr = fmt.Errorf("end SendLoop %v", err)
		}
	}()
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		}
	}
	return returnerr
}
func (scb *ServeConnByte) handleSentPacket(header gos_packet.Header) error {
	scb.sendCounter.Inc()
	switch header.FlowType {
	default:
		return fmt.Errorf("invalid packet type %s %v", scb, header)

	case gos_packet.Request:
		return fmt.Errorf("request packet not supported %s %v", scb, header)

	case gos_packet.Response:
		statOjb := scb.pid2ApiStatObj.Del(header.ID)
		if statOjb != nil {
			statOjb.AfterSendRsp(header)
		} else {
			return fmt.Errorf("send StatObj not found %v", header)
		}
	case gos_packet.Notification:
		scb.notiStat.Add(header)
	}
	return nil
}
func (scb *ServeConnByte) handleRecvPacket(rheader gos_packet.Header, rbody []byte) error {
	scb.recvCounter.Inc()
	if rheader.FlowType != gos_packet.Request {
		return fmt.Errorf("Unexpected rheader packet type: %v", rheader)
	}
	if int(rheader.Cmd) >= len(scb.demuxReq2BytesAPIFnMap) {
		return fmt.Errorf("Invalid rheader command %v", rheader)
	}
	if !scb.authorCmdList.CheckAuth(gos_idcmd.CommandID(rheader.Cmd)) {
		return fmt.Errorf("Not authorized packet %v", rheader)
	}

	statObj, err := scb.apiStat.AfterRecvReqHeader(rheader)
	if err != nil {
		return err
	}
	if err := scb.pid2ApiStatObj.Add(rheader.ID, statObj); err != nil {
		return err
	}
	statObj.BeforeAPICall()

	// timeout api call
	apiResult := scb.callAPI_timed(rheader, rbody)
	sheader, sbody, apierr := apiResult.header, apiResult.body, apiResult.err

	// no timeout api call
	//fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
	//sheader, sbody, apierr := fn(scb, rheader, rbody)

	statObj.AfterAPICall()

	scb.errorStat.Inc(gos_idcmd.CommandID(rheader.Cmd), sheader.ErrorCode)
	if apierr != nil {
		return apierr
	}
	if sbody == nil {
		return fmt.Errorf("Response body nil")
	}
	sheader.FlowType = gos_packet.Response
	sheader.Cmd = rheader.Cmd
	sheader.ID = rheader.ID
	rpk := gos_packet.Packet{
		Header: sheader,
		Body:   sbody,
	}
	return scb.EnqueueSendPacket(rpk)
}

type callAPIResult struct {
	header gos_packet.Header
	body   interface{}
	err    error
}

func (scb *ServeConnByte) callAPI_timed(rheader gos_packet.Header, rbody []byte) callAPIResult {
	rtnCh := make(chan callAPIResult, 1)
	go func(rtnCh chan callAPIResult, rheader gos_packet.Header, rbody []byte) {
		fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
		sheader, sbody, apierr := fn(scb, rheader, rbody)
		rtnCh <- callAPIResult{sheader, sbody, apierr}
	}(rtnCh, rheader, rbody)
	timeoutTk := time.NewTicker(gos_const.ServerAPICallTimeOutDur)
	defer timeoutTk.Stop()
	select {
	case apiResult := <-rtnCh:
		return apiResult
	case <-timeoutTk.C:
		return callAPIResult{rheader, nil, fmt.Errorf("APICall Timeout %v", rheader)}
	}
}
func (scb *ServeConnByte) EnqueueSendPacket(pk gos_packet.Packet) error {
	select {
	case scb.sendCh <- pk:
		return nil
	default:
		return fmt.Errorf("Send channel full %v", scb)
	}
}
func (scb *ServeConnByte) SendNotiPacket(
	cmd gos_idnoti.NotiID, body interface{}) error {
	err := scb.EnqueueSendPacket(gos_packet.Packet{
		gos_packet.Header{
			Cmd:      uint16(cmd),
			FlowType: gos_packet.Notification,
		},
		body,
	})
	if err != nil {
		scb.Disconnect()
	}
	return err
}
