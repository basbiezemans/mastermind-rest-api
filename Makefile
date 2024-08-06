install:
	go install

test:
	go test ./...

http-test:
	GIN_MODE=test go test -v

run:
	go run main.go