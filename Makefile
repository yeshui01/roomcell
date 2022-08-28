## Makefile

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go mod download

.PHONY: fmt
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: check
check: ## Run all the linters
	staticcheck ./...	

.PHONY: build
build: roomcell_account roomcell_hall roomcell_gate roomcell_data roomcell_room roomcell_roommgr

## 各进程编译目标
.PHONY: roomcell_account
roomcell_account:
	go build -o ./bin/roomcell_account ./cmd/account/main.go

.PHONY: roomcell_hall
roomcell_hall:
	go build -o ./bin/roomcell_hall ./cmd/hall_server/main.go

.PHONY: roomcell_gate
roomcell_gate:
	go build -o ./bin/roomcell_gate ./cmd/hall_gate/main.go

.PHONY: roomcell_data
roomcell_data:
	go build -o ./bin/roomcell_data ./cmd/hall_data/main.go

.PHONY: roomcell_room
roomcell_room:
	go build -o ./bin/roomcell_room ./cmd/hall_room/main.go

.PHONY: roomcell_roommgr
roomcell_roommgr:
	go build -o ./bin/roomcell_roommgr ./cmd/hall_roommgr/main.go


.PHONY:clean

clean: ## Remove temporary files
	go clean
	rm -rf bin

.PHONY:code
code:
	go test cmd/generate/*

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
