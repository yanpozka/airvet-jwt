run:
	go run main.go

build:
	go build -o server main.go

rotate:
	go run rotateKeys.go
