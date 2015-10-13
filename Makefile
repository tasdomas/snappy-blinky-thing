build: bin/blink

bin/blink:
	env GOOS=linux GOARCH=arm go build -o bin/blink src/main.go
