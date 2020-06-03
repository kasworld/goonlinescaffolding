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
	"bytes"
	"fmt"
	"net/url"
	"syscall/js"

	"github.com/kasworld/goonlinescaffolding/config/gameconst"
	"github.com/kasworld/goonlinescaffolding/game/stagelist4client"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_version"
	"github.com/kasworld/gowasmlib/jslog"
)

func (app *WasmClient) makeButtons() string {
	var buf bytes.Buffer
	gameOptions.MakeButtonToolTipTop(&buf)
	return buf.String()
}

func (app *WasmClient) DisplayTextInfo() {
	app.updateLeftInfo()
	app.updateRightInfo()
	app.updateCenterInfo()
}

func (app *WasmClient) makeServiceInfo() string {
	msgCopyright := `</hr>Copyright 2019,2020 SeukWon Kang 
		<a href="https://github.com/kasworld/goonlinescaffolding" target="_blank">goonlinescaffolding</a>`

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "goonlinescaffolding webclient<br/>")
	fmt.Fprintf(&buf, "Protocol %v<br/>", gos_version.ProtocolVersion)
	fmt.Fprintf(&buf, "Data %v<br/>", gameconst.DataVersion)
	fmt.Fprintf(&buf, "%v<br/>", msgCopyright)
	return buf.String()
}

func (app *WasmClient) makeDebugInfo() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf,
		"%v<br/>Ping %v<br/>ServerClientTickDiff %v<br/>",
		app.DispInterDur, app.PingDur, app.ServerClientTictDiff,
	)
	return buf.String()
}

func loadStageListHTML() string {
	tlurl := ReplacePathFromHref("stagelist.json")
	aol, err := stagelist4client.LoadFromURL(tlurl)
	if err != nil {
		jslog.Errorf("stagelist load fail %v", err)
		return "fail to load stagelist"
	}
	var buf bytes.Buffer
	buf.WriteString(`
		stage list in server
		<table border=1 style="border-collapse:collapse;">
		<tr>
		<th>Number</th> <th>UUID</th> <th>Info</th> <th>Command</th> 
		</tr>	
		`)
	for i, stg := range aol {
		fmt.Fprintf(&buf, `
		<tr>
		<td>%v</td> <td>%v</td> <td>%v</td> 
		<td><button type="button" style="font-size:20px;" onclick="enterStage('%v')">Enter Stage</button></td> 
		</tr>`,
			i, stg.UUID, stg.Info, stg.UUID,
		)
	}
	buf.WriteString(`
		<tr>
		<th>Number</th> <th>UUID</th> <th>Type</th> <th>Command</th> 
		</tr>
		</table>
		`)
	return buf.String()
}

func ReplacePathFromHref(s string) string {
	loc := js.Global().Get("window").Get("location").Get("href")
	u, err := url.Parse(loc.String())
	if err != nil {
		jslog.Errorf("%v", err)
		return ""
	}
	u.Path = s
	return u.String()
}
