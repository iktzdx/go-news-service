test: build_test up
	@-make -s test_all || (make -s down; exit 1)
	@make -s down

clean_test: build_test_nocache up
	@-make -s test_all || (make -s clean; exit 1)
	@make -s clean

test_all:
	docker-compose -f ./build/docker-compose.yaml exec -T http go test -v ./...

build_test:
	docker build -t service_test -q --target=test -f ./build/Dockerfile .

build_test_nocache:
	docker build -t service_test -q --no-cache --target=test -f ./build/Dockerfile .

down_nocache:
	docker-compose -f ./build/docker-compose.yaml down --rmi all --volumes

clean: down_nocache
	@CONTAINERS=$$(docker ps -a -q); \
	[[ -z $$CONTAINERS ]] || docker rm $$CONTAINERS
	@IMAGES=$$(docker images --filter "dangling=true" -q --no-trunc); \
	[[ -z $$IMAGES ]] || docker rmi $$IMAGES

unit_test:
	go test -v `go list ./... | grep -v test`

build:
	docker build -t service ./build/Dockerfile .

run:
	docker run --publish 8080:8080 service

up:
	docker-compose -f ./build/docker-compose.yaml up -d

down:
	docker-compose -f ./build/docker-compose.yaml down
