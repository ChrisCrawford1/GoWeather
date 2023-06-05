.PHONY: test
test:
	go test ./... -cover

.PHONY: generate-coverage
generate-coverage:
	go test ./... -coverprofile=coverage.out

.PHONY: show-coverage
show-coverage:
	go tool cover -html=coverage.out

.PHONY: dev
dev:
	go build && ./goweather

