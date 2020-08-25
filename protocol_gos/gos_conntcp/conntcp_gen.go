// Code generated by "genprotocol.exe -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_conntcp

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_looptcp"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
)

type Connection struct {
	conn         *net.TCPConn
	sendCh       chan gos_packet.Packet
	sendRecvStop func()

	readTimeoutSec     time.Duration
	writeTimeoutSec    time.Duration
	marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
	handleRecvPacketFn func(header gos_packet.Header, body []byte) error
	handleSentPacketFn func(header gos_packet.Header) error
}

func New(
	readTimeoutSec, writeTimeoutSec time.Duration,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleRecvPacketFn func(header gos_packet.Header, body []byte) error,
	handleSentPacketFn func(header gos_packet.Header) error,
) *Connection {
	tc := &Connection{
		sendCh:             make(chan gos_packet.Packet, 10),
		readTimeoutSec:     readTimeoutSec,
		writeTimeoutSec:    writeTimeoutSec,
		marshalBodyFn:      marshalBodyFn,
		handleRecvPacketFn: handleRecvPacketFn,
		handleSentPacketFn: handleSentPacketFn,
	}

	tc.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", tc)
	}
	return tc
}

func (tc *Connection) ConnectTo(remoteAddr string) error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return err
	}
	tc.conn, err = net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		return err
	}
	return nil
}

func (tc *Connection) Cleanup() {
	tc.sendRecvStop()
	if tc.conn != nil {
		tc.conn.Close()
	}
}

func (tc *Connection) Run(mainctx context.Context) error {
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	tc.sendRecvStop = sendRecvCancel
	var rtnerr error
	var sendRecvWaitGroup sync.WaitGroup
	sendRecvWaitGroup.Add(2)
	go func() {
		defer sendRecvWaitGroup.Done()
		err := gos_looptcp.RecvLoop(
			sendRecvCtx,
			tc.sendRecvStop,
			tc.conn,
			tc.readTimeoutSec,
			tc.handleRecvPacketFn)
		if err != nil {
			rtnerr = err
		}
	}()
	go func() {
		defer sendRecvWaitGroup.Done()
		err := gos_looptcp.SendLoop(
			sendRecvCtx,
			tc.sendRecvStop,
			tc.conn,
			tc.writeTimeoutSec,
			tc.sendCh,
			tc.marshalBodyFn,
			tc.handleSentPacketFn)
		if err != nil {
			rtnerr = err
		}
	}()
	sendRecvWaitGroup.Wait()
	return rtnerr
}

func (tc *Connection) EnqueueSendPacket(pk gos_packet.Packet) error {
	select {
	case tc.sendCh <- pk:
		return nil
	default:
		return fmt.Errorf("Send channel full %v", tc)
	}
}
