.PHONY: test

build:
	go build -o bin/gazur main.go

run:
	go run main.go --cfg-file $(FILE)
test:
	go test ./...
