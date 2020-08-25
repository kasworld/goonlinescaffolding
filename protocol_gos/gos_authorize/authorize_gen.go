// Code generated by "genprotocol -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_authorize

import (
	"bytes"
	"fmt"

	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
)

type AuthorizedCmds [gos_idcmd.CommandID_Count]bool

func (acidl *AuthorizedCmds) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "AuthorizedCmds[")
	for i, v := range acidl {
		if v {
			fmt.Fprintf(&buff, "%v ", gos_idcmd.CommandID(i))
		}
	}
	fmt.Fprintf(&buff, "]")
	return buff.String()
}

func NewAllSet() *AuthorizedCmds {
	rtn := new(AuthorizedCmds)
	for i := 0; i < gos_idcmd.CommandID_Count; i++ {
		rtn[i] = true
	}
	return rtn
}

func NewByCmdIDList(cmdlist []gos_idcmd.CommandID) *AuthorizedCmds {
	rtn := new(AuthorizedCmds)
	for _, id := range cmdlist {
		rtn[id] = true
	}
	return rtn
}

func (acidl *AuthorizedCmds) Union(src *AuthorizedCmds) *AuthorizedCmds {
	for cmdid, auth := range src {
		if auth {
			acidl[cmdid] = true
		}
	}
	return acidl
}

func (acidl *AuthorizedCmds) SubIntersection(src *AuthorizedCmds) *AuthorizedCmds {
	for cmdid, auth := range src {
		if auth {
			acidl[cmdid] = false
		}
	}
	return acidl
}

func (acidl *AuthorizedCmds) Duplicate() *AuthorizedCmds {
	rtn := *acidl
	return &rtn
}

func (acidl *AuthorizedCmds) CheckAuth(cmdid gos_idcmd.CommandID) bool {
	return acidl[cmdid]
}
