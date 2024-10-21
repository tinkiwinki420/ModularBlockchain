build:
	go build -o ./modularBlockchain

run: build
	./modularBlockchain

test:
	go test -v ./...
