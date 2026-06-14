package trade

//go:generate stringer -type=Commodity
type Commodity int

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
