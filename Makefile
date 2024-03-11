test: ### Run unit tests
	go test -v -cover -race ./internal/...

build:
	go build -o service cmd/main.go



# deps:
# 	go get github.com/mayu13/gymshark-task/internal/config

run:
	go run cmd/main.go