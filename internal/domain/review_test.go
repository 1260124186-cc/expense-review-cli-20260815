package domain

import (
	"reflect"
	"testing"
)

func TestReviewBatchPreservesBatchAcrossRepeatedReviews(t *testing.T) {
	batch := Batch{
		Period: "2026-08",
		Policy: DefaultPolicy(),
		Claims: []Claim{
			{ID: "z-approved", Employee: "Ari", Category: "meals", AmountCents: 4200, ReceiptIDs: []string{"r-3"}},
			{ID: "a-review", Employee: "Bo", Category: "meals", AmountCents: 6000},
			{ID: "m-rejected", Employee: "Cy", Category: "travel", AmountCents: 81000, ReceiptIDs: []string{"r-1"}},
		},
	}
	original := batch.Clone()

	first, err := ReviewBatch(batch)
	if err != nil {
		t.Fatalf("first ReviewBatch() error = %v", err)
	}
	want := Review{
		Period: "2026-08",
		Total:  91200,
		Decisions: []Decision{
			{ClaimID: "a-review", Status: StatusReview, Reason: "receipt required"},
			{ClaimID: "m-rejected", Status: StatusRejected, Reason: "category cap exceeded"},
			{ClaimID: "z-approved", Status: StatusApproved},
		},
	}
	if !reflect.DeepEqual(first, want) {
		t.Fatalf("first ReviewBatch() = %#v, want %#v", first, want)
	}
	if !reflect.DeepEqual(batch, original) {
		t.Fatalf("first ReviewBatch() changed batch = %#v, want %#v", batch, original)
	}

	second, err := ReviewBatch(batch)
	if err != nil {
		t.Fatalf("second ReviewBatch() error = %v", err)
	}
	if !reflect.DeepEqual(second, first) {
		t.Fatalf("second ReviewBatch() = %#v, want %#v", second, first)
	}
	if !reflect.DeepEqual(batch, original) {
		t.Fatalf("second ReviewBatch() changed batch = %#v, want %#v", batch, original)
	}
}
