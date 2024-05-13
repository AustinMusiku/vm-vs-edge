include .envrc

.PHONY: run compile

compile:
	@echo "Compiling..."
	go build -o ./bin/benchmark benchmark.go

run:
	@echo "Running..."
	./bin/benchmark ${EDGE_URL} ${VM_URL}

	