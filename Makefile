build:
	go build -o ./ModularBlockchain

run: build
	./ModularBlockchain

test:
	go test -v ./...
