include .envrc

# set default value for "reqs"
reqs ?= 100

.PHONY: run compile

compile:
	@echo "Compiling..."
	go build -o ./bin/benchmark benchmark.go

run:
	@echo "Running..."
	./bin/benchmark -n ${reqs} ${EDGE_URL} ${VM_URL}

	