package entity

type CandleData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data [][9]string `json:"data"`
}

type ChannelinstId struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}
type SocketCandle struct {
	Arg  ChannelinstId `json:"arg"`
	Data []*Candle     `json:"data"`
}
