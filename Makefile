all: go

go:
	go run main.go api

test_go:
	go run main.go api --config ./etc/configs_sqlite.json

test_admin:
	go run main.go admin --config ./etc/configs_sqlite.json

buildx:
	go build -o mycube main.go

