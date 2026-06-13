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
	epsilon := 0.0001

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
					trade.NewReceipt(1, trade.CommodityWood, 75, 10),
				},
				2: []receipt{
					trade.NewReceipt(2, trade.CommodityWood, 75, 10),
				},
			},
		},
	}

	auctionHouse := New(10, rpgMath.NewRNG(42))

	for _, c := range cases {
		var actual ReceiptsByAgentID = auctionHouse.ResolveOffers(trade.CommodityWood, c.asks, c.bids)

		if !equal(actual, c.expected, epsilon) {
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
