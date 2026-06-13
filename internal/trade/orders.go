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
}

func (r *Receipt) Add(receipt Receipt) Receipt {
	if r.AgentID != receipt.AgentID {
		panic("receipt contains mismatching agents")
	}

	if r.Commodity != receipt.Commodity {
		panic("receipt contains mismatching commodities")
	}

	return Receipt{
		Order: Order{
			AgentID:   r.AgentID,
			Commodity: r.Commodity,
			Price:     r.Price + receipt.Price,
			Quantity:  r.Quantity + receipt.Quantity,
		},
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
