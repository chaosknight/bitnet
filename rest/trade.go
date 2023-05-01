package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chaosknight/bitnet/entity"
	requests "github.com/chaosknight/bitnet/rest/request"
	responses "github.com/chaosknight/bitnet/rest/responses"
)

// Trade
//
// https://www.okex.com/docs-v5/en/#rest-api-trade
type Trade struct {
	client *ClientRest
}

// NewTrade returns a pointer to a fresh Trade
func NewTrade(c *ClientRest) *Trade {
	return &Trade{c}
}

// PlaceOrder
// You can place an order only if you have sufficient funds.
//
// https://www.okex.com/docs-v5/en/#rest-api-trade-get-positions
func (c *Trade) PlaceOrder(req []requests.PlaceOrder) (response responses.PlaceOrder, err error) {
	p := "/api/v5/trade/order"
	var m interface{}
	m = entity.S2M(req[0])
	if len(req) > 1 {
		m = entity.S2SM(req)
		p = "/api/v5/trade/batch-orders"
	}
	log.Println("参数:", m)

	res, err := c.client.Do(http.MethodPost, p, true, m)
	if err != nil {
		return
	}
	defer res.Body.Close()
	d := json.NewDecoder(res.Body)
	err = Decode(d, &response)

	return
}

// ClosePosition
// Close all positions of an instrument via a market order.
//
// https://www.okex.com/docs-v5/en/#rest-api-trade-close-positions
func (c *Trade) ClosePosition(req requests.ClosePosition) (response responses.ClosePosition, err error) {
	p := "/api/v5/trade/close-position"
	m := entity.S2M(req)
	res, err := c.client.Do(http.MethodPost, p, true, m)
	if err != nil {
		return
	}
	defer res.Body.Close()
	d := json.NewDecoder(res.Body)
	err = Decode(d, &response)
	return
}

// GetOrderDetail
// Retrieve order details.
//
// https://www.okex.com/docs-v5/en/#rest-api-trade-get-order-details
func (c *Trade) GetOrderDetail(req requests.OrderDetails) (response responses.OrderList, err error) {
	p := "/api/v5/trade/order"
	m := entity.S2M(req)
	res, err := c.client.Do(http.MethodGet, p, true, m)
	if err != nil {
		return
	}
	defer res.Body.Close()
	d := json.NewDecoder(res.Body)
	err = Decode(d, &response)
	return
}
