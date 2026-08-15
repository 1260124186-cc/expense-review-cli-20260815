package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtomicWriterPublishesCompleteContent(t *testing.T) {
	output := filepath.Join(t.TempDir(), "review.txt")
	content := strings.Repeat("approved\n", 2048)

	if err := NewAtomicWriter().Write(output, content); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := string(data); got != content {
		t.Fatalf("published output has %d bytes, want %d", len(got), len(content))
	}
}

func TestAtomicWriterReturnsPublishFailure(t *testing.T) {
	directory := t.TempDir()
	output := filepath.Join(directory, "review.txt")
	if err := os.Mkdir(output, 0o700); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	err := NewAtomicWriter().Write(output, "period=2026-08 total_cents=1\n")
	if err == nil {
		t.Fatal("Write() error = nil, want publish failure")
	}
	if !strings.Contains(err.Error(), "publish output") {
		t.Fatalf("Write() error = %v, want publish output error", err)
	}

	temporary, err := filepath.Glob(filepath.Join(directory, ".expense-review-*"))
	if err != nil {
		t.Fatalf("Glob() error = %v", err)
	}
	if len(temporary) != 0 {
		t.Fatalf("temporary output files = %v, want none", temporary)
	}
}
