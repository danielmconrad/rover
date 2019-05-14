.PHONY: install
install:
	go mod download

.PHONY: start
start:
	go run ./cmd/rover/main.go

.PHONY: deploy
deploy:
	docker build -t danielmconrad/rover:latest . 
	docker push danielmconrad/rover:latest
