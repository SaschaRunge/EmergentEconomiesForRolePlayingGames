package market

type commodity interface {
}

type ask struct {
	commodity commodity
	price     float64
	quantity  int
}

type bid struct {
	item     commodity
	price    float64
	quantity int
}

func CreateAsk(c commodity, limit int) ask {
	return ask{}
}

func CreateBid(c commodity, limit int) bid {
	return bid{}
}
