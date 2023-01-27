BASE_DIR ?= ${PWD}
MONGO_CONTAINER_NAME ?= simpleapiwithmongodb-db
MONGO_VERSION ?= 6.0.3
API_DOCKER_IMG_NAME ?= simpleapiwithmongodb-api
API_CONTAINER_NAME ?= simpleapiwithmongodb-api
API_HOST_PORT ?= 8088


app-delete:
	@ API_RUNNING=`docker container ls --filter name=${API_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ ! -z "$$API_RUNNING" ]; then \
			docker stop ${API_CONTAINER_NAME}; \
		fi && \
		echo "application stoppped" && \
		make mongodb-stop
.PHONY: app-delete

app-run: docker-build mongodb-run
	@ API_EXISTS=`docker container ls -a --filter name=${API_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ -z "$$API_EXISTS" ]; then \
			docker run --rm --name ${API_CONTAINER_NAME} -p ${API_HOST_PORT}:${API_HOST_PORT} \
				-e DB_USERNAME=$$DB_USERNAME -e DB_PASSWORD=$$DB_PASSWORD -e DB_HOST=$$DB_HOST -e DB_PORT=$$DB_PORT \
				--link ${MONGO_CONTAINER_NAME}:$$DB_HOST \
				-d ${API_DOCKER_IMG_NAME} && docker container ls --filter name=${API_CONTAINER_NAME}; \
		fi
.PHONY: app-run

docker-build: go-test
	@ cd ${BASE_DIR} && \
		docker build -t ${API_DOCKER_IMG_NAME}:latest .
.PHONY: docker-build

mongodb-stop:
	@ MONGO_RUNNING=`docker container ls --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ ! -z "$$MONGO_RUNNING" ]; then \
			docker stop ${MONGO_CONTAINER_NAME}; \
		fi && \
		echo "mongodb stoppped"
.PHONY: mongodb-stop

go-run: go-build mongodb-run
	@ cd ${BASE_DIR} && \
		go run main.go
.PHONY: go-run

mongodb-run:
	@ MONGO_EXISTS=`docker container ls -a --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ -z "$$MONGO_EXISTS" ]; then \
			docker run --rm --name ${MONGO_CONTAINER_NAME} -p $$DB_PORT:$$DB_PORT \
				-e MONGO_INITDB_ROOT_USERNAME=$$DB_USERNAME -e MONGO_INITDB_ROOT_PASSWORD=$$DB_PASSWORD \
				-d mongo:${MONGO_VERSION} && docker container ls --filter name=${MONGO_CONTAINER_NAME}; \
		fi
.PHONY: mongodb-run

go-build: go-test
	@ cd ${BASE_DIR} && \
		mkdir -p bin && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o "bin/app" main.go
.PHONY: go-build

go-test: go-lint
	@ cd ${BASE_DIR} && \
		go test ./app -cover
.PHONY: go-test

go-lint: validate-envs
	@ cd ${BASE_DIR} && \
		golangci-lint run --timeout 30s
.PHONY: go-lint

validate-envs:
	@ if [ -z "$$DB_USERNAME" ] || [ -z "$$DB_PASSWORD" ] || [ -z "$$DB_HOST" ] || [ -z "$$DB_PORT" ]; then \
			echo "check for empty envs"; \
		fi
.PHONY: validate-envs