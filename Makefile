build: bin/blink
	snappy build .

bin/blink: src/main.go
	env GOOS=linux GOARCH=arm go build -o bin/blink src/main.go
