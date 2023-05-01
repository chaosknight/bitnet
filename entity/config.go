package entity

import (
	"bytes"
	"encoding/json"
	"os"
)

type FinesseConfig struct {
	InstID string `json:"instId"`
	//底仓
	Firstcount float64 `json:"firstcount"`
	//最大仓位
	Maxcount float64 `json:"maxcount"`
	//加仓条件
	Addposition float64 `json:"addposition"`
	//最大亏损 触发平仓
	Maxloss float64 `json:"maxloss"`
}

type Sysconfig struct {
	User [3]string `json:"user"`

	Insts []*FinesseConfig `json:"insts"`
}

func LoadCnf(fname string) *Sysconfig {
	f, err := os.Open(fname)
	if err != nil {
		return nil
	}
	defer f.Close()
	b := new(bytes.Buffer)
	_, err = b.ReadFrom(f)
	if err != nil {
		return nil
	}
	sf := Sysconfig{}
	err = json.Unmarshal(b.Bytes(), &sf)

	if err != nil {
		return nil
	}

	return &sf
}
