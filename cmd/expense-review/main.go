package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/1260124186-cc/expense-review-cli-20260815/internal/service"
	"github.com/1260124186-cc/expense-review-cli-20260815/internal/store"
)

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("expense-review", flag.ContinueOnError)
	flags.SetOutput(stderr)
	input := flags.String("input", "", "path to a claim batch JSON file")
	output := flags.String("output", "", "optional path for the rendered review")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *input == "" {
		return fmt.Errorf("--input is required")
	}

	reviewer := service.New(store.NewJSONRepository(), store.NewAtomicWriter())
	rendered, err := reviewer.ReviewAndRender(context.Background(), *input)
	if err != nil {
		return err
	}
	if *output != "" {
		return reviewer.Write(*output, rendered)
	}
	_, err = fmt.Fprint(stdout, rendered)
	return err
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
