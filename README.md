# template-service
A service for storing and retrieving text templates

Build the image with the following command:

```bash
docker build . -f Dockerfile.dev -t template-service:debug
```

Run the container with the image we just build with this command:

```bash
docker run -p 8081:8080 -p 4001:4000 -d --name template-service --env-file ./debug.env template-service:debug
```
