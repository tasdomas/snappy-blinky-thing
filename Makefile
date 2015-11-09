build: bin/failer
	snappy build .

bin/failer: src/failer.go
	env GOOS=linux GOARCH=arm go build -o bin/failer src/failer.go
