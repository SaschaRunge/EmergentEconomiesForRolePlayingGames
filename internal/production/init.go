package production

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type commodity = trade.Commodity

var RecipeRegistry map[string]Recipe
var RoleRegistry map[string]Role

// TODO: decouple static seed
func init() {
	RecipeRegistry = loadRecipes()
	RoleRegistry = loadRoles()
}
