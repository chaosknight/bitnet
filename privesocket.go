package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chaosknight/bitnet/entity"
	"github.com/chaosknight/gowebsocket"
)

func (sys *Sys) pwsinit() {
	sys.privews.Init(sys.privemessg, sys.priveconnected)
	sys.privews.Connect()
}

func (sys *Sys) priveconnected(socked gowebsocket.Socket) {
	signstr := `{"op":"login","args":[` + getloginstr(sys.apikey) + `]}`
	log.Println("send login msg:", signstr)
	socked.SendText(signstr)
}

func (sys *Sys) privemessg(message string, socked gowebsocket.Socket) {
	rep, err := entity.UnmarshalResponse([]byte(message))

	if err != nil {
		log.Println("数据处理错误:", message)
		return
	}
	if rep.Code > 0 {
		log.Println("出现错误:", message)
		return
	}

	if entity.IsOpRes(rep) {
		if rep.Event == "login" {
			log.Println("login sucess")
			bpsub := `{"op":"subscribe","args":[{"channel": "positions","instType":"SWAP"}]}`
			socked.SendText(bpsub)
		}
	} else {
		chel, ok := rep.Arg["channel"]
		if ok {
			switch chel {
			case "positions":
				bps := bpmsg(rep.Data)
				sys.pmux.Lock()
				sys.posready = true
				sys.positons = bps
				sys.pmux.Unlock()
				//执行策略判断是否加仓平仓
				sys.net.SendMsg(EmptyActor, "")

				// log.Println(sys.positons)
			default:
				log.Println(message)
			}
		}

	}
}

func bpmsg(b []byte) []*entity.Position {
	var rec []*entity.Position
	err := json.Unmarshal(b, &rec)
	if err != nil {
		log.Fatalln("decode failed %s", err)
		return rec
	}
	return rec
}

func getloginstr(apikey *entity.APIKey) string {
	ts, signature := sign(apikey)
	return `{"apiKey":"` + apikey.ApiKey + `","passphrase":"` + apikey.Passphrase + `","timestamp":"` + ts + `","sign":"` + signature + `"}`
}

func sign(apikey *entity.APIKey) (string, string) {
	method := "GET"
	path := "/users/self/verify"

	t := time.Now().Unix()
	ts := fmt.Sprint(t)
	s := ts + method + path
	p := []byte(s)
	h := hmac.New(sha256.New, apikey.SecretKey)
	h.Write(p)
	return ts, base64.StdEncoding.EncodeToString(h.Sum(nil))
}
