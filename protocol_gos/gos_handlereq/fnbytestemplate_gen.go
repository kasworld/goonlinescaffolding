// Code generated by "genprotocol -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_handlereq

/* bytes base fn map api template , unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error){
	gos_idcmd.Invalid: bytesAPIFn_ReqInvalid,// Invalid not used, make empty packet error
gos_idcmd.Login: bytesAPIFn_ReqLogin,// Login make session with nickname and enter stage
gos_idcmd.Heartbeat: bytesAPIFn_ReqHeartbeat,// Heartbeat prevent connection timeout
gos_idcmd.Chat: bytesAPIFn_ReqChat,// Chat chat to stage
gos_idcmd.Act: bytesAPIFn_ReqAct,// Act send user action

}   // DemuxReq2BytesAPIFnMap

	// Invalid not used, make empty packet error
	func bytesAPIFn_ReqInvalid(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error) {
		// robj, err := gos_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*gos_obj.ReqInvalid_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := gos_packet.Header{
			ErrorCode : gos_error.None,
		}
		sendBody := &gos_obj.RspInvalid_data{
		}
		return sendHeader, sendBody, nil
	}

	// Login make session with nickname and enter stage
	func bytesAPIFn_ReqLogin(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error) {
		// robj, err := gos_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*gos_obj.ReqLogin_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := gos_packet.Header{
			ErrorCode : gos_error.None,
		}
		sendBody := &gos_obj.RspLogin_data{
		}
		return sendHeader, sendBody, nil
	}

	// Heartbeat prevent connection timeout
	func bytesAPIFn_ReqHeartbeat(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error) {
		// robj, err := gos_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*gos_obj.ReqHeartbeat_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := gos_packet.Header{
			ErrorCode : gos_error.None,
		}
		sendBody := &gos_obj.RspHeartbeat_data{
		}
		return sendHeader, sendBody, nil
	}

	// Chat chat to stage
	func bytesAPIFn_ReqChat(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error) {
		// robj, err := gos_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*gos_obj.ReqChat_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := gos_packet.Header{
			ErrorCode : gos_error.None,
		}
		sendBody := &gos_obj.RspChat_data{
		}
		return sendHeader, sendBody, nil
	}

	// Act send user action
	func bytesAPIFn_ReqAct(
		me interface{}, hd gos_packet.Header, rbody []byte) (
		gos_packet.Header, interface{}, error) {
		// robj, err := gos_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*gos_obj.ReqAct_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := gos_packet.Header{
			ErrorCode : gos_error.None,
		}
		sendBody := &gos_obj.RspAct_data{
		}
		return sendHeader, sendBody, nil
	}

*/
