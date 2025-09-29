## HTTP-client

```bash
go run main.go -url https://vk.com -interval 5
```

```bash
docker run -it http-client:0.0.1 ./app -url https://vk.com -interval 3

docker run -p 8080:8080 -it http-server:0.0.1
docker run --network host -it http-client:0.0.1 ./app -url http://localhost:8080/courses -interval 1

docker tag http-client:0.0.1 bdv21/http-client:0.0.1
docker push  bdv21/http-client:0.0.1
```