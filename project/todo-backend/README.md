## Todo app

Run locally with:

```bash
cargo run
```

or with Docker:

```bash
docker build -t todo_backend .
docker run -it --rm todo_backend
```

Deploy in k8s with:

```bash
kubectl apply -f manifests/
```

## Endpoints

- `/todos` - GET, POST todos
