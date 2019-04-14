MAIN = ./cmd/marv/main.go
BIN  = ./bin/marv

.PHONY: build
build:
	go build -o ${BIN} ${MAIN} 

.PHONY: run
run: build
	sudo ${BIN}
