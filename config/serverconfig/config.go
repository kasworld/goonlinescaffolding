// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)

package serverconfig

import (
	"fmt"
	"path/filepath"

	"github.com/kasworld/goonlinescaffolding/lib/goslog"
	"github.com/kasworld/prettystring"
)

type Config struct {
	LogLevel              goslog.LL_Type `default:"7" argname:""`
	SplitLogLevel         goslog.LL_Type `default:"0" argname:""`
	BaseLogDir            string         `default:""  argname:""`
	ServerDataFolder      string         `default:"./serverdata" argname:""`
	ClientDataFolder      string         `default:"./clientdata" argname:""`
	ServicePort           string         `default:":24101"  argname:""`
	AdminPort             string         `default:":24201"  argname:""`
	ConcurrentConnections int            `default:"1000" argname:""`
	StageCount            int            `default:"10" argname:""`
	ActTurnPerSec         float64        `default:"30.0" argname:""`
}

func (config Config) MakeLogDir() string {
	rstr := filepath.Join(
		config.BaseLogDir,
		"gosserver.logfiles",
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) MakePIDFileFullpath() string {
	rstr := filepath.Join(
		config.BaseLogDir,
		"gosserver.pid",
	)
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) MakeOutfileFullpath() string {
	rstr := "gosserver.out"
	rtn, err := filepath.Abs(rstr)
	if err != nil {
		fmt.Println(rstr, rtn, err.Error())
		return rstr
	}
	return rtn
}

func (config Config) StringForm() string {
	return prettystring.PrettyString(config, 4)
}
