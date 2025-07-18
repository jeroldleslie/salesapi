dep: ## Get all dependencies
	go env -w GOPROXY=direct
	go env -w GOSUMDB=off
	go mod download
	go mod tidy

build-run: dep ## Build and run
	go build -mod=mod -race -o salesapi .
	./salesapi
