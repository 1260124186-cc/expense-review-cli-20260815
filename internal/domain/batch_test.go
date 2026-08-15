package domain

import "testing"

func TestBatchRejectsDuplicateClaimIDs(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "meal-1", Employee: "Ari", Category: "meals", AmountCents: 4200},
			{ID: "meal-1", Employee: "Bo", Category: "meals", AmountCents: 4300},
		},
	}
	if err := batch.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want duplicate claim rejection")
	}
}
