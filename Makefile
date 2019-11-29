generate-mocks:
	mockgen -package=util -source=util/client.go -destination=util/mock_client.go
	mockgen -package=util -source=util/error.go -destination=util/mock_error.go
	mockgen -package=util -source=util/print.go -destination=util/mock_print.go

test:
	go test -v ./...
