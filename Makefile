build:
	go build -o ./bin/shark

run: build
	./bin/shark

test:
	go test ./...