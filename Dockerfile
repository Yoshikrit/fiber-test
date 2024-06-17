FROM golang:1.22.1 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /fiber-test ./main.go

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine:3.20.0 AS build-release-stage
WORKDIR /

COPY --from=build-stage /fiber-test /fiber-test

COPY prod.env ./
ENV APP_ENV=prod

EXPOSE 8080

ENTRYPOINT ["/fiber-test", "prod"]