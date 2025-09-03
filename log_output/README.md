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
kubectl apple -f manifests/deployment.yaml
```

## Endpoints

- `/status` - returns the string with current timestamp

## Configuration

The app can run in two modes, determined by the `ROLE` environment variable. With `ROLE=writer` the app writes the timestamped random string into a log file, and without `ROLE` set it runs a web server that serves the log file content at the `/status` endpoint.
