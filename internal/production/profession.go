package production

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type commodity = trade.Commodity

type Role struct {
	rng *rpgMath.RNG

	Name    string
	Recipes []Recipe
}
