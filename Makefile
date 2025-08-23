run: build
	./bin/cox-event-listener

build:
	go build -o bin/cox-event-listener src/main.go
