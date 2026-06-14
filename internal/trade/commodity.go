package trade

type Commodity int

//go:generate stringer -type=Commodity
const (
	CommodityNone Commodity = iota
	CommodityFood
	CommodityWood
	CommodityIronOre
	CommodityIron
	CommodityTools
)

func (c Commodity) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}
