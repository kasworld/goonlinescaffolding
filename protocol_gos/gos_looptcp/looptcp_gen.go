// Code generated by "genprotocol.exe -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_looptcp

import (
	"context"
	"net"
	"time"

	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_const"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
)

var bufPool = gos_packet.NewPool(gos_const.PacketBufferPoolSize)

func SendPacket(conn *net.TCPConn, buf []byte) error {
	toWrite := len(buf)
	for l := 0; l < toWrite; {
		n, err := conn.Write(buf[l:toWrite])
		if err != nil {
			return err
		}
		l += n
	}
	return nil
}

func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	SendCh chan gos_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(header gos_packet.Header) error,
) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		case pk := <-SendCh:
			if err = tcpConn.SetWriteDeadline(time.Now().Add(timeOut)); err != nil {
				break loop
			}
			oldbuf := bufPool.Get()
			sendBuffer, err := gos_packet.Packet2Bytes(&pk, marshalBodyFn, oldbuf)
			if err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = SendPacket(tcpConn, sendBuffer); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = handleSentPacketFn(pk.Header); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			bufPool.Put(oldbuf)
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	HandleRecvPacketFn func(header gos_packet.Header, body []byte) error,
) error {

	defer SendRecvStop()

	pb := gos_packet.NewRecvPacketBuffer()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			return nil

		default:
			if pb.IsPacketComplete() {
				header, rbody, lerr := pb.GetHeaderBody()
				if lerr != nil {
					err = lerr
					break loop
				}
				if err = HandleRecvPacketFn(header, rbody); err != nil {
					break loop
				}
				pb = gos_packet.NewRecvPacketBuffer()
				if err = tcpConn.SetReadDeadline(time.Now().Add(timeOut)); err != nil {
					break loop
				}
			} else {
				err := pb.Read(tcpConn)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}
