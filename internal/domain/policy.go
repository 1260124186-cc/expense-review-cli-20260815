package domain

type Policy struct {
	CategoryCaps     map[string]int64
	ReceiptThreshold int64
}

func DefaultPolicy() Policy {
	return Policy{
		CategoryCaps:     nil,
		ReceiptThreshold: 5_000,
	}
}

func (p Policy) capFor(category string) int64 {
	if cap, ok := p.CategoryCaps[category]; ok {
		return cap
	}
	return p.CategoryCaps["misc"]
}
