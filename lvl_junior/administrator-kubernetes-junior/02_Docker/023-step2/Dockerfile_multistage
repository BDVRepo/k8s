FROM golang:1.20 AS build

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest    
                   
WORKDIR /app/

COPY --from=build ./app .             

CMD ["./app"]