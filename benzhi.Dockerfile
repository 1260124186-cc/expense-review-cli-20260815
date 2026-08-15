FROM golang:1.26.2-bookworm

ENV GOTOOLCHAIN=local
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["go", "run", "./cmd/expense-review", "--input", "examples/claims.json"]
