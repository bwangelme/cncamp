.PHONY: build

local_build:
	go build -o bin/httpserver .
build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/httpserver .
image: build
	docker build -t bwangel/cncamp_http_server:v1.0 .