// Copyright 2015,2016,2017,2018,2019,2020,2021 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stage

import (
	"context"
	"math/rand"
	"time"

	"github.com/kasworld/goonlinescaffolding/config/serverconfig"
	"github.com/kasworld/goonlinescaffolding/lib/goslog"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_connbytemanager"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_idnoti"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_obj"
	"github.com/kasworld/uuidstr"
)

type Stage struct {
	rnd    *rand.Rand      `prettystring:"hide"`
	log    *goslog.LogBase `prettystring:"hide"`
	config serverconfig.Config

	UUID  string
	Conns *gos_connbytemanager.Manager
}

func New(l *goslog.LogBase, config serverconfig.Config) *Stage {
	stg := &Stage{
		UUID:   uuidstr.New(),
		config: config,
		log:    l,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Conns:  gos_connbytemanager.New(),
	}

	return stg
}

func (stg *Stage) Run(ctx context.Context) {

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	turnDur := time.Duration(float64(time.Second) / stg.config.ActTurnPerSec)
	timerTurnTk := time.NewTicker(turnDur)
	defer timerTurnTk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerInfoTk.C:
		case <-timerTurnTk.C:
			si := stg.ToStageInfo()
			conlist := stg.Conns.GetList()
			for _, v := range conlist {
				v.SendNotiPacket(gos_idnoti.StageInfo,
					si,
				)
			}
		}
	}
}

func (stg *Stage) ToStageInfo() *gos_obj.NotiStageInfo_data {
	rtn := &gos_obj.NotiStageInfo_data{
		Tick: time.Now().UnixNano(),
	}
	return rtn
}
