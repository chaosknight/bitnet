package responses

import (
	"github.com/chaosknight/bitnet/entity"
)

type (
	PlaceOrder struct {
		Basic
		PlaceOrders []*entity.PlaceOrder `json:"data"`
	}
	CancelOrder struct {
		Basic
		CancelOrders []*entity.CancelOrder `json:"data"`
	}
	AmendOrder struct {
		Basic
		AmendOrders []*entity.AmendOrder `json:"data"`
	}
	ClosePosition struct {
		Basic
		ClosePositions []*entity.ClosePosition `json:"data"`
	}
	OrderList struct {
		Basic
		Orders []*entity.Order `json:"data"`
	}
	TransactionDetail struct {
		Basic
		TransactionDetails []*entity.TransactionDetail `json:"data"`
	}
	PlaceAlgoOrder struct {
		Basic
		PlaceAlgoOrders []*entity.PlaceAlgoOrder `json:"data"`
	}
	CancelAlgoOrder struct {
		Basic
		CancelAlgoOrders []*entity.CancelAlgoOrder `json:"data"`
	}
	AlgoOrderList struct {
		Basic
		AlgoOrders []*entity.AlgoOrder `json:"data"`
	}
)
