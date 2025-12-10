.PHONY: run
run:
	go run ./cmd/product/main.go

.PHONY: swag
swag: 
	swag init -g ./cmd/product/main.go 

.PHONY: init
init: \
		swag \
    run \

.DEFAULT_GOAL := init