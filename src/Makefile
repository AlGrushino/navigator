BUILD_DIR := build
BIN_DIR := bin

.PHONY: all clean test s21_graph s21_graph_algorithms stack queue app

all: s21_graph s21_graph_algorithms stack queue app

clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)

test:
	go test ./... -cover

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

stack: | $(BUILD_DIR)
	go build -buildmode=archive -o $(BUILD_DIR)/stack.a ./containers/stack

queue: | $(BUILD_DIR)
	go build -buildmode=archive -o $(BUILD_DIR)/queue.a ./containers/queue

s21_graph: | $(BUILD_DIR)
	go build -buildmode=archive -o $(BUILD_DIR)/s21_graph.a ./graph

s21_graph_algorithms: | $(BUILD_DIR)
	go build -buildmode=archive -o $(BUILD_DIR)/s21_graph_algorithms.a ./algorithms

app: | $(BIN_DIR)
	go build -o $(BIN_DIR)/graphcli ./cmd/graphcli
