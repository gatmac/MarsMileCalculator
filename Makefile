.DEFAULT_GOAL := build
BINARY_NAME=MarsMileCalculator

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
	shadow ./...
.PHONY:vet

build: vet
	go build -o ${BINARY_NAME} .
.PHONY:build

run: vet
	go run .
.PHONY:run

clean:
	rm ${BINARY_NAME}
	rm ${BINARY_NAME}-windows-amd64.exe
	go clean
.PHONY:clean

compile: vet
	echo "Compiling for other OSs and Platforms"
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe .
.PHONY:compile