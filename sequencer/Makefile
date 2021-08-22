START_SEQ_FiLE = ../cmd/start_sequencer/start_sequencer.go
LOG_DIR = /tmp/sequencer-tests/
BIN = ./start_seq

test: build test-single-leader clean build test-lfl

test-lfl:
	SEQ_BIN=$(BIN) goreman -f procfiles/lfl.Procfile start &
	go test . -v -race -endpoints :8000

test-single-leader:
	SEQ_BIN=$(BIN) goreman -f procfiles/single_leader.Procfile start &
	go test . -v -race -endpoints :8000

build:
	go build -o $(BIN) ../cmd/start_sequencer/start_sequencer.go

clean:
	sudo pkill -f goreman
	rm $(BIN)
