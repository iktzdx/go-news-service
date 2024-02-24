test: down build_test up e2e_test down

clean_test: clean build_test_nocache up e2e_test clean

e2e_test:
	docker-compose exec -T http go test -tags=e2e -v ./...

build_test:
	docker build . -t service_test -q --target=test

build_test_nocache:
	docker build . -t service_test -q --no-cache --target=test

down_nocache:
	docker-compose down --rmi all --volumes

clean: down_nocache
	[[ -z `docker ps -a -q` ]] || docker rm `docker ps -a -q`
	[[ -z `docker images --filter "dangling=true" -q --no-trunc` ]] \
		|| docker rmi `docker images --filter "dangling=true" -q --no-trunc`

unit_test:
	go test `go list ./... | grep -v test`

build:
	docker build . -t service

run:
	docker run --publish 8080:8080 service

up:
	docker-compose up -d

down:
	docker-compose down
