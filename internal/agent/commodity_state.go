package agent

import (
	"encoding/json"

	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
)

type inventory struct {
	capacity      int
	idealQuantity int
	quantity      int
}

func (i *inventory) AvailableSpace() int {
	return i.capacity - i.quantity
}

func (i *inventory) Excess() int {
	if i.quantity > i.idealQuantity {
		return i.quantity - i.idealQuantity
	}
	return 0
}

type CommodityState struct {
	inventory

	historicalMean float64
	priceBelief    rpgMath.PriceRange
}

func (c *CommodityState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Capacity       int     `json:"capacity"`
		IdealQuantity  int     `json:"ideal_quantity"`
		Quantity       int     `json:"quantity"`
		HistoricalMean float64 `json:"historical_mean"`
		PriceBelief    string  `json:"price_belief"`
	}{
		Capacity:      c.capacity,
		IdealQuantity: c.idealQuantity,
		Quantity:      c.quantity,
	})
}
