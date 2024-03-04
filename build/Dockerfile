# Build the application from source
FROM golang:latest as build-stage

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . /service

RUN chmod +x build/pgsql/entrypoint.sh
RUN apt update -yq
RUN apt install -y postgresql-client

RUN CGO_ENABLED=0 GOOS=linux go build -o /http-service

ENTRYPOINT ["build/pgsql/entrypoint.sh"]

# Run the tests in the container
FROM build-stage as test

# Deploy the application binary into a lean image
FROM alpine:latest as release

RUN apk --no-cache add ca-certificates

COPY --from=build-stage /http-service ./

EXPOSE 8080

ENTRYPOINT ["/http-service"]
