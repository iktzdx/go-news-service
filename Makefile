test: build_test up e2e_test down

clean_test: build_test up e2e_test down_rm clean

unit_test:
	go test `go list ./... | grep -v test`

build:
	docker build . -t service

build_test:
	docker build . -t service_test -q --no-cache --target=test

run:
	docker run --publish 8080:8080 service

up:
	docker-compose up -d

down:
	docker-compose down

down_rm:
	docker-compose down --rmi all --volumes

e2e_test:
	docker-compose exec -T http go test -tags=e2e -v ./...

clean:
	docker images --filter "dangling=true" -q --no-trunc 2>/dev/null \
		&& docker rmi `docker images --filter "dangling=true" -q --no-trunc`
	docker ps -a -q 2>/dev/null \
		&& docker rm `docker ps -a -q`
