.PHONY: build vet tidy clean local

build:
	go build -o bin/bot ./api

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/

local:
	go run ./cmd/local
