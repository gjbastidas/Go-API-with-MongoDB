BASE_DIR ?= ${PWD}

go-run: go-build
	@ cd ${BASE_DIR} && \
		go run main.go
.PHONY: go-run

go-build: go-test
	@ cd ${BASE_DIR} && \
		mkdir -p bin
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o "bin/app" main.go
.PHONY: go-build

go-test: go-gen-mocks
	@ cd ${BASE_DIR} && \
		go test -v ./app
.PHONY: go-test

go-gen-mocks:
	@ cd ${BASE_DIR}/models && \
		${GOPATH}/bin/mockgen -destination=${BASE_DIR}/mocks/mock_post.go -package=mocks . Post
.PHONY: go-gen-mocks

go-lint:
	@ cd ${BASE_DIR} && \
		golangci-lint run --timeout 360s
.PHONY: go-lint