bindata:
	go-bindata -pkg assets -o assets/assets.go config

build: bindata
	go build main.go

run: bindata
	go run main.go

install: bindata
	go install

test: bindata
	GO111MODULE=on go test -race -covermode=atomic -coverpkg=./... ./...
