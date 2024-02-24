test: build_test
	docker-compose up -d
	docker-compose exec -T http go test -v ./...
	docker-compose down --rmi all --volumes

unit_test:
	go test `go list ./... | grep -v e2e_test`

build:
	docker build . -t service

build_test:
	docker build . -t service_test -q --no-cache --target=test

run:
	docker run --publish 8080:8080 service
