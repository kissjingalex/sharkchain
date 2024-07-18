build:
	go build -o ./bin/shark

run: build
	./bin/shark

build2:
	go build -o ./bin/shark2 ./nodes/server2/main.go

run2: build2
	./bin/shark2

test:
	go test ./...