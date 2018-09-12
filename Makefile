TESTABLE_PACKAGES = `go list ./... | grep -v examples | grep -v constants | grep -v testing | grep -v cmd | grep -v protos`

setup:
	@dep ensure

run:
	@go run main.go start

run-grpc:
	@go run main.go start -g true -r 3939

ensure-test-bin:
	@[ -f testing/server ] || go build -o testing/server examples/main.go

test-dep: ensure-test-bin
	@cd testing/ && docker-compose up -d
	@testing/server 2>/dev/null & echo $$! > testserver.pid

go-test:
	@echo "=========RUNNING UNIT TESTS==========="
	@sleep 2
	@-go test $(TESTABLE_PACKAGES) -coverprofile coverprofile.out -failfast -count=1

.SILENT:
kill-test-server:
	if [ -e testserver.pid ]; then \
		kill -TERM $$(cat testserver.pid) || true; \
		rm testserver.pid || true; \
	fi;

kill-testing-deps:
	@cd testing/ && docker-compose down; true

setup-ci:
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u github.com/wadey/gocovmerge
	@dep ensure

test: kill-testing-deps test-dep go-test kill-test-server
