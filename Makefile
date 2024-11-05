# make
.PHONY: build
build: sec lint test

# updates dependencies in go.mod and rebuilds docs
.PHONY: update
update:
	rm -rf vendor
	go get -u ./...
	go mod tidy
	go mod vendor

# find the lint
.PHONY: lint
lint:
	@which golangci-lint || (echo "golangci-lint not in path" && false)
	# only scan source folders
	golangci-lint run --verbose ./...

# run units
.PHONY: test
test:
	go test -race -count=1 -shuffle=on -covermode=atomic -coverprofile=coverage.out ./...

# render and view coverage report in browser
.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out

# check the sec
.PHONY: sec
sec:
	@which govulncheck || (echo "govulncheck not in path" && false)
	govulncheck ./...
