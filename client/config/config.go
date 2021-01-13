package config

import (
	"sharpshooterTunnel/common"
)

var CFG *config

type config struct {
	C          string `json:"c,omitempty" flag:"c" default:"" usage:""`
	RemoteAddr string `json:"remote_addr" flag:"remote_addr" usage:""`
	LocalAddr  string `json:"local_addr" flag:"local_addr" usage:""`
	Key        string `json:"key" flag:"key" usage:"" default:"sharpshooter"`
	ConNum     int    `json:"con_num" flag:"con_num" usage:"" default:"10"`
	SendWin    int    `json:"send_win" flag:"send_win" default:"1024" usage:""`
	RecWin     int    `json:"rec_win" flag:"rec_win" default:"1024" usage:""`
	Interval   int64  `json:"interval" flag:"interval" default:"100" usage:""`
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
