
ROOT := $(CURDIR)
BIN_NAME ?= candy
VERSION ?= 0.1
IMAGE_NAME ?= $(BIN_NAME):$(VERSION)
DOCKER_ID_USER ?= naughtytao

FULLNAME=$(DOCKER_ID_USER)/${BIN_NAME}:${VERSION}

install:
	export GO111MODULE=on
	go get github.com/golang/protobuf/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc
	#export PATH="$PATH:$(go env GOPATH)/bin"

gen:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		protobuf/service.proto
clean:
	rm -f protobuf/*.go server/server client/client zk/zk 

build:
	cd ./server ; env GOOS=linux GOARCH=amd64 go build
	cd ./client ; env GOOS=linux GOARCH=amd64 go build

docker: Dockerfile build
	docker build -t $(IMAGE_NAME) .

push:
	docker tag $(IMAGE_NAME) ${FULLNAME}
	docker push ${FULLNAME}docker run -it --name etcd bitnami/etcd

backend:
	docker run -d \
    	--publish 3000:2379 \
    	--publish 3001:2380 \
    	--env ALLOW_NONE_AUTHENTICATION=yes \
    	bitnami/etcd
