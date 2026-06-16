package production

type Role struct {
	Name    string
	Recipes []Recipe
}

func loadRoles() map[string]Role {
	return map[string]Role{
		"Blacksmith": {
			Name: "Blacksmith",

			Recipes: []Recipe{
				RecipeRegistry["Iron"],
			},
		},
		"Farmer": {
			Name: "Farmer",

			Recipes: []Recipe{
				RecipeRegistry["Food"],
			},
		},
	}
}
