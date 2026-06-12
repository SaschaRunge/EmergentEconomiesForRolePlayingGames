package auction

import (
	"cmp"
	"slices"

	_ "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/agent"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type ask = trade.Ask
type bid = trade.Bid
type receipt = trade.Receipt
type commodity = trade.Commodity

type House struct {
	rng *rpgMath.RNG

	daysToArchive int
	statistics    []Statistics
}

type Statistics struct {
	Commodity            commodity
	Supply               int
	Demand               int
	UnitsSold            int
	AverageClearingPrice float64
}

func New(rng *rpgMath.RNG) House {
	return House{
		rng: rng,
	}
}

func (h *House) ResolveOffers(c commodity, asks []ask, bids []bid) map[int][]receipt {
	receipts := make(map[int][]receipt)

	rpgMath.Shuffle(h.rng, asks)
	rpgMath.Shuffle(h.rng, bids)

	slices.SortStableFunc(asks, func(a, b ask) int {
		return cmp.Compare(a.Price, b.Price)
	})
	slices.SortStableFunc(bids, func(a, b bid) int {
		return cmp.Compare(b.Price, a.Price)
	})

	for len(asks) > 0 && len(bids) > 0 {
		buyer := &bids[0]
		seller := &asks[0]

		//might need revision, contrary to the paper this will cause cancelation of trades if there's no buyer willing to match a sellers price
		if buyer.Price < seller.Price {
			break
		}

		quantityTraded := min(buyer.Quantity, seller.Quantity)
		clearingPrice := (buyer.Price + seller.Price) / 2

		if quantityTraded > 0 {
			buyer.Quantity -= quantityTraded
			seller.Quantity -= quantityTraded

			receipts[buyer.AgentID] = append(receipts[buyer.AgentID], trade.NewReceipt(buyer.AgentID, c, clearingPrice, quantityTraded))
			receipts[seller.AgentID] = append(receipts[seller.AgentID], trade.NewReceipt(seller.AgentID, c, clearingPrice, quantityTraded))
		}

		if buyer.Quantity == 0 {
			bids = bids[1:]
		}
		if seller.Quantity == 0 {
			asks = asks[1:]
		}
	}

	return receipts
}
