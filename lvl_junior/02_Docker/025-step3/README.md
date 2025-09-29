## HTTP-server

Accepts requests and log URL path.

```bash
docker build . -t http-server:0.0.1
docker run -p 8080:8080 -it http-server:0.0.1
docker tag http-server:0.0.1 bdv21/http-server:0.0.1
docker push  bdv21/http-server:0.0.1
```