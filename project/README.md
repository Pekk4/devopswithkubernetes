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
kubectl create deployment todo_app --image=ghcr.io/pekk4/devopswithkubernetes-todo_app:1.2.
```
