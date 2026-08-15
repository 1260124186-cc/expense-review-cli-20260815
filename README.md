# Expense Review CLI

Expense Review CLI evaluates a JSON batch of employee expense claims against
simple category and receipt policies. It runs entirely locally and requires no
network services or database.

```sh
go run ./cmd/expense-review --input examples/claims.json
go test ./...
```

Use `--output review.txt` to publish the rendered review atomically.
