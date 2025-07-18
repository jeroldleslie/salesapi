dep: ## Get all dependencies
	go env -w GOPROXY=direct
	go env -w GOSUMDB=off
	go mod download
	go mod tidy

setup-mockgen:
	go install github.com/golang/mock/mockgen@latest
setup: setup-mockgen dep

build-run: dep ## Build and run
	go build -mod=mod -race -o salesapi .
	./salesapi
