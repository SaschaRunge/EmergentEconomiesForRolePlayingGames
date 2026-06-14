package production

type Recipe struct {
	Name         string
	Input        map[commodity]int
	Output       map[commodity]int
	OutputChance map[commodity]float64
}
