package domain

import "testing"

func TestReviewBatchDoesNotReorderClaims(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "z-last", Employee: "Ari", Category: "meals", AmountCents: 4200},
			{ID: "a-first", Employee: "Bo", Category: "meals", AmountCents: 4300},
		},
	}
	if _, err := ReviewBatch(batch); err != nil {
		t.Fatalf("ReviewBatch() error = %v", err)
	}
	if got, want := batch.Claims[0].ID, "z-last"; got != want {
		t.Fatalf("claim order changed to %q, want %q", got, want)
	}
}
