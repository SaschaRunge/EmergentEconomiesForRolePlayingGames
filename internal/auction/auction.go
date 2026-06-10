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
	rng     *rpgMath.RNG
	askBook []ask
	bidBook []bid

	receipts map[int][]receipt
}

func (h *House) ResolveOffers(c commodity) {
	rpgMath.Shuffle(h.rng, h.askBook)
	rpgMath.Shuffle(h.rng, h.bidBook)

	slices.SortStableFunc(h.askBook, func(a, b ask) int {
		return cmp.Compare(a.Price, b.Price)
	})
	slices.SortStableFunc(h.bidBook, func(a, b bid) int {
		return cmp.Compare(b.Price, a.Price)
	})

	for len(h.askBook) > 0 && len(h.bidBook) > 0 {
		buyer := &h.bidBook[0]
		seller := &h.askBook[0]

		//might need revision, contrary to the paper this will cause cancelation of trades if there's no buyer willing to match a sellers price
		if buyer.Price < seller.Price {
			break
		}

		quantityTraded := min(buyer.Quantity, seller.Quantity)
		clearingPrice := (buyer.Price + seller.Price) / 2

		if quantityTraded > 0 {
			buyer.Quantity -= quantityTraded
			seller.Quantity -= quantityTraded

			if _, exists := h.receipts[buyer.AgentID]; !exists {
				h.receipts[buyer.AgentID] = []receipt{}
			}
			if _, exists := h.receipts[seller.AgentID]; !exists {
				h.receipts[seller.AgentID] = []receipt{}
			}

			h.receipts[buyer.AgentID] = append(h.receipts[buyer.AgentID], trade.NewReceipt(buyer.AgentID, c, clearingPrice, quantityTraded))
			h.receipts[seller.AgentID] = append(h.receipts[seller.AgentID], trade.NewReceipt(seller.AgentID, c, clearingPrice, quantityTraded))
		}

		if buyer.Quantity == 0 {
			h.bidBook = h.bidBook[1:]
		}
		if seller.Quantity == 0 {
			h.askBook = h.askBook[1:]
		}
	}
}

func (h *House) PlaceAsk(a ask) {
	h.askBook = append(h.askBook, a)
}

func (h *House) PlaceBid(b bid) {
	h.bidBook = append(h.bidBook, b)
}
