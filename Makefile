run:
   	air -c dev.air.toml
test:
   	go test ./...
test-integration_test:
	go test github.com/Yoshikrit/fiber-test/integration_test -v
build:
   	docker build -t yoshikrit/fiber-test . 
compose:
	docker-compose up -d
