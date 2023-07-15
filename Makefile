all: go


go:
	go run main.go 14023

build:
	go build -v -o mycube main.go