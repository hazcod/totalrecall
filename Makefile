all: run

prepare:
	go install github.com/goreleaser/goreleaser/v2@latest

run:
	go run ./cmd/...

build:
	$$GOPATH/bin/goreleaser build --config=.github/goreleaser.yml --clean --snapshot

clean:
	rm -r dist/ totalrecall || true
