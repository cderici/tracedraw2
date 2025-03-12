.PHONY: all cli test clean

all: cli

# Build the CLI binary
cli:
	go build -o tracedraw2-cli ./cmd/tracedraw2-cli

# Build the server binary
# server:
# 	go build -o tracedraw2-server ./cmd/tracedraw2-server

# Run all tests
test:
	go test ./...

# Remove built binaries
clean:
	rm -f tracedraw2-cli

