setup:
	@dep ensure

run:
	@go run main.go start

run-grpc:
	@go run main.go start -g true -r 3939

