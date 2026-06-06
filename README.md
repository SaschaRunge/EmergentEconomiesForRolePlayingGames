### Emergent Economies for Role Playing Games

Go-Implementation of Emergent Economies for Role Playing Games:

https://ianparberry.com/pubs/econ.pdf

## Definitions

mean: historical mean price of Commodity
favorability: position of mean within observered trading range (by individual agent), favorability = (mean - ObservedMin)/(ObservedMax - ObservedMin) // <- price belief rather than observed trading range

## Deviations from the Source

The original paper uses the observed trading range as a measure of favorability when deciding on the quantity of goods to buy/sell while speaking of using the price belief in the paragraph prior. I suspect the written text is correct, as using this measure seems more reasonable and also avoids needing an additional measure. As such, I use price belief to determine those quantities.

