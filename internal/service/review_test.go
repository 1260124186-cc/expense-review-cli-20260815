package service_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1260124186-cc/expense-review-cli-20260815/internal/service"
	"github.com/1260124186-cc/expense-review-cli-20260815/internal/store"
)

func TestReviewAndRender(t *testing.T) {
	input := writeInput(t, `{
		"period": "2026-08",
		"claims": [
			{"id":"meal-1","employee":"Ari","category":"meals","amount_cents":6000,"receipt_ids":[]},
			{"id":"trip-1","employee":"Bo","category":"travel","amount_cents":81000,"receipt_ids":["r-9"]}
		]
	}`)
	reviewer := service.New(store.NewJSONRepository(), store.NewAtomicWriter())

	rendered, err := reviewer.ReviewAndRender(context.Background(), input)
	if err != nil {
		t.Fatalf("ReviewAndRender() error = %v", err)
	}
	for _, want := range []string{
		"period=2026-08 total_cents=87000",
		"meal-1=review (receipt required)",
		"trip-1=rejected (category cap exceeded)",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered output missing %q:\n%s", want, rendered)
		}
	}
}

func TestWritePublishesRenderedReview(t *testing.T) {
	output := filepath.Join(t.TempDir(), "review.txt")
	reviewer := service.New(store.NewJSONRepository(), store.NewAtomicWriter())
	if err := reviewer.Write(output, "period=2026-08 total_cents=1\n"); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got, want := string(data), "period=2026-08 total_cents=1\n"; got != want {
		t.Fatalf("published output = %q, want %q", got, want)
	}
}

func TestReviewRejectsDuplicateClaims(t *testing.T) {
	input := writeInput(t, `{
		"period": "2026-08",
		"claims": [
			{"id":"meal-1","employee":"Ari","category":"meals","amount_cents":4200,"receipt_ids":[]},
			{"id":"meal-1","employee":"Bo","category":"meals","amount_cents":4300,"receipt_ids":[]}
		]
	}`)
	reviewer := service.New(store.NewJSONRepository(), store.NewAtomicWriter())

	rendered, err := reviewer.ReviewAndRender(context.Background(), input)
	if err == nil {
		t.Fatal("ReviewAndRender() error = nil, want duplicate input rejection")
	}
	if rendered != "" {
		t.Fatalf("ReviewAndRender() output = %q, want no successful review", rendered)
	}
	if got, want := err.Error(), `duplicate claim ID: "meal-1"`; !strings.Contains(got, want) {
		t.Fatalf("ReviewAndRender() error = %q, want message containing %q", got, want)
	}
}

func writeInput(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "claims.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
