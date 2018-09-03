TESTABLE_PACKAGES = `go list ./... | grep -v examples | grep -v constants | grep -v testing | grep -v cmd | grep -v protos`


setup:
	@dep ensure

run:
	@go run main.go start

run-grpc:
	@go run main.go start -g true -r 3939

test-dep:
	@cd examples/ && go run main.go 2>/dev/null && cd - & sleep 5

test: test-dep
	@echo "=========RUNNING UNIT TESTS==========="
	@go test $(TESTABLE_PACKAGES) -coverprofile coverprofile.out -failfast && pkill main
