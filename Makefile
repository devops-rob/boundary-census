generate_mocks:
	go install github.com/vektra/mockery/v2@v2.20.0
	go generate ./clients/boundary

