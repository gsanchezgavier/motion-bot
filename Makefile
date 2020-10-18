all: build

build: clean test compile 

clean:
	@echo "==== Clean ===="
	rm -rf bin/

test:
	@echo "==== Test ===="
	go test -v ./...

compile:
	@echo "==== Compile ===="
	go build -v -o bin/bot cmd/motion-bot/motion-bot.go