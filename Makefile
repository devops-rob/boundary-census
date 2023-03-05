generate_mocks:
	go install github.com/vektra/mockery/v2@v2.20.0
	go generate ./clients/boundary

run_dev:
	go run main.go -config=./example_config.hcl

build_app:
	go run ./dagger/*.go
