package domain

import "fmt"

type Claim struct {
	ID          string
	Employee    string
	Category    string
	AmountCents int64
	ReceiptIDs  []string
}

func (c Claim) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("claim id is required")
	}
	if c.Employee == "" {
		return fmt.Errorf("employee is required for claim %s", c.ID)
	}
	if c.Category == "" {
		return fmt.Errorf("category is required for claim %s", c.ID)
	}
	if c.AmountCents <= 0 {
		return fmt.Errorf("amount must be positive for claim %s", c.ID)
	}
	return nil
}

func cloneClaims(claims []Claim) []Claim {
	cloned := make([]Claim, len(claims))
	for i, claim := range claims {
		cloned[i] = claim
		cloned[i].ReceiptIDs = append([]string(nil), claim.ReceiptIDs...)
	}
	return cloned
}
