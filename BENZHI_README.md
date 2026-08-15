# BENZHI Docker Verification

The Docker image retains the Go toolchain so that validation happens inside the
container rather than against a host binary.

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

Each command builds the selected platform image, runs `go build ./...` inside
the container, and starts the CLI with the included sample claim batch.
