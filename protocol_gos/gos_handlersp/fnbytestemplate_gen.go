// Code generated by "genprotocol -ver=ee0c2a9542c23df14c7723030f3c8ad53da871d9e014d150b2e62d1ce8399fde -basedir=. -prefix=gos -statstype=int"

package gos_handlersp

/* bytes base demux fn map template

var DemuxRsp2BytesFnMap = [...]func(me interface{}, hd gos_packet.Header, rbody []byte) error {
gos_idcmd.Invalid : bytesRecvRspFn_Invalid, // Invalid not used, make empty packet error
gos_idcmd.Login : bytesRecvRspFn_Login, // Login make session with nickname and enter stage
gos_idcmd.Heartbeat : bytesRecvRspFn_Heartbeat, // Heartbeat prevent connection timeout
gos_idcmd.Chat : bytesRecvRspFn_Chat, // Chat chat to stage
gos_idcmd.Act : bytesRecvRspFn_Act, // Act send user action

}

	// Invalid not used, make empty packet error
	func bytesRecvRspFn_Invalid(me interface{}, hd gos_packet.Header, rbody []byte) error {
		robj, err := gos_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*gos_obj.RspInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	// Login make session with nickname and enter stage
	func bytesRecvRspFn_Login(me interface{}, hd gos_packet.Header, rbody []byte) error {
		robj, err := gos_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*gos_obj.RspLogin_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	// Heartbeat prevent connection timeout
	func bytesRecvRspFn_Heartbeat(me interface{}, hd gos_packet.Header, rbody []byte) error {
		robj, err := gos_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*gos_obj.RspHeartbeat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	// Chat chat to stage
	func bytesRecvRspFn_Chat(me interface{}, hd gos_packet.Header, rbody []byte) error {
		robj, err := gos_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*gos_obj.RspChat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

	// Act send user action
	func bytesRecvRspFn_Act(me interface{}, hd gos_packet.Header, rbody []byte) error {
		robj, err := gos_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*gos_obj.RspAct_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

*/
