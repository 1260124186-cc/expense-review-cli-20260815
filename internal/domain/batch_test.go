package domain

import (
	"strings"
	"testing"
)

func TestBatchRejectsDuplicateClaimIDs(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "meal-1", Employee: "Ari", Category: "meals", AmountCents: 4200},
			{ID: "meal-1", Employee: "Bo", Category: "meals", AmountCents: 4300},
		},
	}
	err := batch.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want duplicate claim rejection")
	}
	if got, want := err.Error(), `duplicate claim ID: "meal-1"`; !strings.Contains(got, want) {
		t.Fatalf("Validate() error = %q, want message containing %q", got, want)
	}
}
