package cell

import (
	"log"
	"math"
	"sort"
	"sync"
	"time"

	// "log"
	"github.com/chaosknight/bitnet/entity"
	"github.com/chaosknight/bitnet/rest"

	requests "github.com/chaosknight/bitnet/rest/request"

	// responses "github.com/chaosknight/bitnet/rest/responses"
	"github.com/chaosknight/bitnet/util"
)

type Breed struct {
	InstId   string
	kcandles []*entity.Candle
	tdcount  []int
	mutex    sync.RWMutex
}

func NewBreed(instId string) *Breed {
	return &Breed{
		InstId:   instId,
		kcandles: []*entity.Candle{},
		tdcount:  make([]int, 500),
	}
}

func (brd *Breed) LoadeK(client *rest.ClientRest) bool {
	//时间判断是否真正需要加载数据
	brd.mutex.Lock()
	defer brd.mutex.Unlock()
	brd.kcandles = brd.kcandles[0:0]
	list, ok := client.Market.GetCandlesticks(requests.GetCandlesticks{
		InstID: brd.InstId,
		Bar:    entity.Bar5m,
	})

	if ok == nil && list.Code == 0 {
		brd.kcandles = list.Candles
		sort.SliceStable(brd.kcandles, func(i, j int) bool {
			return ((time.Time)(brd.kcandles[i].TS)).Before((time.Time)(brd.kcandles[j].TS))
		})
		brd.tdcount = brd.tdcount[0:len(brd.kcandles)]
		util.SetTdcunt(brd.kcandles, brd.tdcount)
		// log.Println(brd.kcandles)
		log.Println("加载k线数据完成:", brd.InstId, len(brd.kcandles), "条")

	} else {
		log.Println(ok.Error())
	}

	return len(brd.kcandles) > 0
}

func (brd *Breed) SetNewK(k *entity.Candle) bool {
	brd.mutex.Lock()
	defer brd.mutex.Unlock()
	klen := len(brd.kcandles)
	if klen < 1 {
		return false
	}
	lastk := brd.kcandles[klen-1]

	if k.TS == lastk.TS {
		lastk.O = k.O
		lastk.H = k.H
		lastk.L = k.L
		lastk.C = k.C
		lastk.Confirm = k.Confirm
	} else {
		brd.kcandles = append(brd.kcandles, k)
		brd.tdcount = append(brd.tdcount, 0)
		if klen > 200 {
			//删除过去无用的数据
		}
	}
	util.SetTdcuntinde(brd.kcandles, brd.tdcount, len(brd.kcandles)-1)
	return true
}

func (brd *Breed) IsTiming() bool {
	brd.mutex.RLock()
	defer brd.mutex.RUnlock()
	size := len(brd.tdcount)
	abss := math.Abs(float64(brd.tdcount[size-1]))
	return abss >= 9
}

func (brd *Breed) GetTdc() (tdc int, ma7 float64, lastc float64) {
	brd.mutex.RLock()
	defer brd.mutex.RUnlock()
	tdc = 0
	ma7 = 0
	size := len(brd.tdcount)
	if size > 7 {
		tdc = brd.tdcount[size-1]
		for i := 1; i < 8; i++ {
			ma7 += brd.kcandles[size-i].C
		}
		ma7 = ma7 / 7
		lastc = brd.kcandles[size-1].C
	}
	return
}
