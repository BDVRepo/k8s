## Commands

```bash
docker build . -t docsify-cli
docker run -v .:/docs -it docsify-cli docsify init ./docs
sudo chown -R vboxuser:vboxuser docs
docker run -v .:/docs -p 3000:3000 --init -it docsify-cli docsify serve ./docs

## uncomment lines in Dockerfile
docker build . -t frontend:0.0.1
docker run -it -p 3000:3000 frontend:0.0.1
docker tag frontend:0.0.1 bdv21/frontend:0.0.1
docker push bdv21/frontend:0.0.1
```