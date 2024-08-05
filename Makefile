init:
	go mod download
	go install github.com/segmentio/golines@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

lint: 
	golangci-lint version
	golangci-lint -c ".golangci.yml" run --allow-parallel-runners ./...

lint-fix:
	golangci-lint -c ".golangci.yml" run --fix --allow-parallel-runners ./...

format:
	golines --base-formatter="goimports" -w -m 120 .
	gofumpt -w .

test-coverage-out:
	ENV=test go test -race -coverprofile=profile.cov -covermode=atomic `go list ./internal/./... | grep -v /mocks`
	go tool cover -func profile.cov

generate-proto:
	cd plugins/protobuf && protoc --go_out=./generated protobuf_response.proto