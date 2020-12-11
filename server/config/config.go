package config

import (
	"sharpshooterTunnel/common"
)

var CFG *config

type config struct {
	C          string `json:"c" flag:"c" default:"" usage:""`
	LocalAddr  string `json:"local_addr" flag:"local_addr" usage:""`
	Key        string `json:"key" flag:"key" usage:"" default:"sharpshooter"`
	SendWin    int    `json:"send_win" flag:"send_win" default:"1024" usage:""`
	RecWin     int    `json:"rec_win" flag:"rec_win" default:"1024" usage:""`
	Interval   int64    `json:"interval" flag:"interval" default:"100" usage:""`
	ListenPort int    `json:"listen_port" flag:"listen_port" default:"0" usage:""`
	Debug      bool   `json:"debug" flag:"debug" default:"false" usage:""`
	PPort      int    `json:"p_port" flag:"p_port" default:"8888" usage:""`
}

func init() {
	c := &config{}
	common.AutoParse(c)
	common.GenConf(c)
	if c.C != "" {
		common.Load(c.C, c)
	}
	CFG = c
}
