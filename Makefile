SRC := src
BASE := $(SRC)/capsule
CAPSULE_MAIN := $(BASE)/main.go

all:
	@make build
	@make run

build:
	@go build $(CAPSULE_MAIN)

run:
	@go run $(CAPSULE_MAIN)