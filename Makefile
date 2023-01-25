BASE_DIR ?= ${PWD}
MONGO_CONTAINER_NAME ?= simpleapiwithmongodb
MONGO_HOST_PORT ?= 27018
MONGO_VERSION ?= 6.0.3
MONGO_USERNAME ?= admin
MONGO_PASSWORD ?= secret

go-run: go-build
	@ cd ${BASE_DIR} && \
		go run main.go
.PHONY: go-run

go-build: go-test
	@ cd ${BASE_DIR} && \
		mkdir -p bin && \
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

mongodb-stop:
	@ MONGO_RUNNING=`docker container ls --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ ! -z "$$MONGO_RUNNING" ]; then \
			docker stop ${MONGO_CONTAINER_NAME}; \
		fi
.PHONY: mongodb-stop

mongodb-run:
	@ MONGO_EXISTS=`docker container ls -a --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ -z "$$MONGO_EXISTS" ]; then \
			docker run --rm --name ${MONGO_CONTAINER_NAME} -p ${MONGO_HOST_PORT}:27017 \
				-e MONGO_INITDB_ROOT_USERNAME=${MONGO_USERNAME} -e MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWORD} \
				-d mongo:${MONGO_VERSION} && docker container ls --filter name=${MONGO_CONTAINER_NAME}; \
		fi
.PHONY: mongodb-run