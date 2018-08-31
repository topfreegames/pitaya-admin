TESTABLE_PACKAGES = `go list ./... | grep -v examples | grep -v constants | grep -v testing | grep -v cmd | grep -v protos`


setup:
	@dep ensure

run:
	@go run main.go start

run-grpc:
	@go run main.go start -g true -r 3939


test:
	cd ../pitaya/examples/demo/cluster/ && go run main.go && cd -   && sleep 1 &
	go test $(TESTABLE_PACKAGES) -coverprofile coverprofile.out && pkill main
