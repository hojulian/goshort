init: 
	rm -Rf .git
	find . -type f ! -name 'Makefile' -exec sed -i 's/{{pkgName}}/${pkgName}/g;' {} \;
	git init
	git add .
	git commit -m 'initial commit'

bench:
	go test -bench=. -benchmem ./...

test:
	GOCACHE=off go test $(shell go list ./... | grep -v /examples/ ) -covermode=count

test-race:
	GOCACHE=off go test -race $(shell go list ./... | grep -v /examples/ )

coverage:
	go test $(shell go list ./... | grep -v /examples/ ) -covermode=count -coverprofile=coverage.out && go tool cover -func=coverage.out

coverage-html:
	go test $(shell go list ./... | grep -v /examples/ ) -covermode=count -coverprofile=coverage.out && go tool cover -html=coverage.out

lint: 
	golint -set_exit_status $(shell (go list ./... | grep -v /vendor/))

.PHONY: bench test test-race coverage coverage-html lint