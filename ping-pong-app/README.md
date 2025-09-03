## Ping pong app

Run locally with:

```bash
go run main.go
```

or with Docker:

```bash
docker build -t pingpong_app .
docker run -it --rm pingpong_app
```

Deploy in k8s with:

```bash
kubectl apply -f manifests/
```

## Endpoints

- `/pingpong` - returns the count of pongs requested
