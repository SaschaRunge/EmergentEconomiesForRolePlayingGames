package market

type Ask struct {
	Commodity Commodity
	Price     float64
	Quantity  int
}

type Bid struct {
	Commodity Commodity
	Price     float64
	Quantity  int
}

type Receipt struct {
	Commodity Commodity
	Price     float64
	Quantity  int
}
