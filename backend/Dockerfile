FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY ./pkg ./modules ./
COPY ./main.go ./main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /app


# Deploy the application binary into a lean image
FROM alpine:3.19

WORKDIR /

COPY ./pkg/database/migrations ./pkg/database/migrations
COPY --from=build-stage /app /app

EXPOSE 8080


CMD ["./app"]