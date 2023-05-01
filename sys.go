package main

import (
	"sync"

	"github.com/chaosknight/bitnet/bsocket"
	"github.com/chaosknight/bitnet/entity"
	"github.com/chaosknight/bitnet/rest"
	"github.com/chaosknight/skynet/skynet"
	"github.com/chaosknight/skynet/types"
)

type Sys struct {
	net     *skynet.SkyNet
	privews *bsocket.Bsocket
	pubws   *bsocket.Bsocket
	rest    *rest.ClientRest

	apikey   *entity.APIKey
	posready bool
	//仓位信息
	positons []*entity.Position

	confs []*entity.FinesseConfig
	// usdt
	pmux     sync.Mutex
	BreedMap sync.Map
}

func NewSys(apikey *entity.APIKey, cons []*entity.FinesseConfig) *Sys {
	return &Sys{
		apikey:   apikey,
		confs:    cons,
		rest:     rest.NewClient(apikey, entity.AwsRestURL, entity.DemoServer),
		net:      &skynet.SkyNet{},
		posready: false,
		privews:  bsocket.NewBsocket(string(entity.AwsPrivateWsURL)),
		pubws:    bsocket.NewBsocket(string(entity.AwsPublicWsURL)),
	}
}

func (sys *Sys) Init() {
	sys.net.Init(types.SkyNetInitOptions{})
	sys.actorinit()
	sys.pubinit()
	sys.pwsinit()
}
