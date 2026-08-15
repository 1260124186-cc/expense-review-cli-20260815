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
	sort.Slice(batch.Claims, func(i, j int) bool {
		return batch.Claims[i].ID < batch.Claims[j].ID
	})

	result := Review{
		Period:    batch.Period,
		Decisions: make([]Decision, 0, len(batch.Claims)),
	}
	for _, claim := range batch.Claims {
		decision := Decision{ClaimID: claim.ID, Status: StatusApproved}
		if claim.AmountCents > batch.Policy.capFor(claim.Category) {
			decision.Status = StatusRejected
			decision.Reason = "category cap exceeded"
		} else if claim.AmountCents >= batch.Policy.ReceiptThreshold && len(claim.ReceiptIDs) == 0 {
			decision.Status = StatusReview
			decision.Reason = "receipt required"
		}
		result.Decisions = append(result.Decisions, decision)
		result.Total += claim.AmountCents
	}
	return result, nil
}
