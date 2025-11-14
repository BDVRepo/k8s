# syntax=docker/dockerfile:1

FROM golang:1.24 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o orders ./cmd/orders

FROM gcr.io/distroless/base-debian12
WORKDIR /srv

COPY --from=builder /app/orders ./orders

ENV ORDERS_GRPC_PORT=:50051
ENV PAYMENTS_GRPC_ADDR=payments:50052

EXPOSE 50051

ENTRYPOINT ["/srv/orders"]
