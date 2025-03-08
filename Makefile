CC = go
OUTPUT = bin/tokino

build:
	$(CC) build -o $(OUTPUT) src/main.go
run:
	$(MAKE) build
	$(OUTPUT)
deps:
	go mod download