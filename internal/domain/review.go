package domain

import "sort"

type Status string

const (
	StatusApproved Status = "approved"
	StatusReview   Status = "review"
	StatusRejected Status = "rejected"
)

type Decision struct {
	ClaimID string
	Status  Status
	Reason  string
}

type Review struct {
	Period    string
	Decisions []Decision
	Total     int64
}

func ReviewBatch(batch Batch) (Review, error) {
	if err := batch.Validate(); err != nil {
		return Review{}, err
	}
	reviewBatch := batch.Clone()
	sort.Slice(reviewBatch.Claims, func(i, j int) bool {
		return reviewBatch.Claims[i].ID < reviewBatch.Claims[j].ID
	})

	result := Review{
		Period:    reviewBatch.Period,
		Decisions: make([]Decision, 0, len(reviewBatch.Claims)),
	}
	for _, claim := range reviewBatch.Claims {
		decision := Decision{ClaimID: claim.ID, Status: StatusApproved}
		if claim.AmountCents > reviewBatch.Policy.capFor(claim.Category) {
			decision.Status = StatusRejected
			decision.Reason = "category cap exceeded"
		} else if claim.AmountCents >= reviewBatch.Policy.ReceiptThreshold && len(claim.ReceiptIDs) == 0 {
			decision.Status = StatusReview
			decision.Reason = "receipt required"
		}
		result.Decisions = append(result.Decisions, decision)
		result.Total += claim.AmountCents
	}
	return result, nil
}
