APP=file-duplicates-finder
PATHTOMAIN=cmd/${APP}/main.go

.PHONY: clear
clear:
	rm -f ${APP} || true

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race ./...

.PHONY: build-mac
build-mac: clear lint test
	env GOOS=darwin GOARCH=amd64 go build -o ${APP} ${PATHTOMAIN}

.PHONY: build-linux
build-linux: clear lint test
	env GOOS=linux GOARCH=386 go build -o ${APP} ${PATHTOMAIN}