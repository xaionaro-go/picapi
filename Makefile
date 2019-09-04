all:
	@echo "Makefile is only for developers' needs"

coveralls:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	go test ./... -v -covermode=count -coverprofile=/tmp/picapi-coverage.out
	$(shell go env GOPATH)/bin/goveralls -coverprofile=/tmp/picapi-coverage.out -service=travis-ci -repotoken u8WyCVZQCpLA8fvWSlLVroeOR48erIxJ0
