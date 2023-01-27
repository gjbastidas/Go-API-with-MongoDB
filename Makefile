BASE_DIR ?= ${PWD}
MONGO_CONTAINER_NAME ?= simpleapiwithmongodb-db
MONGO_HOST_PORT ?= 27018
MONGO_VERSION ?= 6.0.3
MONGO_USERNAME ?= admin
MONGO_PASSWORD ?= secret
DOCKER_IMG_NAME ?= simpleapiwithmongodb-api
API_CONTAINER_NAME ?= simpleapiwithmongodb-api
API_HOST_PORT =? 8088

mongodb-stop:
	@ MONGO_RUNNING=`docker container ls --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ ! -z "$$MONGO_RUNNING" ]; then \
			docker stop ${MONGO_CONTAINER_NAME}; \
		fi
.PHONY: mongodb-stop

go-run: go-build mongodb-run
	@ cd ${BASE_DIR} && \
		go run main.go
.PHONY: go-run

mongodb-run:
	@ MONGO_EXISTS=`docker container ls -a --filter name=${MONGO_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ -z "$$MONGO_EXISTS" ]; then \
			docker run --rm --name ${MONGO_CONTAINER_NAME} -p ${MONGO_HOST_PORT}:27017 \
				-e MONGO_INITDB_ROOT_USERNAME=${MONGO_USERNAME} -e MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWORD} \
				-d mongo:${MONGO_VERSION} && docker container ls --filter name=${MONGO_CONTAINER_NAME}; \
		fi
.PHONY: mongodb-run

go-build: go-test
	@ cd ${BASE_DIR} && \
		mkdir -p bin && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o "bin/app" main.go
.PHONY: go-build

docker-run: # TODO
	@ API_EXISTS=`docker container ls -a --filter name=${API_CONTAINER_NAME} --format '{{.Names}}'` && \
		if [ -z "$$API_EXISTS" ]; then \
			docker run --rm --name ${API_CONTAINER_NAME} -p ${API_HOST_PORT}:${API_HOST_PORT} \
				# -e SVR_ADDR=$$SVR_ADDR -e DB_USERNAME=$$DB_USERNAME -e DB_PASSWORD=$$DB_PASSWORD -e DB_HOST=$$DB_HOST -e DB_NAME=$$DB_NAME \
				-d ${DOCKER_IMG_NAME} && docker container ls --filter name=${API_CONTAINER_NAME}; \
		fi
.PHONY: docker-run

docker-build: # TODO
	@ cd ${BASE_DIR} && \
		docker build -t ${DOCKER_IMG_NAME}:latest .
.PHONY: docker-build

go-test:
	@ cd ${BASE_DIR} && \
		go test -v ./app -cover
.PHONY: go-test

go-lint:
	@ cd ${BASE_DIR} && \
		golangci-lint run --timeout 360s
.PHONY: go-lint