NOW=`date +%Y%m%d%H%M%S`
OS=`uname -n -m`
AFTER_COMMIT=`git rev-parse HEAD`
GOPATH_DIR=`go env GOPATH`

install:
	go install .

check:
	@echo "running tests..."
	@go test -count 1 -coverprofile=coverage.txt -covermode=atomic -race -v ./test/...
	@echo "everything is OK"

.PHONY: check
