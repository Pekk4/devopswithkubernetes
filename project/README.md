## Todo app

Run locally with:

```bash
go run main.go
```

or with Docker:

```bash
docker build -t todo_app .
docker run -it --rm todo_app
```

Deploy in k8s with:

```bash
kubectl apply -f manifests/
```
