package trade

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

	PriceMean      float64
	Demand         int
	Supply         int
	TotalUnitsSold int
}

func NewAsk(agentID int, commodity Commodity, price float64, quantity int) Ask {
	return Ask{
		Order{
			AgentID:   agentID,
			Commodity: commodity,
			Price:     price,
			Quantity:  quantity,
		},
	}
}

func NewBid(agentID int, commodity Commodity, price float64, quantity int) Bid {
	return Bid{
		Order{
			AgentID:   agentID,
			Commodity: commodity,
			Price:     price,
			Quantity:  quantity,
		},
	}
}
