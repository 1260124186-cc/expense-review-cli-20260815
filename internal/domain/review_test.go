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

func TestReviewBatchOrdersDecisionsByID(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "z-last", Employee: "Ari", Category: "meals", AmountCents: 4200},
			{ID: "a-first", Employee: "Bo", Category: "meals", AmountCents: 4300},
		},
	}
	review, err := ReviewBatch(batch)
	if err != nil {
		t.Fatalf("ReviewBatch() error = %v", err)
	}
	if got, want := review.Decisions[0].ClaimID, "a-first"; got != want {
		t.Fatalf("first decision = %q, want %q", got, want)
	}
	if got, want := review.Decisions[1].ClaimID, "z-last"; got != want {
		t.Fatalf("second decision = %q, want %q", got, want)
	}
	if got, want := review.Total, int64(8500); got != want {
		t.Fatalf("total = %d, want %d", got, want)
	}
}

func TestReviewBatchRepeatedReviewKeepsOriginalOrderStable(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "z-last", Employee: "Ari", Category: "meals", AmountCents: 4200},
			{ID: "a-first", Employee: "Bo", Category: "meals", AmountCents: 4300},
			{ID: "m-mid", Employee: "Cy", Category: "meals", AmountCents: 4400},
		},
	}
	want := []string{"z-last", "a-first", "m-mid"}

	for i := 0; i < 3; i++ {
		review, err := ReviewBatch(batch)
		if err != nil {
			t.Fatalf("ReviewBatch() pass %d error = %v", i, err)
		}
		if got := claimIDs(batch.Claims); !equal(got, want) {
			t.Fatalf("pass %d original order = %v, want %v", i, got, want)
		}
		if got := decisionIDs(review.Decisions); !equal(got, []string{"a-first", "m-mid", "z-last"}) {
			t.Fatalf("pass %d decision order = %v, want sorted", i, got)
		}
	}
}

func claimIDs(claims []Claim) []string {
	ids := make([]string, len(claims))
	for i, claim := range claims {
		ids[i] = claim.ID
	}
	return ids
}

func decisionIDs(decisions []Decision) []string {
	ids := make([]string, len(decisions))
	for i, decision := range decisions {
		ids[i] = decision.ClaimID
	}
	return ids
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
