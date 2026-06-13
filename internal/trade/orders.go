package trade

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
)

type Order struct {
	AgentID   int
	Commodity Commodity
	Price     float64
	Quantity  int
}

type Ask struct {
	Order
}

type Bid struct {
	Order
}

type Receipt struct {
	Order

	//TODO: remove
	PriceMean      float64
	Demand         int
	Supply         int
	TotalUnitsSold int
}

func (r *Receipt) Add(receipt Receipt) {
	if r.AgentID != receipt.AgentID {
		panic("receipt contains mismatching agents")
	}

	if r.Commodity != receipt.Commodity {
		panic("receipt contains mismatching commodities")
	}

	if r.Quantity+receipt.Quantity > 0 {
		r.Price = rpgMath.WeightedMean(r.Price, receipt.Price, float64(r.Quantity), float64(receipt.Quantity))
		r.Quantity += receipt.Quantity
	}
}

func (r *Receipt) MergeInto(receipts []Receipt) {
	for _, receipt := range receipts {
		r.Add(receipt)
	}
}

func NewAsk(agentID int, commodity Commodity, price float64, quantity int) Ask {
	return Ask{
		Order: Order{
			AgentID:   agentID,
			Commodity: commodity,
			Price:     price,
			Quantity:  quantity,
		},
	}
}

func NewBid(agentID int, commodity Commodity, price float64, quantity int) Bid {
	return Bid{
		Order: Order{
			AgentID:   agentID,
			Commodity: commodity,
			Price:     price,
			Quantity:  quantity,
		},
	}
}

func NewReceipt(agentID int, commodity Commodity, price float64, quantity int) Receipt {
	return Receipt{
		Order: Order{
			AgentID:   agentID,
			Commodity: commodity,
			Price:     price,
			Quantity:  quantity,
		},
	}
}
