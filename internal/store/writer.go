package store

import (
	"fmt"
	"os"
	"path/filepath"
)

type AtomicWriter struct{}

func NewAtomicWriter() AtomicWriter {
	return AtomicWriter{}
}

func (AtomicWriter) Write(path, content string) error {
	directory := filepath.Dir(path)
	file, err := os.CreateTemp(directory, ".expense-review-*")
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	temporary := file.Name()
	defer os.Remove(temporary)

	if _, err := file.WriteString(content); err != nil {
		file.Close()
		return fmt.Errorf("write output: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close output: %w", err)
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("publish output: %w", err)
	}
	return nil
}
