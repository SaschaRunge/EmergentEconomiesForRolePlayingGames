package auction

import (
	"encoding/json"
	"testing"

	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type ReceiptsByAgentID map[int][]receipt

func (r ReceiptsByAgentID) String() string {
	asString, _ := json.MarshalIndent(r, "", "    ")
	return string(asString)
}

func TestResolveOffers(t *testing.T) {
	cases := []struct {
		name     string
		asks     []ask
		bids     []bid
		expected ReceiptsByAgentID
	}{
		{
			name:     "empty input",
			asks:     []ask{},
			bids:     []bid{},
			expected: ReceiptsByAgentID{},
		},
		{
			name: "simple",
			asks: []ask{trade.NewAsk(1, trade.CommodityWood, 50, 10)},
			bids: []bid{trade.NewBid(2, trade.CommodityWood, 100, 10)},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 75, -10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 75, 10),
				},
			},
		},
		{
			name: "buyer half-matched",
			asks: []ask{trade.NewAsk(1, trade.CommodityWood, 50, 10)},
			bids: []bid{trade.NewBid(2, trade.CommodityWood, 100, 20)},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 75, -10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 75, 10),
				},
			},
		},
		{
			name: "seller half-matched",
			asks: []ask{trade.NewAsk(1, trade.CommodityWood, 50, 20)},
			bids: []bid{trade.NewBid(2, trade.CommodityWood, 100, 10)},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 75, -10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 75, 10),
				},
			},
		},
		{
			name: "multiple agents, perfect matches",
			asks: []ask{
				trade.NewAsk(1, trade.CommodityWood, 50, 10),
				trade.NewAsk(2, trade.CommodityWood, 80, 20),
				trade.NewAsk(3, trade.CommodityWood, 120, 30),
			},
			bids: []bid{
				trade.NewBid(4, trade.CommodityWood, 150, 10),
				trade.NewBid(5, trade.CommodityWood, 140, 20),
				trade.NewBid(6, trade.CommodityWood, 130, 30),
			},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 100, -10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 110, -20),
				},
				3: []receipt{
					trade.NewReceipt(3, trade.CommodityWood, 125, -30),
				},
				4: []receipt{
					trade.NewReceipt(4, trade.CommodityWood, 100, 10),
				},
				5: []receipt{
					trade.NewReceipt(5, trade.CommodityWood, 110, 20),
				},
				6: []receipt{
					trade.NewReceipt(6, trade.CommodityWood, 125, 30),
				},
			},
		},
		{
			name: "single seller, multiple buyers, partial fills",
			asks: []ask{
				trade.NewAsk(1, trade.CommodityWood, 50, 30),
			},
			bids: []bid{
				trade.NewBid(2, trade.CommodityWood, 150, 10),
				trade.NewBid(3, trade.CommodityWood, 130, 20),
			},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 100, -10),

					trade.NewReceipt(1, trade.CommodityWood, 90, -20),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 100, 10),
				},
				3: []receipt{
					trade.NewReceipt(3, trade.CommodityWood, 90, 20),
				},
			},
		},
		{
			name: "8 agents, partial fills, unfulfilled leftovers",
			asks: []ask{
				trade.NewAsk(1, trade.CommodityWood, 50, 10),
				trade.NewAsk(2, trade.CommodityWood, 60, 15),
				trade.NewAsk(3, trade.CommodityWood, 80, 20),
				trade.NewAsk(4, trade.CommodityWood, 130, 10),
			},
			bids: []bid{
				trade.NewBid(5, trade.CommodityWood, 150, 15),
				trade.NewBid(6, trade.CommodityWood, 110, 15),
				trade.NewBid(7, trade.CommodityWood, 90, 10),
				trade.NewBid(8, trade.CommodityWood, 40, 20),
			},
			expected: ReceiptsByAgentID{
				1: []receipt{
					trade.NewReceipt(1, trade.CommodityWood, 100, -10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 105, -5),
					trade.NewReceipt(2, trade.CommodityWood, 85, -10),
				},
				3: []receipt{
					trade.NewReceipt(3, trade.CommodityWood, 95, -5),
					trade.NewReceipt(3, trade.CommodityWood, 85, -10),
				},
				5: []receipt{
					trade.NewReceipt(5, trade.CommodityWood, 100, 10),
					trade.NewReceipt(5, trade.CommodityWood, 105, 5),
				},
				6: []receipt{
					trade.NewReceipt(6, trade.CommodityWood, 85, 10),
					trade.NewReceipt(6, trade.CommodityWood, 95, 5),
				},
				7: []receipt{
					trade.NewReceipt(7, trade.CommodityWood, 85, 10),
				},
			},
		},
	}

	for _, c := range cases {
		auctionHouse := New(10, rpgMath.NewRNG(42))

		var actual ReceiptsByAgentID = auctionHouse.ResolveOffers(trade.CommodityWood, c.asks, c.bids)

		if !equal(actual, c.expected, rpgMath.Epsilon) {
			t.Errorf("test %q failed:\nexpected: \n%v \nactual: \n%v", c.name, c.expected, actual)
		}
	}
}

func equal(r1, r2 ReceiptsByAgentID, epsilon float64) bool {
	if len(r1) != len(r2) {
		return false
	}

	for id := range r1 {
		if _, exists := r2[id]; !exists {
			return false
		}

		if len(r1[id]) != len(r2[id]) {
			return false
		}

		for i := range r1[id] {
			differentID := r1[id][i].AgentID != r2[id][i].AgentID
			differentCommodity := r1[id][i].Commodity != r2[id][i].Commodity
			differentPrice := !rpgMath.AlmostEquals(r1[id][i].Price, r2[id][i].Price, epsilon)
			differentQuantity := r1[id][i].Quantity != r2[id][i].Quantity

			if differentID || differentCommodity || differentPrice || differentQuantity {
				return false
			}
		}
	}

	return true
}
