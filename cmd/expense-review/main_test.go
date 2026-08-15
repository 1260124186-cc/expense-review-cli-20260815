package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRejectsDuplicateClaimIDsWithoutPublishingOutput(t *testing.T) {
	input := writeClaims(t, `{
		"period": "2026-08",
		"claims": [
			{"id":"meal-1","employee":"Ari","category":"meals","amount_cents":4200},
			{"id":"meal-1","employee":"Bo","category":"meals","amount_cents":4300}
		]
	}`)
	output := filepath.Join(t.TempDir(), "review.txt")
	if err := os.WriteFile(output, []byte("existing review\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout, stderr bytes.Buffer
	err := run([]string{"--input", input, "--output", output}, &stdout, &stderr)

	if err == nil {
		t.Fatal("run() error = nil, want duplicate claim rejection")
	}
	if got, want := err.Error(), `duplicate claim ID: "meal-1"`; !strings.Contains(got, want) {
		t.Fatalf("run() error = %q, want message containing %q", got, want)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want no successful review", got)
	}
	data, readErr := os.ReadFile(output)
	if readErr != nil {
		t.Fatalf("ReadFile() error = %v", readErr)
	}
	if got, want := string(data), "existing review\n"; got != want {
		t.Fatalf("output = %q, want unchanged %q", got, want)
	}
}

func TestRunPublishesReviewForUniqueClaimIDs(t *testing.T) {
	input := writeClaims(t, `{
		"period": "2026-08",
		"claims": [
			{"id":"meal-1","employee":"Ari","category":"meals","amount_cents":4200},
			{"id":"travel-1","employee":"Bo","category":"travel","amount_cents":81000,"receipt_ids":["r-1"]}
		]
	}`)
	output := filepath.Join(t.TempDir(), "review.txt")

	var stdout, stderr bytes.Buffer
	if err := run([]string{"--input", input, "--output", output}, &stdout, &stderr); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want output written to file only", got)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	for _, want := range []string{
		"period=2026-08 total_cents=85200",
		"meal-1=approved",
		"travel-1=rejected (category cap exceeded)",
	} {
		if got := string(data); !strings.Contains(got, want) {
			t.Errorf("published review missing %q:\n%s", want, got)
		}
	}
}

func writeClaims(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "claims.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
