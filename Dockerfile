# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /transaction-processor-rest ./cmd/rest

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /transaction-processor-rest /transaction-processor-rest
COPY --from=build-stage /app/migrations /migrations

EXPOSE 8080 8081

ENTRYPOINT ["/transaction-processor-rest"]