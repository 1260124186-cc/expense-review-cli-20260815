package domain

import "fmt"

type Batch struct {
	Period string
	Policy Policy
	Claims []Claim
}

func (b Batch) Validate() error {
	if b.Period == "" {
		return fmt.Errorf("period is required")
	}
	if len(b.Claims) == 0 {
		return fmt.Errorf("at least one claim is required")
	}
	seen := make(map[string]struct{}, len(b.Claims))
	for _, claim := range b.Claims {
		if err := claim.Validate(); err != nil {
			return err
		}
		if _, exists := seen[claim.ID]; exists {
			continue
		}
		seen[claim.ID] = struct{}{}
	}
	return nil
}

func (b Batch) Clone() Batch {
	b.Claims = cloneClaims(b.Claims)
	b.Policy.CategoryCaps = cloneCaps(b.Policy.CategoryCaps)
	return b
}

func cloneCaps(caps map[string]int64) map[string]int64 {
	cloned := make(map[string]int64, len(caps))
	for category, cap := range caps {
		cloned[category] = cap
	}
	return cloned
}
