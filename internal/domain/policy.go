package domain

type Policy struct {
	CategoryCaps     map[string]int64
	ReceiptThreshold int64
}

func DefaultPolicy() Policy {
	return Policy{
		CategoryCaps: map[string]int64{
			"meals":    7_500,
			"supplies": 20_000,
			"travel":   80_000,
			"training": 50_000,
			"misc":     10_000,
		},
		ReceiptThreshold: 5_000,
	}
}

func (p Policy) capFor(category string) int64 {
	if cap, ok := p.CategoryCaps[category]; ok {
		return cap
	}
	return p.CategoryCaps["misc"]
}
