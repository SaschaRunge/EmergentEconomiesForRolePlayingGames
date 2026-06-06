package economy

type ask struct {
	commodity Commodity
	price     float64
	quantity  int
}

type bid struct {
	commodity Commodity
	price     float64
	quantity  int
}
