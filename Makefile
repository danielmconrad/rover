MAIN = ./cmd/marv/main.go
BIN  = ./bin/marv

.PHONY: run
run:
	go run ${MAIN} 

.PHONY: build
build:
	go build -o ${BIN} ${MAIN} 

.PHONY: build-and-run
build-and-run: build
	sudo ${BIN}

.PHONY: docker-build
docker-build: 
	docker build -t marv:local . 

.PHONY: docker-start
docker-start: docker-build
	docker run -it --privileged --restart unless-stopped marv:local
