## Log output app

Run locally with:

```bash
go run main.go
```

or with Docker:

```bash
docker build -t log_output_app .
docker run -it --rm log_output_app
```

Deploy in k8s with:

```bash
kubectl create deployment log-output --image=ghcr.io/pekk4/devopswithkubernetes-log_output:1.1.
```
