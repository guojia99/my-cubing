all: go

go:
	go run main.go api


buildx:
	go build -o mycube main.go