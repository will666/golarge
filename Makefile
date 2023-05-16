PROD_FLAGS = -ldflags='-w -s'
DEBUG_FLAGS = -gcflags='-m -l'
OUTPUT = ./build
BIN = golarge
TARGET = .

.PHONY: build run prod bench debug clean

build:
	@go build -o $(OUTPUT)/$(BIN) $(TARGET)

run:
	@go run $(TARGET) $(DIR)

prod:
	@go build $(PROD_FLAGS) -o $(OUTPUT)/$(BIN) $(TARGET)

bench:
	go test -bench $(TARGET) -benchmem

debug:
	go run $(DEBUG_FLAGS) $(TARGET) $(DIR)

clean:
	rm -rf logs *.txt *.json $(OUTPUT)/*
