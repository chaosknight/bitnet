package main

import (
	"encoding/json"
	"log"

	// "log"

	"strings"

	"github.com/chaosknight/bitnet/cell"
	"github.com/chaosknight/bitnet/entity"
	"github.com/chaosknight/gowebsocket"
)

func (sys *Sys) pubinit() {
	sys.pubws.Init(sys.pubmessg, sys.pubconnected)
	sys.pubws.Connect()
}

func (sys *Sys) pubconnected(socked gowebsocket.Socket) {
	socked.SendText(`{"op":"subscribe","args":[` + getsubstr(sys.confs) + `]}`)
}

func (sys *Sys) pubmessg(message string, socket gowebsocket.Socket) {
	// log.Println(message)
	if strings.Contains(message, `"channel":"candle5m"`) {
		var msg entity.SocketCandle
		json.Unmarshal([]byte(message), &msg)

		if len(msg.Data) != 1 {
			return
		}
		instId := msg.Arg.InstId

		v, ok := sys.BreedMap.Load(instId)
		if !ok {
			v = cell.NewBreed(instId)
			sys.BreedMap.Store(instId, v)
			sys.net.SendMsg(LoadBreed, instId)
		} else {
			brd := v.(*cell.Breed)
			if brd.SetNewK(msg.Data[0]) {
				if sys.posready {
					tdc, ma7, lastc := brd.GetTdc()
					if lastc == 0 {
						return
					}
					bios := (lastc - ma7) * 100 / lastc
					log.Println(instId, "bios:", bios)

					//建仓
					if brd.IsTiming() {

						if tdc > 0 && bios >= float64(0.75) {
							sys.net.SendMsg(EmptyActor, brd.InstId, tdc)
						}

						if tdc < 0 && bios <= float64(-0.75) {
							sys.net.SendMsg(EmptyActor, brd.InstId, tdc)
						}

						if tdc >= 14 && bios >= float64(0.3) {
							sys.net.SendMsg(EmptyActor, brd.InstId, tdc)
						}

						if tdc <= -14 && bios <= float64(-0.3) {
							sys.net.SendMsg(EmptyActor, brd.InstId, tdc)
						}

					}

				}
			}
		}
	}
}

func getsubstr(insts []*entity.FinesseConfig) string {
	str := ""
	for i, v := range insts {
		if i > 0 {
			str = str + ","
		}
		str = str + `{"channel":"candle5m","instId":"` + v.InstID + `"}`
	}
	return str
}
