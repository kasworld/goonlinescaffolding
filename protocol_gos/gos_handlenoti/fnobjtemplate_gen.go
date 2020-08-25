// Code generated by "genprotocol -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_handlenoti

/* obj base demux fn map template

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd gos_packet.Header, body interface{}) error {
gos_idnoti.Invalid : objRecvNotiFn_Invalid, // Invalid not used, make empty packet error
gos_idnoti.StageInfo : objRecvNotiFn_StageInfo, // StageInfo for client display
gos_idnoti.StageChat : objRecvNotiFn_StageChat, // StageChat broadcasted chat

}

	// Invalid not used, make empty packet error
	func objRecvNotiFn_Invalid(me interface{}, hd gos_packet.Header, body interface{}) error {
		robj , ok := body.(*gos_obj.NotiInvalid_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	// StageInfo for client display
	func objRecvNotiFn_StageInfo(me interface{}, hd gos_packet.Header, body interface{}) error {
		robj , ok := body.(*gos_obj.NotiStageInfo_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	// StageChat broadcasted chat
	func objRecvNotiFn_StageChat(me interface{}, hd gos_packet.Header, body interface{}) error {
		robj , ok := body.(*gos_obj.NotiStageChat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

*/
