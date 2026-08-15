package domain

import "context"

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

func ReviewBatch(ctx context.Context, batch Batch) (Review, error) {
	if err := batch.Validate(); err != nil {
		return Review{}, err
	}

	result := Review{
		Period:    batch.Period,
		Decisions: make([]Decision, 0, len(batch.Claims)),
	}
	for _, claim := range batch.Claims {
		// 调用方取消后及时停止，不再继续计算或产生 decision
		if err := ctx.Err(); err != nil {
			return Review{}, err
		}
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
