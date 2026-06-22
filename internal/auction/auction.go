package auction

import (
	"cmp"
	"slices"

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
	archive       map[commodity][]Report
}

type Report struct {
	Supply               int
	Demand               int
	UnitsSold            int
	AverageClearingPrice float64
}

func New(daysToArchive int, rng *rpgMath.RNG) *House {
	return &House{
		rng: rng,

		daysToArchive: daysToArchive,
		archive:       make(map[commodity][]Report),
	}
}

func (h *House) ResolveOffers(c commodity, asks []ask, bids []bid) map[int][]receipt {
	for _, a := range asks {
		if a.Commodity != c {
			panic("auction house: ask contained wrong commodity")
		}
	}
	for _, b := range bids {
		if b.Commodity != c {
			panic("auction house: bid contained wrong commodity")
		}
	}

	receiptsByAgentID := make(map[int][]receipt)

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

		if quantityTraded > 0 {
			clearingPrice := (buyer.Price + seller.Price) / 2

			buyer.Quantity -= quantityTraded
			seller.Quantity -= quantityTraded

			// TODO: test invariant all receipts need to have correct agent id
			receiptsByAgentID[buyer.AgentID] = append(receiptsByAgentID[buyer.AgentID], trade.NewReceipt(buyer.AgentID, c, clearingPrice, quantityTraded))
			receiptsByAgentID[seller.AgentID] = append(receiptsByAgentID[seller.AgentID], trade.NewReceipt(seller.AgentID, c, clearingPrice, -quantityTraded))
		}

		if buyer.Quantity == 0 {
			bids = bids[1:]
		}
		if seller.Quantity == 0 {
			asks = asks[1:]
		}
	}

	return receiptsByAgentID
}

func (h *House) archiveReport(c commodity, dailyReport Report) {
	h.archive[c] = append(h.archive[c], dailyReport)

	if len(h.archive[c]) > h.daysToArchive {
		overflow := len(h.archive[c]) - h.daysToArchive
		h.archive[c] = h.archive[c][overflow:]
	}
}
