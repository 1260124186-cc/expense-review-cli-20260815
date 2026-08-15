package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/1260124186-cc/expense-review-cli-20260815/internal/domain"
)

type JSONRepository struct{}

func NewJSONRepository() JSONRepository {
	return JSONRepository{}
}

type inputBatch struct {
	Period string       `json:"period"`
	Policy *inputPolicy `json:"policy"`
	Claims []inputClaim `json:"claims"`
}

type inputPolicy struct {
	CategoryCaps     map[string]int64 `json:"category_caps"`
	ReceiptThreshold *int64           `json:"receipt_threshold"`
}

type inputClaim struct {
	ID          string   `json:"id"`
	Employee    string   `json:"employee"`
	Category    string   `json:"category"`
	AmountCents int64    `json:"amount_cents"`
	ReceiptIDs  []string `json:"receipt_ids"`
}

func (JSONRepository) Load(ctx context.Context, path string) (domain.Batch, error) {
	if err := ctx.Err(); err != nil {
		return domain.Batch{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Batch{}, fmt.Errorf("read claim batch: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.Batch{}, err
	}

	var input inputBatch
	decoder := json.NewDecoder(bytesReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return domain.Batch{}, fmt.Errorf("decode claim batch: %w", err)
	}
	return toDomain(input), nil
}

func toDomain(input inputBatch) domain.Batch {
	policy := domain.DefaultPolicy()
	if input.Policy != nil {
		if input.Policy.CategoryCaps != nil {
			policy.CategoryCaps = input.Policy.CategoryCaps
		}
		if input.Policy.ReceiptThreshold != nil {
			policy.ReceiptThreshold = *input.Policy.ReceiptThreshold
		}
	}
	claims := make([]domain.Claim, 0, len(input.Claims))
	for _, claim := range input.Claims {
		claims = append(claims, domain.Claim{
			ID:          claim.ID,
			Employee:    claim.Employee,
			Category:    claim.Category,
			AmountCents: claim.AmountCents,
			ReceiptIDs:  append([]string(nil), claim.ReceiptIDs...),
		})
	}
	return domain.Batch{Period: input.Period, Policy: policy, Claims: claims}
}
