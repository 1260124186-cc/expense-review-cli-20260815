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
	// 在克隆副本上排序，避免改写调用方/仓储保留的原始批次顺序
	ordered := batch.Clone()
	sort.Slice(ordered.Claims, func(i, j int) bool {
		return ordered.Claims[i].ID < ordered.Claims[j].ID
	})

	result := Review{
		Period:    ordered.Period,
		Decisions: make([]Decision, 0, len(ordered.Claims)),
	}
	for _, claim := range ordered.Claims {
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
