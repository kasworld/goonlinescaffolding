// Code generated by "genprotocol -ver=fa962a76ad7b14946f492eb8876e2f538e89415bc44d01f1655f1ad6b962a045 -basedir=protocol_gos -prefix=gos -statstype=int"

package gos_statcallapi

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idcmd"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_packet"
)

func (cps *StatCallAPI) String() string {
	return fmt.Sprintf("StatCallAPI[%v]",
		len(cps))
}

type StatCallAPI [gos_idcmd.CommandID_Count]StatRow

func New() *StatCallAPI {
	cps := new(StatCallAPI)
	for i := 0; i < gos_idcmd.CommandID_Count; i++ {
		cps[i].Name = gos_idcmd.CommandID(i).String()
	}
	return cps
}
func (cps *StatCallAPI) BeforeSendReq(header gos_packet.Header) (*statObj, error) {
	if int(header.Cmd) >= gos_idcmd.CommandID_Count {
		return nil, fmt.Errorf("CommandID out of range %v %v",
			header, gos_idcmd.CommandID_Count)
	}
	return cps[header.Cmd].open(), nil
}
func (cps *StatCallAPI) AfterSendReq(header gos_packet.Header) error {
	if int(header.Cmd) >= gos_idcmd.CommandID_Count {
		return fmt.Errorf("CommandID out of range %v %v", header, gos_idcmd.CommandID_Count)
	}
	n := int(header.BodyLen()) + gos_packet.HeaderLen
	cps[header.Cmd].addTx(n)
	return nil
}
func (cps *StatCallAPI) AfterRecvRsp(header gos_packet.Header) error {
	if int(header.Cmd) >= gos_idcmd.CommandID_Count {
		return fmt.Errorf("CommandID out of range %v %v", header, gos_idcmd.CommandID_Count)
	}
	n := int(header.BodyLen()) + gos_packet.HeaderLen
	cps[header.Cmd].addRx(n)
	return nil
}
func (ws *StatCallAPI) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Parse(`
	<html><head><title>Call API Stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table><br/>
	</body></html>`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, ws); err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
type statObj struct {
	StartTime time.Time
	StatRef   *StatRow
}

func (so *statObj) CallServerEnd(success bool) {
	so.StatRef.close(success, so.StartTime)
}

////////////////////////////////////////////////////////////////////////////////
type PacketID2StatObj struct {
	mutex sync.RWMutex
	stats map[uint32]*statObj
}

func NewPacketID2StatObj() *PacketID2StatObj {
	return &PacketID2StatObj{
		stats: make(map[uint32]*statObj),
	}
}
func (som *PacketID2StatObj) Add(pkid uint32, so *statObj) error {
	som.mutex.Lock()
	defer som.mutex.Unlock()
	if _, exist := som.stats[pkid]; exist {
		return fmt.Errorf("pkid exist %v", pkid)
	}
	som.stats[pkid] = so
	return nil
}
func (som *PacketID2StatObj) Del(pkid uint32) *statObj {
	som.mutex.Lock()
	defer som.mutex.Unlock()
	so := som.stats[pkid]
	delete(som.stats, pkid)
	return so
}
func (som *PacketID2StatObj) Get(pkid uint32) *statObj {
	som.mutex.RLock()
	defer som.mutex.RUnlock()
	return som.stats[pkid]
}

////////////////////////////////////////////////////////////////////////////////
const (
	HTML_tableheader = `<tr>
	<th>Name</th>
	<th>Start</th>
	<th>End</th>
	<th>Success</th>
	<th>Running</th>
	<th>Fail</th>
	<th>Avg ms</th>
	<th>TxAvg Byte</th>
	<th>RxAvg Byte</th>
	</tr>`
	HTML_row = `<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.StartCount}}</td>
	<td>{{$v.EndCount}}</td>
	<td>{{$v.SuccessCount}}</td>
	<td>{{$v.RunCount}}</td>
	<td>{{$v.FailCount}}</td>
	<td>{{printf "%13.6f" $v.Avgms }}</td>
	<td>{{printf "%10.3f" $v.AvgTx }}</td>
	<td>{{printf "%10.3f" $v.AvgRx }}</td>
	</tr>
	`
)

type StatRow struct {
	mutex        sync.Mutex
	Name         string
	TxCount      int
	TxByte       int
	RxCount      int
	RxByte       int
	StartCount   int
	EndCount     int
	SuccessCount int
	Sum          time.Duration
}

func (sr *StatRow) open() *statObj {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.StartCount++
	return &statObj{
		StartTime: time.Now(),
		StatRef:   sr,
	}
}
func (sr *StatRow) close(success bool, startTime time.Time) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.EndCount++
	if success {
		sr.SuccessCount++
		sr.Sum += time.Now().Sub(startTime)
	}
}
func (sr *StatRow) addTx(n int) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.TxCount++
	sr.TxByte += n
}
func (sr *StatRow) addRx(n int) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	sr.RxCount++
	sr.RxByte += n
}
func (sr *StatRow) RunCount() int {
	return sr.StartCount - sr.EndCount
}
func (sr *StatRow) FailCount() int {
	return sr.EndCount - sr.SuccessCount
}
func (sr *StatRow) Avgms() float64 {
	if sr.EndCount != 0 {
		return float64(sr.Sum) / float64(sr.EndCount*1000000)
	}
	return 0.0
}
func (sr *StatRow) AvgRx() float64 {
	if sr.EndCount != 0 {
		return float64(sr.RxByte) / float64(sr.RxCount)
	}
	return 0.0
}
func (sr *StatRow) AvgTx() float64 {
	if sr.EndCount != 0 {
		return float64(sr.TxByte) / float64(sr.TxCount)
	}
	return 0.0
}
