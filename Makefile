build:
	go build main.go

run:
	go run main.go

install:
	go install

test:
	GO111MODULE=on go test -race -covermode=atomic -coverpkg=./... ./...
