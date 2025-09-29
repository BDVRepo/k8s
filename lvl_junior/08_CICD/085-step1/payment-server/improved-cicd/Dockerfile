FROM golang:1.22 AS build
WORKDIR /app 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -buildvcs=false -o app .

FROM build AS unit-tests
RUN go test ./...

FROM build AS lint
RUN echo "lint stub"

FROM alpine:latest AS runtime
WORKDIR /app
COPY --from=build app .    
CMD ["./app"]
