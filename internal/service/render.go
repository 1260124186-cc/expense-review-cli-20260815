package service

import (
	"fmt"
	"strings"

	"github.com/1260124186-cc/expense-review-cli-20260815/internal/domain"
)

func Render(review domain.Review) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "period=%s total_cents=%d\n", review.Period, review.Total)
	for _, decision := range review.Decisions {
		if decision.Reason == "" {
			fmt.Fprintf(&builder, "%s=%s\n", decision.ClaimID, decision.Status)
			continue
		}
		fmt.Fprintf(&builder, "%s=%s (%s)\n", decision.ClaimID, decision.Status, decision.Reason)
	}
	return builder.String()
}
