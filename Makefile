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
