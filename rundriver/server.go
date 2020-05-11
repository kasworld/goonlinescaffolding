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

package main

import (
	"flag"
	"fmt"

	"github.com/kasworld/argdefault"
	"github.com/kasworld/configutil"
	"github.com/kasworld/go-profile"
	"github.com/kasworld/goonlinescaffolding/config/gameconst"
	"github.com/kasworld/goonlinescaffolding/config/serverconfig"
	"github.com/kasworld/goonlinescaffolding/game/server"
	"github.com/kasworld/goonlinescaffolding/lib/goslog"
	"github.com/kasworld/goonlinescaffolding/protocol_gos/gos_version"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/signalhandle"
	"github.com/kasworld/version"
)

var Ver = ""

func init() {
	version.Set(Ver)
}

func printVersion() {
	fmt.Println("goonlinescaffolding")
	fmt.Println("Build     ", version.GetVersion())
	fmt.Println("Data      ", gameconst.DataVersion)
	fmt.Println("Protocol  ", gos_version.ProtocolVersion)
	fmt.Println()
}

func main() {
	printVersion()

	configurl := flag.String("i", "", "server config file or url")
	signalhandle.AddArgs()
	profile.AddArgs()

	ads := argdefault.New(&serverconfig.Config{})
	ads.RegisterFlag()
	flag.Parse()
	config := &serverconfig.Config{
		LogLevel:      goslog.LL_All,
		SplitLogLevel: 0,
	}
	ads.SetDefaultToNonZeroField(config)
	if *configurl != "" {
		if err := configutil.LoadIni(*configurl, &config); err != nil {
			goslog.Fatal("%v", err)
		}
	}
	ads.ApplyFlagTo(config)
	if *configurl == "" {
		configutil.SaveIni("gosserver.ini", &config)
	}
	if profile.IsCpu() {
		fn := profile.StartCPUProfile()
		defer fn()
	}

	l, err := goslog.NewWithDstDir(
		"",
		config.MakeLogDir(),
		logflags.DefaultValue(false).BitClear(logflags.LF_functionname),
		config.LogLevel,
		config.SplitLogLevel,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	goslog.GlobalLogger = l

	svr := server.New(*config)
	if err := signalhandle.StartByArgs(svr); err != nil {
		goslog.Error("%v", err)
	}

	if profile.IsMem() {
		profile.WriteHeapProfile()
	}
}
