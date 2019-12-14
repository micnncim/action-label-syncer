export GO111MODULE := on

test:
	go test -race ./...

mock:
	mockgen -source pkg/github -destination pkg/github/mock
