package main

import (
	"log"
	"math"
	"time"

	"github.com/chaosknight/bitnet/cell"
	"github.com/chaosknight/bitnet/entity"

	requests "github.com/chaosknight/bitnet/rest/request"

	"github.com/chaosknight/skynet/actor"
	"github.com/chaosknight/skynet/types"
)

const LoadBreed = "LoadBreed"

const OrderActor = "OrderActor"

const EmptyActor = "EmptyActor"

func (sys *Sys) actorinit() {
	//加载历史k线数据
	breedactor := actor.NewFromReducer(LoadBreed, 20, func(a types.Actor, msg *types.MasterMsg) {
		log.Println("加载k线数据 :", msg.Cmd)
		brd, ok := sys.BreedMap.Load(msg.Cmd)
		if ok {
			(brd.(*cell.Breed)).LoadeK(sys.rest)
			time.Sleep(time.Duration(10) * time.Second)
		}
	})
	sys.net.Rigist(breedactor, 1)

	//平仓
	emptyactor := actor.NewFromReducer(EmptyActor, 1024, func(a types.Actor, msg *types.MasterMsg) {
		log.Println("策略执行中 :", msg.Cmd)
		sys.pmux.Lock()
		var hasdo = false
		var hashttp = false

		//
		for _, v := range sys.positons {
			//排除非正常数据
			if v.Pos == 0 {
				break
			}
			if v.InstID == msg.Cmd {
				hasdo = true
			}
			cnf := getConfigbyid(v.InstID, sys.confs)

			brdi, ok := sys.BreedMap.Load(v.InstID)
			//获取k线数据失败
			if !ok {
				log.Println("k线数据失败")
				break
			}
			brd := (brdi.(*cell.Breed))

			//正常平仓
			_, ma7, lastc := brd.GetTdc()

			isclose := false
			log.Println("lastc:", lastc, "ma7:", ma7, "bool:", lastc < ma7)
			//平仓条件
			if float64(v.Pos) < 0 && lastc < ma7 {

				isclose = true
			}
			if float64(v.Pos) > 0 && lastc > ma7 {
				isclose = true
			}
			//平仓
			if isclose || float64(v.Upl)+cnf.Maxloss < 0 {
				hashttp = true
				sys.rest.Trade.ClosePosition(requests.ClosePosition{InstID: v.InstID, MgnMode: v.MgnMode})
				break
			}

			//加仓
			if float64(v.UplRatio)*100 < cnf.Addposition && math.Abs(float64(v.Pos))*2 < cnf.Maxcount {
				hashttp = true
				poside := entity.OrderBuy

				if float64(v.Pos) < 0 {
					poside = entity.OrderSell
				}

				sys.rest.Trade.PlaceOrder([]requests.PlaceOrder{
					requests.PlaceOrder{
						InstID:  v.InstID,
						TdMode:  entity.TradeCrossMode,
						Side:    poside,
						OrdType: entity.OrderMarket,
						Sz:      math.Abs(float64(v.Pos)),
					},
				})

				break
			}
		}

		//底仓
		if !hasdo && msg.Cmd != "" {
			log.Println("建仓开始")

			tdc, ok := msg.Args[0].(int)
			if len(msg.Args) > 0 && ok {
				hashttp = true
				poside := entity.OrderBuy

				if tdc > 0 {
					poside = entity.OrderSell
				}
				cnf := getConfigbyid(msg.Cmd, sys.confs)
				_, err := sys.rest.Trade.PlaceOrder([]requests.PlaceOrder{
					requests.PlaceOrder{
						InstID:  msg.Cmd,
						TdMode:  entity.TradeCrossMode,
						OrdType: entity.OrderMarket,
						Side:    poside,
						Sz:      cnf.Firstcount,
					},
				})
				if err != nil {
					log.Println(err.Error())
				} else {
					log.Println("建仓完成")
				}
			}
		}

		sys.pmux.Unlock()
		if hashttp {
			time.Sleep(time.Duration(5) * time.Second)
		}
	})

	sys.net.Rigist(emptyactor, 1)

}

func getConfigbyid(id string, ids []*entity.FinesseConfig) *entity.FinesseConfig {
	for _, v := range ids {
		if v.InstID == id {
			return v
		}
	}
	return nil
}
